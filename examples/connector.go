package examples

import (
	"fmt"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/connector"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/launch"
	"github.com/jqiris/kungfu/session"
	"github.com/jqiris/kungfu/tcpserver"
	"github.com/jqiris/kungfu/treaty"
)

type MyConnector struct {
	connector.BaseConnector

}

func (b *MyConnector) EventHandleSelf(req []byte) []byte {
	fmt.Printf("MyConnector EventHandleSelf received: %+v \n", string(req))
	return nil
}

func (b *MyConnector) EventHandleBroadcast(req []byte) []byte {
	fmt.Printf("MyConnector EventHandleBroadcast received: %+v \n", string(req))
	return nil
}
//Login 登录操作
func (b *MyConnector) Login(request tcpserver.IRequest) {

	//先读取客户端的数据
	logger.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//回复信息
	resp := &treaty.LoginResponse{}
	//回复对象
	conn := request.GetConnection()
	//解析登录数据
	req := &treaty.LoginRequest{}
	if err := GetRequest(request, req); err != nil {
		resp.Code = 1
		resp.Msg = err.Error()
		SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
		return
	}
	logger.Printf("login request is:%+v", req)
	//判断登录信息的正确性
	uid, nickname := req.Uid, req.Nickname
	tokenkey := conf.GetConnectorConf().TokenKey
	token := helper.Md5(fmt.Sprintf("%d|%s|%s", uid, nickname, tokenkey))
	if req.Token != token {
		resp.Code = 1
		resp.Msg = "token不正确"
		SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
		return
	}
	//检查会话信息
	sess := session.GetSession(uid)
	if sess != nil {
		if sess.Connector != nil && sess.Connector.ServerId != request.GetServerID() {
			//之前在其他客户端登录，通知其他connetor登出
		}
		//如果连接了后端服务器，进行重连
		if sess.Backend != nil && sess.Backend != req.Backend{
			resp.Code = 2
			resp.Msg = "请进行重连服务器"
			resp.Backend = sess.Backend
			return
		}
	}

	if req.Backend != nil{
		//与后端服务器建立连接
		backend := discover.GetServerById(req.Backend.ServerId)
		if backend == nil{
			//查找同类服务器
			backend  = discover.GetServerByType(req.Backend.ServerType,conn.RemoteAddr().String())
		}
		if backend == nil{
			resp.Code = 3
			resp.Msg = "找不到服务器"
			return
		}
	} else {
		//加入大厅用户

	}
	resp.Code = 0
	resp.Msg = "登录成功"
	SendMsg(conn, treaty.MsgId_Msg_Login_Response, resp)
}

func init() {
	srv := &MyConnector{}
	routers := map[int32]tcpserver.IHandler{
		int32(treaty.MsgId_Msg_Login_Request): srv.Login,
	}
	srv.SetServerId("connector_2001")
	srv.RegEventHandlerSelf(srv.EventHandlerSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	srv.RegRouters(routers)
	launch.RegisterServer(srv)

	srv2 := &MyConnector{}
	srv2.SetServerId("connector_2002")
	srv2.RegEventHandlerSelf(srv2.EventHandlerSelf)
	srv2.RegEventHandlerBroadcast(srv2.EventHandleBroadcast)
	srv2.RegRouters(routers)
	launch.RegisterServer(srv2)
}
