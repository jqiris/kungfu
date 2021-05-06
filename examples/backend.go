package examples

import (
	"fmt"

	"github.com/jqiris/kungfu/backend"
	"github.com/jqiris/kungfu/launch"
)

type MyBackend struct {
	backend.BaseBackEnd
}

func (b *MyBackend) EventHandleSelf(req []byte) []byte {
	fmt.Printf("MyBackend EventHandleSelf received: %+v \n", string(req))
	return nil
}

func (b *MyBackend) EventHandleBroadcast(req []byte) []byte {
	fmt.Printf("MyBackend EventHandleBroadcast received: %+v \n", string(req))
	return nil
}

func init() {
	srv := &MyBackend{}
	srv.SetServerId("backend_3001")
	srv.RegEventHandlerSelf(srv.EventHandlerSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroadcast)
	launch.RegisterServer(srv)
}
