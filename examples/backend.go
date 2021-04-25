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

func (b *MyBackend) EventHandleBroasdcast(req []byte) []byte {
	fmt.Printf("MyBackend EventHandleBroasdcast received: %+v \n", string(req))
	return nil
}

func init() {
	srv := &MyBackend{}
	srv.SetServerId(3001)
	srv.RegEventHandlerSelf(srv.EventHandlerSelf)
	srv.RegEventHandlerBroadcast(srv.EventHandleBroasdcast)
	launch.RegisterServer(srv)
}
