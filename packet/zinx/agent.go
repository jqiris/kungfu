package zinx

import (
	"fmt"
	"github.com/apex/log"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/packet"
	"net"
	"sync/atomic"
)

type Agent struct {
	conn  net.Conn
	state int32 // current Agent state
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
	decoder     *Decoder // binary decoder
	lastAt      int64    // last msg time stamp
}

func newAgent(conn net.Conn) *Agent {
	cfg := config.GetConnectorConf()
	return &Agent{
		conn:        conn,
		state:       packet.StatusStart,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, cfg.MaxMsgChanLen),
		decoder:     NewDecoder(),
	}
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
	return a.conn.Close()
}

// RemoteAddr, implementation for session.NetworkEntity interface
// returns the remote network address.
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
