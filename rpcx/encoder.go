package rpcx

import "github.com/jqiris/kungfu/serialize"

const (
	NATS_ENCODER = "nats"
)

type NatsEncoder struct {
	useType string
	encoder serialize.Serializer
}

func NewNatsEncoder(useType string) *NatsEncoder {
	var encoder serialize.Serializer
	switch useType {
	case "json":
		encoder = serialize.NewJsonSerializer()
	case "proto":
		encoder = serialize.NewProtoSerializer()
	default:
		logger.Fatal("not support ")
	}
	return &NatsEncoder{
		useType: useType,
		encoder: encoder,
	}
}

func (n *NatsEncoder) Encode(subject string, v interface{}) ([]byte, error) {
	logger.Infof("nats encoder subject: %s", subject)
	return n.encoder.Marshal(v)
}

func (n *NatsEncoder) Decode(subject string, data []byte, vPtr interface{}) error {
	logger.Infof("nats decoder subject: %s", subject)
	return n.encoder.Unmarshal(data, vPtr)
}
