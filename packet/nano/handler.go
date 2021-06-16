package nano

import (
	"encoding/json"
	"fmt"
	"github.com/jqiris/kungfu/component"
	"github.com/jqiris/kungfu/tcpface"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"time"

	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/packet"
)

var (
	logger = logrus.WithField("package", "zinx")
)

type MsgHandle struct {
	services       map[string]*component.Service // all registered service
	handlers       map[string]*component.Handler // all handler method
	WorkerPoolSize int                           //业务工作Worker池的数量
	TaskQueue      []chan *Request               //Worker负责取任务的消息队列
	// serialized data
	hrd []byte // handshake response data
	hbd []byte // heartbeat packet data
}

func NewMsgHandle() *MsgHandle {
	cfg := config.GetConnectorConf()
	h := &MsgHandle{
		Apis:           make(map[int32]Router),
		WorkerPoolSize: cfg.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan *Request, cfg.WorkerPoolSize),
	}
	h.hbdEncode()
	return h
}

func (h *MsgHandle) hbdEncode() {
	data, err := json.Marshal(map[string]interface{}{
		"code": 200,
		"sys":  map[string]float64{"heartbeat": 30},
	})
	if err != nil {
		panic(err)
	}

	h.hrd, err = Encode(Handshake, data)
	if err != nil {
		panic(err)
	}

	h.hbd, err = Encode(Heartbeat, nil)
	if err != nil {
		panic(err)
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (h *MsgHandle) SendMsgToTaskQueue(request *Request) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnID() % h.WorkerPoolSize
	//fmt.Println("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	h.TaskQueue[workerID] <- request
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (h *MsgHandle) DoMsgHandler(request *Request) {
	handler, ok := h.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler(request)
}

func (h *MsgHandle) register(comp component.Component, opts []component.Option) error {
	s := component.NewService(comp, opts)

	if _, ok := h.services[s.Name]; ok {
		return fmt.Errorf("handler: service already defined: %s", s.Name)
	}

	if err := s.ExtractHandler(); err != nil {
		return err
	}

	// register all handlers
	h.services[s.Name] = s
	for name, handler := range s.Handlers {
		h.handlers[fmt.Sprintf("%s.%s", s.Name, name)] = handler
	}
	return nil
}

// StartOneWorker 启动一个Worker工作流程
func (h *MsgHandle) StartOneWorker(workerID int, taskQueue chan *Request) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			h.DoMsgHandler(request)
		}
	}
}

// StartWorkerPool 启动worker工作池
func (h *MsgHandle) StartWorkerPool() {
	cfg := config.GetConnectorConf()
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(h.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		h.TaskQueue[i] = make(chan *Request, cfg.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go h.StartOneWorker(i, h.TaskQueue[i])
	}
}

func (h *MsgHandle) Handle(iConn tcpface.IConnection) {
	agent := iConn.(*Agent)
	go agent.StartWriter()
	defer func() {
		err := agent.Close()
		if err != nil {
			logger.Error(err)
		}
	}()
	conn := agent.GetConn()
	// read loop
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			logger.Println(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
			return
		}

		packets, err := agent.decoder.Decode(buf[:n])
		if err != nil {
			logger.Println(err.Error())
			return
		}

		if len(packets) < 1 {
			continue
		}

		// process all packet
		for i := range packets {
			if err := h.processPacket(agent, packets[i]); err != nil {
				logger.Println(err.Error())
				return
			}
		}
	}
}

func (h *MsgHandle) processPacket(agent *Agent, p *Packet) error {

	switch p.Type {
	case Handshake:
		if _, err := agent.conn.Write(h.hbd); err != nil {
			return err
		}

		agent.setStatus(packet.StatusHandshake)
		//if env.debug {
		//	logger.Println(fmt.Sprintf("Session handshake Id=%d, Remote=%s", agent.session.ID(), agent.conn.RemoteAddr()))
		//}

	case HandshakeAck:
		agent.setStatus(packet.StatusWorking)
		//if env.debug {
		//	logger.Println(fmt.Sprintf("Receive handshake ACK Id=%d, Remote=%s", agent.session.ID(), agent.conn.RemoteAddr()))
		//}

	case Data:
		if agent.status() < packet.StatusWorking {
			return fmt.Errorf("receive data on socket which not yet ACK, session will be closed immediately, remote=%s",
				agent.conn.RemoteAddr().String())
		}

		msg, err := MsgDecode(p.Data)
		if err != nil {
			return err
		}
		h.processMessage(agent, msg)

	case Heartbeat:
		// expected
	}

	agent.lastAt = time.Now().Unix()
	return nil
}

func (h *MsgHandle) processMessage(agent *Agent, msg *Message) {
	var lastMid uint
	switch msg.Type {
	case Request:
		lastMid = msg.ID
	case Notify:
		lastMid = 0
	}

	handler, ok := h.handlers[msg.Route]
	if !ok {
		logger.Println(fmt.Sprintf("handler: %s not found(forgot registered?)", msg.Route))
		return
	}
	var payload = msg.Data
	var data interface{}
	if handler.IsRawArg {
		data = payload
	} else {
		data = reflect.New(handler.Type.Elem()).Interface()
		err := serializer.Unmarshal(payload, data)
		if err != nil {
			logger.Println("deserialize error", err.Error())
			return
		}
	}

	if h.WorkerPoolSize > 0 {
		//已经启动工作池机制，将消息交给Worker处理
		h.SendMsgToTaskQueue(req)
	} else {
		//从绑定好的消息和对应的处理方法中执行对应的Handle方法
		go h.DoMsgHandler(req)
	}

}
