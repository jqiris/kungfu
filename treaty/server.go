package treaty

import (
	"encoding/json"
)

func RegSeverItem(x *Server) string {
	return x.ServerType + "/" + x.ServerId
}

func RegServerType(x *Server) string {
	return x.ServerType
}

func RegSerialize(x *Server) string {
	if res, err := json.Marshal(x); err == nil {
		return string(res)
	}
	return ""
}

func RegUnSerialize(s string) (*Server, error) {
	x := &Server{}
	if err := json.Unmarshal([]byte(s), x); err != nil {
		return nil, err
	}
	return x, nil
}
