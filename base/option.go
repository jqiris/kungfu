package base

import "github.com/jqiris/kungfu/v2/rpc"

type Option func(b *ServerBase)

func WithSelfEventHandler(handler rpc.CallbackFunc) Option {
	return func(b *ServerBase) {
		b.selfEventHandler = handler
	}
}

func WithBroadcastEventHandler(handler rpc.CallbackFunc) Option {
	return func(b *ServerBase) {
		b.broadcastEventHandler = handler
	}
}

func WithInnerMsgHandler(handler rpc.MsgHandler) Option {
	return func(b *ServerBase) {
		b.innerMsgHandler = handler
	}
}

func WithPlugin(plugin rpc.ServerPlugin) Option {
	return func(b *ServerBase) {
		b.plugins = append(b.plugins, plugin)
	}
}

func WithPlugins(plugins []rpc.ServerPlugin) Option {
	return func(b *ServerBase) {
		b.plugins = plugins
	}
}
