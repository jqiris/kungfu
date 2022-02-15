package tests

import (
	"context"
	"fmt"
	"github.com/jqiris/kungfu/v2/logger"
	"testing"
	"time"
)

var (
//client, cancel = NewRouteClient("client default")
)

//func SetClient(c *RouteClient) {
//	cancel()
//	client = c
//}

type RouteClient struct {
	Name string
}

func NewRouteClient(name string) (*RouteClient, context.CancelFunc) {
	c := &RouteClient{Name: name}
	ctx, cancel := context.WithCancel(context.Background())
	go c.RoutePrint(ctx)
	return c, cancel
}

func (r *RouteClient) RoutePrint(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println(r.Name, " is done")
			return
		case <-ticker.C:
			fmt.Println(time.Now(), "--", r.Name, ":ticked")
		}
	}
}

type RouteManager struct {
	routes map[string]*RouteClient
}

func TestRouteManager(t *testing.T) {
	manager := &RouteManager{routes: make(map[string]*RouteClient)}
	manager.routes["c1"], _ = NewRouteClient("c1")
	manager.routes["c2"], _ = NewRouteClient("c2")
	time.AfterFunc(5*time.Second, func() {
		delete(manager.routes, "c1")
		logger.Info("delete routers c1")
	})
	select {}
}

func TestRoutineRun(t *testing.T) {
	//c, tmp := NewRouteClient("client one")
	//SetClient(c)
	//cancel = tmp
	//select {}
}
