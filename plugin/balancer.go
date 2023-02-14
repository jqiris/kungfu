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

package plugin

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/golang/protobuf/proto"
	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/serialize"
	"github.com/jqiris/kungfu/v2/session"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
)

type ServerBalancer struct {
	ClientServer *http.Server
	ClientCoder  serialize.Serializer
}

func NewServerBalancer() *ServerBalancer {
	return &ServerBalancer{
		ClientCoder: serialize.NewProtoSerializer(),
	}
}

func (b *ServerBalancer) Init(s *rpc.ServerBase) {

	//run the server
	go utils.SafeRun(func() {
		b.Run(s)
	})
}
func (b *ServerBalancer) Run(s *rpc.ServerBase) {
	//set the server
	addr := fmt.Sprintf("%s:%d", s.Server.ServerIp, s.Server.ClientPort)
	b.ClientServer = &http.Server{Addr: addr}
	//handle the balance
	http.HandleFunc("/balance", b.HandleBalance)
	logger.Warnf("ServerBalancer start at:%v", addr)
	err := b.ClientServer.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
	}
}

func (b *ServerBalancer) AfterInit(s *rpc.ServerBase) {

}

func (b *ServerBalancer) BeforeShutdown(s *rpc.ServerBase) {
}

func (b *ServerBalancer) Shutdown(s *rpc.ServerBase) {
	if b.ClientServer != nil {
		if err := b.ClientServer.Close(); err != nil {
			logger.Error(err)
		}
	}
}

func (b *ServerBalancer) HandleBalance(w http.ResponseWriter, r *http.Request) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	serverType, uid := "", 0
	if err == nil {
		if len(queryForm["server_type"]) > 0 {
			serverType = queryForm["server_type"][0]
		}
		if len(queryForm["uid"]) > 0 {
			uid = utils.StringToInt(queryForm["uid"][0])
		}
	}
	if len(serverType) < 1 {
		res := &treaty.BalanceResult{
			Code: treaty.CodeType_CodeChooseBackendLogin,
		}
		b.WriteResponse(w, res)
		return
	}
	connector, err := b.Balance(r.RemoteAddr)
	if err != nil {
		res := &treaty.BalanceResult{
			Code: treaty.CodeType_CodeFailed,
		}
		b.WriteResponse(w, res)
		return
	}
	backend := discover.GetServerByType(serverType, r.RemoteAddr)
	var backendPre *treaty.Server
	sess := session.GetSession(int32(uid))
	if sess != nil {
		backendPre = sess.Backend
	}
	res := &treaty.BalanceResult{
		Code:       treaty.CodeType_CodeSuccess,
		Connector:  connector,
		Backend:    backend,
		BackendPre: backendPre,
	}
	b.WriteResponse(w, res)
}

func (b *ServerBalancer) WriteResponse(w http.ResponseWriter, msg proto.Message) {
	if v, e := b.ClientCoder.Marshal(msg); e == nil {
		if _, e2 := w.Write(v); e2 != nil {
			logger.Error(e2)
		}
	}
}
func (b *ServerBalancer) Balance(remoteAddr string) (*treaty.Server, error) {
	if server := discover.GetServerByType(rpc.Connector, remoteAddr); server != nil {
		return server, nil
	}
	return nil, errors.New("no suitable connector found")
}
