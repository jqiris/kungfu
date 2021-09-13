package rpcx

import (
	"fmt"
	"github.com/jqiris/kungfu/logger"
	"reflect"
	"time"
)

type HandlerItem struct {
	MsgType MessageType
	InType  reflect.Type
	Func    reflect.Value
}
type Handler struct {
	handlers     map[int32]HandlerItem
	maxPoolSize  int
	cSemaphore   chan struct{}
	reqPerSecond int
	rateLimiter  <-chan time.Time
}

func NewHandler(maxPoolSize int, reqPerSec int) *Handler {
	var semaphore chan struct{} = nil
	if maxPoolSize > 0 {
		semaphore = make(chan struct{}, maxPoolSize) // Buffered channel to act as a semaphore
	}

	var emitter <-chan time.Time = nil
	if reqPerSec > 0 {
		emitter = time.NewTicker(time.Second / time.Duration(reqPerSec)).C // x req/s == 1s/x req (inverse)
	}
	return &Handler{
		handlers:     make(map[int32]HandlerItem),
		maxPoolSize:  maxPoolSize,
		cSemaphore:   semaphore,
		reqPerSecond: reqPerSec,
		rateLimiter:  emitter,
	}
}

func (h *Handler) isSuitHandler(tf reflect.Type) bool {
	if tf.NumIn() != 1 {
		return false
	}
	if tf.In(0).Kind() != reflect.Ptr {
		return false
	}
	if tf.NumOut() > 1 || (tf.NumOut() == 1 && tf.Out(0).Kind() != reflect.Ptr) {
		return false
	}
	return true
}

func (h *Handler) Register(msgId int32, v interface{}) {
	if _, ok := h.handlers[msgId]; ok {
		logger.Errorf("msgId has already been registered:%v", msgId)
		return
	}
	vf, tf := reflect.ValueOf(v), reflect.TypeOf(v)
	if !h.isSuitHandler(tf) {
		logger.Errorf("not suit handler:%+v", v)
		return
	}
	msgType := MessageType(Publish)
	if tf.NumOut() == 1 {
		msgType = Request
	}
	h.handlers[msgId] = HandlerItem{
		MsgType: msgType,
		InType:  tf.In(0),
		Func:    vf,
	}
}

func (h *Handler) DealMsg(server RpcServer, req *RpcMsg) ([]byte, error) {
	if h.maxPoolSize > 0 {
		h.cSemaphore <- struct{}{} // Grab a connection from our pool
		defer func() {
			<-h.cSemaphore // Defer release our connection back to the pool
		}()
	}

	if h.reqPerSecond > 0 {
		<-h.rateLimiter // Block until a signal is emitted from the rateLimiter
	}
	msgId, msgData := req.MsgId, req.MsgData.([]byte)
	if handler, ok := h.handlers[msgId]; ok {
		if handler.MsgType != req.MsgType {
			return nil, fmt.Errorf("req msg type not suit handler msg type, msgId:%v, req:%+v", msgId, req)
		}
		inElem := reflect.New(handler.InType.Elem()).Interface()
		err := server.DecodeMsg(msgData, inElem)
		if err != nil {
			return nil, err
		}
		args := []reflect.Value{reflect.ValueOf(inElem)}
		resp := handler.Func.Call(args)
		if handler.MsgType == Request && len(resp) > 0 {
			outItem := resp[0].Interface()
			return server.Response(outItem), nil
		}
		return nil, nil
	}
	return nil, fmt.Errorf("req msg not suit handler, msgId:%v, req:%+v", msgId, req)
}

func (h *Handler) DealJsonMsg(server RpcServer, req *RpcMsg) ([]byte, error) {
	if h.maxPoolSize > 0 {
		h.cSemaphore <- struct{}{} // Grab a connection from our pool
		defer func() {
			<-h.cSemaphore // Defer release our connection back to the pool
		}()
	}

	if h.reqPerSecond > 0 {
		<-h.rateLimiter // Block until a signal is emitted from the rateLimiter
	}
	msgId, msgData := req.MsgId, req.MsgData.([]byte)
	if handler, ok := h.handlers[msgId]; ok {
		if handler.MsgType != req.MsgType {
			return nil, fmt.Errorf("req msg type not suit handler msg type, msgId:%v, req:%+v", msgId, req)
		}
		inElem := reflect.New(handler.InType.Elem()).Interface()
		err := server.DecodeJsonMsg(msgData, inElem)
		if err != nil {
			return nil, err
		}
		args := []reflect.Value{reflect.ValueOf(inElem)}
		resp := handler.Func.Call(args)
		if handler.MsgType == Request && len(resp) > 0 {
			outItem := resp[0].Interface()
			return server.ResponseJson(outItem), nil
		}
		return nil, nil
	}
	return nil, fmt.Errorf("req msg not suit handler, msgId:%v, req:%+v", msgId, req)
}
