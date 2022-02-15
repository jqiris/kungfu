package base

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/golang/protobuf/proto"
	"github.com/jqiris/kungfu/v2/discover"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/rpc"
	"github.com/jqiris/kungfu/v2/serialize"
	"github.com/jqiris/kungfu/v2/session"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
	"net/http"
	"net/url"
)

type ServerBalancer struct {
	*ServerBase
	ClientServer *http.Server
	ClientCoder  serialize.Serializer
}

func NewServerBalancer(s *treaty.Server) *ServerBalancer {
	return &ServerBalancer{
		ServerBase: NewServerBase(s),
	}
}

func (s *ServerBalancer) Init() {
	s.ServerBase.Init()
	//init the coder
	s.ClientCoder = serialize.NewProtoSerializer()
	//set the server
	s.ClientServer = &http.Server{Addr: fmt.Sprintf("%s:%d", s.Server.ServerIp, s.Server.ClientPort)}
	//handle the balance
	http.HandleFunc("/balance", s.HandleBalance)
	//run the server
	go func() {
		err := s.ClientServer.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
		}
	}()
}
func (s *ServerBalancer) Shutdown() {
	s.ServerBase.Shutdown()
	if s.ClientServer != nil {
		if err := s.ClientServer.Close(); err != nil {
			logger.Error(err)
		}
	}
}

func (s *ServerBalancer) HandleBalance(w http.ResponseWriter, r *http.Request) {
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
		s.WriteResponse(w, res)
		return
	}
	connector, err := s.Balance(r.RemoteAddr)
	if err != nil {
		res := &treaty.BalanceResult{
			Code: treaty.CodeType_CodeFailed,
		}
		s.WriteResponse(w, res)
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
	s.WriteResponse(w, res)
}

func (s *ServerBalancer) WriteResponse(w http.ResponseWriter, msg proto.Message) {
	if v, e := s.ClientCoder.Marshal(msg); e == nil {
		if _, e2 := w.Write(v); e2 != nil {
			logger.Error(e2)
		}
	}
}
func (s *ServerBalancer) Balance(remoteAddr string) (*treaty.Server, error) {
	if server := discover.GetServerByType(rpc.Connector, remoteAddr); server != nil {
		return server, nil
	}
	return nil, errors.New("no suitable connector found")
}
