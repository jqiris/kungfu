package zinx

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/packet"
	"github.com/jqiris/kungfu/tcpface"
	"net"
	"sync"
	"sync/atomic"
)

type Agent struct {
	sync.RWMutex
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
		state:       packet.StatusWorking,
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
	logger.Info("[Writer Goroutine is running]")
	defer func() {
		err := a.Close()
		if err != nil {
			log.Error(err.Error())
		}
		logger.Info(a.conn.RemoteAddr().String(), "[conn Writer exit!]")
	}()
	for {
		select {
		case data := <-a.msgChan:
			//有数据要写给客户端
			if _, err := a.conn.Write(data); err != nil {
				logger.Info("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			//logger.Infof("Send data succ! data = %+v\n", data)
		case data, ok := <-a.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := a.conn.Write(data); err != nil {
					logger.Info("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				logger.Info("msgBuffChan is Closed")
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

//SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (a *Agent) SendMsg(msgID int32, data []byte) error {
	a.RLock()
	defer a.RUnlock()
	if a.status() == packet.StatusClosed {
		return errors.New("connection closed when send msg")
	}
	pk, err := Encode(&Message{
		Id:   msgID,
		Data: data,
	})
	if err != nil {
		return err
	}
	//写回客户端
	a.msgChan <- pk
	return nil
}

//SendBuffMsg  发生BuffMsg
func (a *Agent) SendBuffMsg(msgID int32, data []byte) error {
	a.RLock()
	defer a.RUnlock()
	if a.status() == packet.StatusClosed {
		return errors.New("connection closed when send msg")
	}

	pk, err := Encode(&Message{
		Id:   msgID,
		Data: data,
	})
	if err != nil {
		return err
	}
	//写回客户端
	a.msgBuffChan <- pk
	return nil
}
