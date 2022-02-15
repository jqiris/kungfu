package rpc

import (
	"fmt"
	"github.com/jqiris/kungfu/v2/logger"
	"reflect"
)

type HandlerItem struct {
	MsgType MessageType
	InType  reflect.Type
	Func    reflect.Value
}
type Handler struct {
	handlers map[int32]HandlerItem
}

func NewHandler() *Handler {
	return &Handler{
		handlers: make(map[int32]HandlerItem),
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

func (h *Handler) DealMsg(codeType string, server ServerRpc, req *MsgRpc) ([]byte, error) {
	msgId, msgData := req.MsgId, req.MsgData.([]byte)
	if handler, ok := h.handlers[msgId]; ok {
		if handler.MsgType != req.MsgType {
			return nil, fmt.Errorf("req msg type not suit handler msg type, msgId:%v, req:%+v", msgId, req)
		}
		inElem := reflect.New(handler.InType.Elem()).Interface()
		err := server.DecodeMsg(codeType, msgData, inElem)
		if err != nil {
			return nil, err
		}
		args := []reflect.Value{reflect.ValueOf(inElem)}
		resp := handler.Func.Call(args)
		if handler.MsgType == Request && len(resp) > 0 {
			outItem := resp[0].Interface()
			return server.Response(codeType, outItem), nil
		}
		return nil, nil
	}
	return nil, fmt.Errorf("req msg not suit handler, msgId:%v, req:%+v", msgId, req)
}
