package tcpserver

import (
	"fmt"
	"net"
	"os"

	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/treaty"
)

//Server 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	Name string
	//服务器ID
	ID string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int32
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler IMsgHandle
	//当前Server的链接管理器
	ConnMgr IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn IConnection)
	//客户端参数
	Config conf.ConnectorConf
}

func MergeConf(cfg conf.ConnectorConf) conf.ConnectorConf {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	config := conf.ConnectorConf{
		Version:          "V0.11",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		LogDir:           pwd + "/log",
		LogFile:          "",
	}
	if cfg.MaxConn > 0 {
		config.MaxConn = cfg.MaxConn
	}
	if cfg.MaxPacketSize > 0 {
		config.MaxPacketSize = cfg.MaxPacketSize
	}
	if cfg.WorkerPoolSize > 0 {
		config.WorkerPoolSize = cfg.WorkerPoolSize
	}
	if cfg.MaxWorkerTaskLen > 0 {
		config.MaxWorkerTaskLen = cfg.MaxWorkerTaskLen
	}
	if cfg.MaxMsgChanLen > 0 {
		config.MaxMsgChanLen = cfg.MaxMsgChanLen
	}
	if len(cfg.LogDir) > 0 {
		config.LogDir = cfg.LogDir
	}
	if len(cfg.LogFile) > 0 {
		config.LogFile = cfg.LogFile
	}
	config.LogDebugClose = cfg.LogDebugClose
	return config
}

//NewServer 创建一个服务器句柄
func NewServer(server *treaty.Server, cfg conf.ConnectorConf) IServer {
	//合并配置
	cfg = MergeConf(cfg)
	//其他操作
	s := &Server{
		Name:       server.ServerName,
		ID:         server.ServerId,
		IPVersion:  "tcp4",
		IP:         "0.0.0.0",
		Port:       server.ClientPort,
		msgHandler: NewMsgHandle(cfg),
		ConnMgr:    NewConnManager(),
		Config:     cfg,
	}
	return s
}

//============== 实现 IServer 里的全部接口方法 ========

//Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	//开启一个go去做服务端Linster业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		//已经监听成功
		fmt.Println("start Zinx server  ", s.Name, " succ, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cID uint32 = 0

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= s.Config.MaxConn {
				conn.Close()
				continue
			}

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConntion(s, conn, cID, s.msgHandler, s.Config)
			cID++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) GetServerID() string {
	return s.ID
}

//Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

//Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select {}
}

//AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

//AddRouters 批量注册消息
func (s *Server) AddRouters(routers map[uint32]IRouter) {
	s.msgHandler.AddRouters(routers)
}

//GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() IConnManager {
	return s.ConnMgr
}

//SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(IConnection)) {
	s.OnConnStart = hookFunc
}

//SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(IConnection)) {
	s.OnConnStop = hookFunc
}

//CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
