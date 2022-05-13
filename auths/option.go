package auths

import (
	"github.com/farmerx/gorsa"
	"github.com/jqiris/kungfu/v2/logger"
)

type Option func(e *Encipherer)

func WithRsaPubKey(pubKey string) Option {
	return func(e *Encipherer) {
		if err := gorsa.RSA.SetPublicKey(pubKey); err != nil {
			logger.Fatalf("set rsa public key err:%v", err)
		}
		e.rsaPubKey = pubKey
	}
}

func WithRsaPriKey(priKey string) Option {
	return func(e *Encipherer) {
		if err := gorsa.RSA.SetPrivateKey(priKey); err != nil {
			logger.Fatalf("set rsa pri key err:%v", err)
		}
		e.rsaPriKey = priKey
	}
}

func WithAesKey(aesKey string) Option {
	return func(e *Encipherer) {
		e.aesKey = []byte(aesKey)
	}
}

func WithAesIv(aesIv string) Option {
	return func(e *Encipherer) {
		e.aesIv = []byte(aesIv)
	}
}
