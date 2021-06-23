package nano

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/packet"
	"github.com/jqiris/kungfu/serialize"
	"github.com/jqiris/kungfu/session"
	"github.com/jqiris/kungfu/tcpface"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type Agent struct {
	sync.RWMutex
	session *session.Session // session
	//当前Server的链接管理器
	server  tcpface.IServer
	conn    net.Conn
	connId  int
	lastMid uint  // last message id
	state   int32 // current Agent state
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan       chan []byte
	decoder           *Decoder             // binary decoder
	lastAt            int64                // last msg time stamp
	srv               reflect.Value        // cached session reflect.Value
	serializer        serialize.Serializer //序列化对象
	heartbeatInterval int                  //心跳间隔
	hbd               []byte               // heartbeat packet data
	chDie             chan struct{}        // wait for close
}

type pendingMessage struct {
	typ     MsgType     // message type
	route   string      // message route(push)
	mid     uint        // response message id(response)
	payload interface{} // payload
}

func (a *Agent) MID() uint {
	return a.lastMid
}

// Push Push, implementation for session.NetworkEntity interface
func (a *Agent) Push(route string, v interface{}) error {
	if a.status() == packet.StatusClosed {
		return packet.ErrBrokenPipe
	}

	return a.SendBuffMsg(pendingMessage{typ: Push, route: route, payload: v})
}

// Response Response, implementation for session.NetworkEntity interface
// Response message to session
func (a *Agent) Response(v interface{}) error {
	return a.ResponseMID(a.lastMid, v)
}

// ResponseMID Response, implementation for session.NetworkEntity interface
// Response message to session
func (a *Agent) ResponseMID(mid uint, v interface{}) error {
	if a.status() == packet.StatusClosed {
		return packet.ErrBrokenPipe
	}

	if mid <= 0 {
		return packet.ErrSessionOnNotify
	}
	return a.SendBuffMsg(pendingMessage{typ: Response, mid: mid, payload: v})
}

func NewAgent(server tcpface.IServer, conn net.Conn, connId int) *Agent {
	cfg := config.GetConnectorConf()
	maxMsgChanLen := 1024
	if cfg.MaxMsgChanLen > 0 {
		maxMsgChanLen = int(cfg.MaxMsgChanLen)
	}
	a := &Agent{
		server:            server,
		conn:              conn,
		connId:            connId,
		state:             packet.StatusStart,
		msgChan:           make(chan []byte),
		msgBuffChan:       make(chan []byte, maxMsgChanLen),
		decoder:           NewDecoder(),
		heartbeatInterval: cfg.HeartbeatInterval,
		lastAt:            time.Now().Unix(),
		chDie:             make(chan struct{}),
	}
	a.server.GetConnMgr().Add(a)
	a.server.CallOnConnStart(a)
	s := session.NewSession(a.connId, a)
	a.session = s
	a.srv = reflect.ValueOf(s)
	switch cfg.UseSerializer {
	case "proto":
		a.serializer = serialize.NewProtoSerializer()
	case "json":
		a.serializer = serialize.NewJsonSerializer()
	default:
		logger.Fatalf("no suitable serializer:%v", cfg.UseSerializer)
	}

	return a
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
	heartbeat := time.Duration(a.heartbeatInterval) * time.Second
	ticker := time.NewTicker(heartbeat)
	defer func() {
		ticker.Stop()
		err := a.Close()
		if err != nil {
			log.Error(err.Error())
		}
		logger.Info(a.conn.RemoteAddr().String(), "[conn Writer exit!]")
	}()

	for {
		select {
		case <-ticker.C:
			deadline := time.Now().Add(-2 * heartbeat).Unix()
			if a.lastAt < deadline {
				logger.Infof("Session heartbeat timeout, LastTime=%d, Deadline=%d", a.lastAt, deadline)
				return
			}
			a.msgBuffChan <- hbd
		case data, ok := <-a.msgChan:
			//有数据要写给客户端
			if ok {
				if _, err := a.conn.Write(data); err != nil {
					logger.Info("Send Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				logger.Info("msgChan is Closed")
				return
			}

			//fmt.Printf("Send data succ! data = %+v\n", data)
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
		case <-a.chDie:
			logger.Info("received agent closed signal")
			return
		}
	}
}

func (a *Agent) Close() error {
	if a.status() == packet.StatusClosed {
		return packet.ErrCloseClosedSession
	}
	a.setStatus(packet.StatusClosed)
	close(a.chDie)
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

func (a *Agent) SendRawMessage(useBuffer bool, data []byte) error {
	a.RLock()
	defer a.RUnlock()
	if a.status() == packet.StatusClosed {
		return errors.New("connection closed when send msg")
	}
	//写回客户端
	if useBuffer {
		a.msgBuffChan <- data
	} else {
		a.msgChan <- data
	}
	return nil
}

//SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (a *Agent) SendMsg(data pendingMessage) error {
	a.RLock()
	defer a.RUnlock()
	if a.status() == packet.StatusClosed {
		return errors.New("connection closed when send msg")
	}
	pk, err := a.serializeOrRaw(data)
	if err != nil {
		return err
	}
	//写回客户端
	a.msgChan <- pk
	return nil
}

//SendBuffMsg  发生BuffMsg
func (a *Agent) SendBuffMsg(data pendingMessage) error {
	a.RLock()
	defer a.RUnlock()
	if a.status() == packet.StatusClosed {
		return errors.New("connection closed when send msg")
	}

	pk, err := a.serializeOrRaw(data)
	if err != nil {
		return err
	}
	//写回客户端
	a.msgBuffChan <- pk
	return nil
}

func (a *Agent) serializeOrRaw(data pendingMessage) ([]byte, error) {
	var payload []byte
	var err error
	switch v := data.payload.(type) {
	case []byte:
		payload = v
	default:
		payload, err = a.serializer.Marshal(v)
		if err != nil {
			return nil, err
		}
	}

	m := &Message{
		Type:  data.typ,
		Data:  payload,
		Route: data.route,
		ID:    data.mid,
	}

	em, err := m.Encode()
	if err != nil {
		return nil, err
	}

	// packet encode
	p, err := Encode(Data, em)
	if err != nil {
		return nil, err
	}
	return p, nil
}
