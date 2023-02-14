/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package zinx

import (
	"fmt"
	"github.com/jqiris/kungfu/v2/tcpface"
	"strconv"
	"time"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/packet"
)

type MsgHandle struct {
	Apis           map[int32]Router //存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize int              //业务工作Worker池的数量
	TaskQueue      []chan *Request  //Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	cfg := config.GetConnectorConf()
	workerPoolSize := 10
	if cfg.WorkerPoolSize > 0 {
		workerPoolSize = cfg.WorkerPoolSize
	}
	return &MsgHandle{
		Apis:           make(map[int32]Router),
		WorkerPoolSize: workerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan *Request, workerPoolSize),
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (h *MsgHandle) SendMsgToTaskQueue(request *Request) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnID() % h.WorkerPoolSize
	//logger.Info("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	h.TaskQueue[workerID] <- request
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (h *MsgHandle) DoMsgHandler(request *Request) {
	handler, ok := h.Apis[request.GetMsgID()]
	if !ok {
		logger.Error("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (h *MsgHandle) AddRouter(msgId int32, router Router) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := h.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	h.Apis[msgId] = router
	logger.Info("Add api msgId = ", msgId)
}

// StartOneWorker 启动一个Worker工作流程
func (h *MsgHandle) StartOneWorker(workerID int, taskQueue chan *Request) {
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
	var maxWorkerTaskLen int32 = 1024
	if cfg.MaxMsgChanLen > 0 {
		maxWorkerTaskLen = cfg.MaxWorkerTaskLen
	}
	logger.Infof("start worker pool:%v， one pool size:%v", h.WorkerPoolSize, maxWorkerTaskLen)
	for i := 0; i < int(h.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		h.TaskQueue[i] = make(chan *Request, maxWorkerTaskLen)
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
			logger.Info(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
			return
		}

		packets, err := agent.decoder.Decode(buf[:n])
		if err != nil {
			logger.Info(err.Error())
			return
		}

		if len(packets) < 1 {
			continue
		}

		// process all packet
		for i := range packets {
			if err := h.processPacket(agent, packets[i]); err != nil {
				logger.Info(err.Error())
				return
			}
		}
	}
}

func (h *MsgHandle) processPacket(agent *Agent, p *Packet) error {

	if agent.status() < packet.StatusWorking {
		return fmt.Errorf("receive data on socket which not yet ACK, session will be closed immediately, remote=%s",
			agent.conn.RemoteAddr().String())
	}

	msg, err := MsgDecode(p.Data)
	if err != nil {
		return err
	}
	h.processMessage(agent, msg)

	agent.lastAt = time.Now().Unix()
	return nil
}

func (h *MsgHandle) processMessage(agent *Agent, msg *Message) {
	req := &Request{
		agent: agent,
		msg:   msg,
	}

	if h.WorkerPoolSize > 0 {
		//已经启动工作池机制，将消息交给Worker处理
		h.SendMsgToTaskQueue(req)
	} else {
		//从绑定好的消息和对应的处理方法中执行对应的Handle方法
		go h.DoMsgHandler(req)
	}

}
