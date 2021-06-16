package zinx

import (
	"fmt"
	"github.com/apex/log"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/packet"
	"github.com/jqiris/kungfu/tcpface"
	"net"
	"sync/atomic"
)

type Agent struct {
	//当前Server的链接管理器
	server tcpface.IServer
	conn   net.Conn
	connId int
	state  int32 // current Agent state
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
	decoder     *Decoder // binary decoder
	lastAt      int64    // last msg time stamp
}

func NewAgent(server tcpface.IServer, conn net.Conn, connId int) *Agent {
	cfg := config.GetConnectorConf()
	agent := &Agent{
		server:      server,
		conn:        conn,
		connId:      connId,
		state:       packet.StatusStart,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, cfg.MaxMsgChanLen),
		decoder:     NewDecoder(),
	}
	agent.server.GetConnMgr().Add(agent)
	agent.server.CallOnConnStart(agent)
	return agent
}

func (a *Agent) GetConnID() int {
	return a.connId
}

func (a *Agent) GetConn() net.Conn {
	return a.conn
}

/*
	写消息Goroutine， 用户将数据发送给客户端
*/
func (a *Agent) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer func() {
		err := a.Close()
		if err != nil {
			log.Error(err.Error())
		}
		fmt.Println(a.conn.RemoteAddr().String(), "[conn Writer exit!]")
	}()
	for {
		select {
		case data := <-a.msgChan:
			//有数据要写给客户端
			if _, err := a.conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			//fmt.Printf("Send data succ! data = %+v\n", data)
		case data, ok := <-a.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := a.conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				return
			}
		}
	}
}

func (a *Agent) Close() error {
	if a.status() == packet.StatusClosed {
		return packet.ErrCloseClosedSession
	}
	a.setStatus(packet.StatusClosed)
	close(a.msgChan)
	close(a.msgBuffChan)
	a.server.GetConnMgr().Remove(a) //从管理器移除
	a.server.CallOnConnStop(a)      //连接关闭事件
	return a.conn.Close()
}

// RemoteAddr  implementation for session.NetworkEntity interface
func (a *Agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

// String, implementation for Stringer interface
func (a *Agent) String() string {
	return fmt.Sprintf("Remote=%s, LastTime=%d", a.conn.RemoteAddr().String(), a.lastAt)
}

func (a *Agent) status() int32 {
	return atomic.LoadInt32(&a.state)
}

func (a *Agent) setStatus(state int32) {
	atomic.StoreInt32(&a.state, state)
}
