package tcpserver

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/packet/nano"
	"github.com/jqiris/kungfu/v2/packet/zinx"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/tcpface"
	"github.com/jqiris/kungfu/v2/treaty"
)

// Server 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//配置
	Config config.ConnectorConf
	//当前Server的链接管理器
	ConnMgr tcpface.IConnManager
	//消息处理器
	MsgHandler tcpface.IMsgHandle
	//连接处理器
	ConnHandler tcpface.IConnHandler
	//该Server的连接创建时Hook函数
	OnConnStart func(conn tcpface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn tcpface.IConnection)
	//服务配置信息
	Server *treaty.Server
}

// NewServer 创建一个服务器句柄
func NewServer(server *treaty.Server) tcpface.IServer {
	cfg := config.GetConnectorConf()
	s := &Server{
		Name:      server.ServerName,
		IPVersion: "tcp4",
		IP:        server.ServerIp,
		Port:      int(server.ClientPort),
		ConnMgr:   NewConnManager(),
		Config:    cfg,
		Server:    server,
	}
	//开启一个go去做服务端Lister业务
	var msgHandler tcpface.IMsgHandle
	var connHandler tcpface.IConnHandler
	switch cfg.UseType {
	case "zinx":
		msgHandler = zinx.NewMsgHandle()
		connHandler = func(server tcpface.IServer, conn net.Conn, connId int) tcpface.IConnection {
			return zinx.NewAgent(server, conn, connId)
		}
	case "nano":
		msgHandler = nano.NewMsgHandle()
		connHandler = func(server tcpface.IServer, conn net.Conn, connId int) tcpface.IConnection {
			return nano.NewAgent(server, conn, connId)
		}
	default:
		logger.Fatalf("no suitable connector type:%v", cfg.UseType)
	}
	s.MsgHandler = msgHandler
	s.ConnHandler = connHandler
	return s
}

//============== 实现 tcpface.IServer 里的全部接口方法 ========

// Start 开启网络服务
func (s *Server) Start() {
	logger.Infof("[START] Server name: %s,listenner at IP: %s, Port %d is starting", s.Name, s.IP, s.Port)
	go func() {
		if s.Config.UseWebsocket {
			s.ListenAndServeWs(s.MsgHandler, s.ConnHandler)
		} else {
			s.ListenAndServe(s.MsgHandler, s.ConnHandler)
		}
	}()
}

func (s *Server) ListenAndServe(msgHandler tcpface.IMsgHandle, connHandler tcpface.IConnHandler) {
	//0 启动worker工作池机制
	msgHandler.StartWorkerPool()

	//1 获取一个TCP的Addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf(":%d", s.Port))
	if err != nil {
		logger.Info("resolve tcp addr err: ", err)
		return
	}

	//2 监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		logger.Info("listen", s.IPVersion, "err", err)
		return
	}

	//已经监听成功
	logger.Info("start tcpserver server  ", s.Name, " succ, now listenning...")

	cid := 0

	//3 启动server网络连接业务
	for {
		//3.1 阻塞等待客户端建立连接请求
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.Info("Accept err ", err)
			continue
		}
		logger.Info("Get conn remote addr = ", conn.RemoteAddr().String())

		//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
		if s.Config.MaxConn > 0 && s.ConnMgr.Len() >= s.Config.MaxConn {
			err = conn.Close()
			if err != nil {
				logger.Error(err.Error())
			}
			continue
		}
		cid++

		agent := connHandler(s, conn, cid)
		go msgHandler.Handle(agent)
	}
}

func (s *Server) ListenAndServeWs(msgHandler tcpface.IMsgHandle, connHandler tcpface.IConnHandler) {
	//0 启动worker工作池机制
	msgHandler.StartWorkerPool()

	//1 获取一个TCP的Addr
	addr := fmt.Sprintf(":%d", s.Port)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(_ *http.Request) bool { return true },
	}
	cid := 0
	http.HandleFunc("/"+strings.TrimPrefix(s.Config.WebsocketPath, "/"), func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Infof("Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
			return
		}

		if wc, err := newWSConn(conn); err == nil {
			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.Config.MaxConn > 0 && s.ConnMgr.Len() >= s.Config.MaxConn {
				err = wc.Close()
				if err != nil {
					logger.Error(err.Error())
				}
				return
			}
			cid++
			agent := connHandler(s, wc, cid)
			go msgHandler.Handle(agent)
		} else {
			logger.Fatal(err)
		}
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Fatal(err.Error())
	}
}

// Stop 停止服务
func (s *Server) Stop() {
	logger.Info("[STOP] tcpserver server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() tcpface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(tcpface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(tcpface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn tcpface.IConnection) {
	if s.OnConnStart != nil {
		logger.Info("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn tcpface.IConnection) {
	if s.OnConnStop != nil {
		logger.Info("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

func (s *Server) GetMsgHandler() tcpface.IMsgHandle {
	return s.MsgHandler
}

func (s *Server) GetServerID() string {
	return s.Server.ServerId
}
