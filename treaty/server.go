package treaty

import "fmt"

func (x *Server) RegId() string {
	return fmt.Sprintf("%d/%d", x.ServerType, x.ServerId)
}

func (x *Server) RegType() string {
	return fmt.Sprintf("%d", x.ServerType)
}
