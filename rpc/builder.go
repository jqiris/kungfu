package rpc

import (
	"time"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
)

func DefaultCallback(req *MsgRpc) []byte {
	logger.Info("DefaultCallback")
	return nil
}

type RssBuilder struct {
	queue       string
	server      *treaty.Server
	callback    CallbackFunc
	codeType    string
	suffix      string
	parallel    bool
	dialTimeout time.Duration
	//for rabbitmq
	exName string //交换机名称
	exType string //交换类型
	rtKey  string //绑定key

}

func NewRssBuilder(server *treaty.Server) *RssBuilder {
	parallel := true
	if server.Serial { //串行处理
		parallel = false
	}
	return &RssBuilder{
		queue:    DefaultQueue,
		server:   server,
		callback: DefaultCallback,
		codeType: CodeTypeProto,
		suffix:   DefaultSuffix,
		parallel: parallel,
		exName:   DefaultExName,
		exType:   DefaultExType,
		rtKey:    DefaultRtKey,
	}
}
func (r *RssBuilder) SetDialTimeout(d time.Duration) *RssBuilder {
	r.dialTimeout = d
	return r
}
func (r *RssBuilder) SetExName(name string) *RssBuilder {
	r.exName = name
	return r
}

func (r *RssBuilder) SetExType(typ string) *RssBuilder {
	r.exType = typ
	return r
}

func (r *RssBuilder) SetRtKey(rtKey string) *RssBuilder {
	r.rtKey = rtKey
	return r
}

func (r *RssBuilder) SetQueue(queue string) *RssBuilder {
	r.queue = queue
	return r
}

func (r *RssBuilder) SetServer(server *treaty.Server) *RssBuilder {
	r.server = server
	return r
}
func (r *RssBuilder) SetCallback(callback CallbackFunc) *RssBuilder {
	r.callback = callback
	return r
}
func (r *RssBuilder) SetCodeType(codeType string) *RssBuilder {
	r.codeType = codeType
	return r
}
func (r *RssBuilder) SetSuffix(suffix string) *RssBuilder {
	r.suffix = suffix
	return r
}
func (r *RssBuilder) SetParallel(parallel bool) *RssBuilder {
	r.parallel = parallel
	return r
}

func (r *RssBuilder) Build() RssBuilder {
	return RssBuilder{
		queue:    r.queue,
		server:   r.server,
		callback: r.callback,
		codeType: r.codeType,
		suffix:   r.suffix,
		parallel: r.parallel,
		exName:   r.exName,
		exType:   r.exType,
		rtKey:    r.rtKey,
	}
}

type ReqBuilder struct {
	queue       string
	codeType    string
	suffix      string
	server      *treaty.Server
	msgId       int32
	req         any
	resp        any
	serverType  string
	dialTimeout time.Duration
	//for rabbitmq
	exName string //交换机名称
	exType string //交换类型
	rtKey  string //绑定key
}

func NewReqBuilder(server *treaty.Server) *ReqBuilder {
	serverType := ""
	if server != nil {
		serverType = server.ServerType
	}
	return &ReqBuilder{
		queue:      DefaultQueue,
		codeType:   CodeTypeProto,
		suffix:     DefaultSuffix,
		server:     server,
		serverType: serverType,
		exName:     DefaultExName,
		exType:     DefaultExType,
		rtKey:      DefaultRtKey,
	}
}

func DefaultReqBuilder() *ReqBuilder {
	return &ReqBuilder{
		queue:    DefaultQueue,
		codeType: CodeTypeProto,
		suffix:   DefaultSuffix,
		exName:   DefaultExName,
		exType:   DefaultExType,
		rtKey:    DefaultRtKey,
	}
}

func (r *ReqBuilder) SetExName(name string) *ReqBuilder {
	r.exName = name
	return r
}

func (r *ReqBuilder) SetExType(typ string) *ReqBuilder {
	r.exType = typ
	return r
}

func (r *ReqBuilder) SetRtKey(rtKey string) *ReqBuilder {
	r.rtKey = rtKey
	return r
}

func (r *ReqBuilder) SetQueue(queue string) *ReqBuilder {
	r.queue = queue
	return r
}
func (r *ReqBuilder) SetCodeType(codeType string) *ReqBuilder {
	r.codeType = codeType
	return r
}
func (r *ReqBuilder) SetSuffix(suffix string) *ReqBuilder {
	r.suffix = suffix
	return r
}
func (r *ReqBuilder) SetServer(server *treaty.Server) *ReqBuilder {
	r.server = server
	r.SetServerType(server.ServerType)
	return r
}
func (r *ReqBuilder) SetMsgId(msgId int32) *ReqBuilder {
	r.msgId = msgId
	return r
}
func (r *ReqBuilder) SetReq(req any) *ReqBuilder {
	r.req = req
	return r
}
func (r *ReqBuilder) SetResp(resp any) *ReqBuilder {
	r.resp = resp
	return r
}
func (r *ReqBuilder) SetServerType(serverType string) *ReqBuilder {
	r.serverType = serverType
	return r
}
func (r *ReqBuilder) SetDialTimeout(d time.Duration) *ReqBuilder {
	r.dialTimeout = d
	return r
}

func (r *ReqBuilder) Build() ReqBuilder {
	return ReqBuilder{
		queue:       r.queue,
		codeType:    r.codeType,
		suffix:      r.suffix,
		server:      r.server,
		msgId:       r.msgId,
		req:         r.req,
		resp:        r.resp,
		serverType:  r.serverType,
		dialTimeout: r.dialTimeout,
		exName:      r.exName,
		exType:      r.exType,
		rtKey:       r.rtKey,
	}
}
