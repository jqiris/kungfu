package treaty

import (
	"encoding/json"
	"fmt"
)

func (x *Server) RegId() string {
	return fmt.Sprintf("%d/%d", x.ServerType, x.ServerId)
}

func (x *Server) RegType() string {
	return fmt.Sprintf("%d", x.ServerType)
}

func (x *Server) Serialize() string{
	if res, err := json.Marshal(x);err == nil{
		return string(res)
	}
	return ""
}
