package tests

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var (
	client, cancel = NewRouteClient("client default")
)

func SetClient(c *RouteClient) {
	cancel()
	client = c
}

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

func TestRoutineRun(t *testing.T) {
	c, tmp := NewRouteClient("client one")
	SetClient(c)
	cancel = tmp
	select {}
}
