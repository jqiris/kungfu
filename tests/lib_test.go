package tests

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/nats-io/nats.go"
	"testing"
	"time"
)

func TestLogs(t *testing.T) {
	ctx := log.WithFields(log.Fields{
		"file": "something.png",
		"type": "image/png",
		"user": "tobi",
	})

	for range time.Tick(time.Millisecond * 200) {
		ctx.Info("upload")
		ctx.Info("upload complete")
		ctx.Warn("upload retry")
		ctx.WithError(errors.New("unauthorized")).Error("upload failed")
		ctx.Errorf("failed to upload %s", "img.png")
	}
}

func TestNats(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.Fatal(err)
	}
	// Simple Publisher
	nc.Publish("foo", []byte("Hello World"))

	// Simple Async Subscriber
	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	// Responding to a request message
	nc.Subscribe("request", func(m *nats.Msg) {
		m.Respond([]byte("answer is 42"))
	})

	// Simple Sync Subscriber
	//sub, err := nc.SubscribeSync("foo")
	//if err != nil {
	//	logger.Error(err)
	//}
	//timeout := 10 * time.Second
	//m, err := sub.NextMsg(timeout)
	//logger.Debug(m)
	// Channel Subscriber
	//ch := make(chan *nats.Msg, 64)
	//sub, err := nc.ChanSubscribe("foo", ch)
	//msg := <-ch
	//logger.Debug(msg)
	//// Unsubscribe
	//sub.Unsubscribe()

	// Drain
	//sub.Drain()
	nc.Subscribe("help", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("I can help!"))
	})
	// Requests
	msg, err := nc.Request("help", []byte("help me"), 10*time.Millisecond)
	logger.Debug(msg, err)
	// Replies

	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	nc.Drain()

	// Close connection
	nc.Close()
}
