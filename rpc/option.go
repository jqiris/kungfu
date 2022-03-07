package rpc

type Option func(b *ServerBase)

func WithSelfEventHandler(handler CallbackFunc) Option {
	return func(b *ServerBase) {
		b.selfEventHandler = handler
	}
}

func WithBroadcastEventHandler(handler CallbackFunc) Option {
	return func(b *ServerBase) {
		b.broadcastEventHandler = handler
	}
}

func WithInnerMsgHandler(handler MsgHandler) Option {
	return func(b *ServerBase) {
		b.innerMsgHandler = handler
	}
}

func WithPlugin(plugin ServerPlugin) Option {
	return func(b *ServerBase) {
		b.plugins = append(b.plugins, plugin)
	}
}

func WithPlugins(plugins []ServerPlugin) Option {
	return func(b *ServerBase) {
		b.plugins = plugins
	}
}
