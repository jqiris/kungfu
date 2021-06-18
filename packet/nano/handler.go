package nano

import (
	"encoding/json"
	"fmt"
	"github.com/jqiris/kungfu/component"
	"github.com/jqiris/kungfu/serialize"
	"github.com/jqiris/kungfu/tcpface"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/packet"
)

var (
	logger = logrus.WithField("package", "nano")
)

type MsgHandle struct {
	services       map[string]*component.Service // all registered service
	handlers       map[string]*component.Handler // all handler method
	WorkerPoolSize int                           //业务工作Worker池的数量
	Serializer     serialize.Serializer          //序列化对象
	TaskQueue      []chan unhandledMessage       //Worker负责取任务的消息队列
	// serialized data
	hrd []byte // handshake response data
	hbd []byte // heartbeat packet data
}
type unhandledMessage struct {
	agent   *Agent
	lastMid uint
	handler reflect.Method
	args    []reflect.Value
}

func NewMsgHandle() *MsgHandle {
	cfg := config.GetConnectorConf()
	h := &MsgHandle{
		services:       make(map[string]*component.Service),
		handlers:       make(map[string]*component.Handler),
		WorkerPoolSize: cfg.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan unhandledMessage, cfg.WorkerPoolSize),
	}
	h.hbdEncode()
	switch cfg.UseSerializer {
	case "proto":
		h.Serializer = serialize.NewProtoSerializer()
	case "json":
		h.Serializer = serialize.NewJsonSerializer()
	default:
		logger.Fatalf("no suitable serializer:%v", cfg.UseSerializer)
	}
	return h
}

func (h *MsgHandle) hbdEncode() {
	data, err := json.Marshal(map[string]interface{}{
		"code": 200,
		"sys": map[string]interface{}{
			"heartbeat": 30,
			"protos": map[string]interface{}{
				"client": map[string]interface{}{
					"UserConnector.Login": map[string]interface{}{
						"required uInt32 uid":      1,
						"required string nickname": 2,
						"required string token":    3,
						"message Server": map[string]interface{}{
							"required string server_id":   1,
							"required string server_type": 2,
							"required string server_name": 3,
							"required string server_ip":   4,
							"required uInt32 client_port": 5,
						},
						"required Server backend": 4,
					},
				},
				"server": map[string]interface{}{
					"UserConnector.Login": map[string]interface{}{
						"required uInt32 code": 1,
						"required string msg":  2,
						"message Server": map[string]interface{}{
							"required string server_id":   1,
							"required string server_type": 2,
							"required string server_name": 3,
							"required string server_ip":   4,
							"required uInt32 client_port": 5,
						},
						"required Server backend": 3,
					},
				},
			},
		},
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
func (h *MsgHandle) SendMsgToTaskQueue(request unhandledMessage) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.agent.connId % h.WorkerPoolSize
	//fmt.Println("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	h.TaskQueue[workerID] <- request
}
func stack() string {
	buf := make([]byte, 10000)
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	s := string(buf)

	// skip nano frames lines
	const skip = 7
	count := 0
	index := strings.IndexFunc(s, func(c rune) bool {
		if c != '\n' {
			return false
		}
		count++
		return count == skip
	})
	return s[index+1:]
}

// call handler with protected
func pCall(method reflect.Method, args []reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("dispatch: %v", err)
			println(stack())
		}
	}()

	if r := method.Func.Call(args); len(r) > 0 {
		if err := r[0].Interface(); err != nil {
			logger.Errorf(err.(error).Error())
		}
	}
}

// call handler with protected
func pinvoke(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println(fmt.Sprintf("nano/invoke: %v", err))
			println(stack())
		}
	}()

	fn()
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (h *MsgHandle) DoMsgHandler(request unhandledMessage) {
	request.agent.lastMid = request.lastMid
	pCall(request.handler, request.args)
}

func (h *MsgHandle) Register(comp component.Component, opts ...component.Option) error {
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
	h.DumpServices()
	return nil
}

// StartOneWorker 启动一个Worker工作流程
func (h *MsgHandle) StartOneWorker(workerID int, taskQueue chan unhandledMessage) {
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
		h.TaskQueue[i] = make(chan unhandledMessage, cfg.MaxWorkerTaskLen)
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
		if err := agent.SendRawMessage(true, h.hrd); err != nil {
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
		err := h.Serializer.Unmarshal(payload, data)
		if err != nil {
			logger.Println("deserialize error", err.Error())
			return
		}
	}
	args := []reflect.Value{handler.Receiver, agent.srv, reflect.ValueOf(data)}
	request := unhandledMessage{agent, lastMid, handler.Method, args}
	if h.WorkerPoolSize > 0 {
		//已经启动工作池机制，将消息交给Worker处理
		h.SendMsgToTaskQueue(request)
	} else {
		//从绑定好的消息和对应的处理方法中执行对应的Handle方法
		go h.DoMsgHandler(request)
	}
}

// DumpServices outputs all registered services
func (h *MsgHandle) DumpServices() {
	logger.Println("DumpServices:")
	for name := range h.handlers {
		logger.Println("registered service：", name)
	}
}
