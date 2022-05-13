package auths

import (
	"encoding/base64"

	"github.com/farmerx/gorsa"
	"github.com/wumansgy/goEncrypt"
)

type Encipherer struct {
	rsaPriKey string
	rsaPubKey string
	aesKey    []byte
	aesIv     []byte
}

func NewEncipherer(options ...Option) *Encipherer {
	encipherer := &Encipherer{}
	for _, option := range options {
		option(encipherer)
	}
	return encipherer
}

func (e *Encipherer) RsaPrikeyEncrypt(data []byte) ([]byte, error) {
	return gorsa.RSA.PriKeyENCTYPT(data)
}

func (e *Encipherer) RsaPrikeyDecrypt(data []byte) ([]byte, error) {
	return gorsa.RSA.PriKeyDECRYPT(data)
}

func (e *Encipherer) RsaPubkeyEncrypt(data []byte) ([]byte, error) {
	return gorsa.RSA.PubKeyENCTYPT(data)
}

func (e *Encipherer) RsaPubkeyDecrypt(data []byte) ([]byte, error) {
	return gorsa.RSA.PubKeyDECRYPT(data)
}

func (e *Encipherer) AesCbcDecrypt(data []byte) ([]byte, error) {
	bs, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return goEncrypt.AesCbcDecrypt(bs, e.aesKey, e.aesIv)
}

func (e *Encipherer) AesCbcEncrypt(data []byte) ([]byte, error) {
	aesData, err := goEncrypt.AesCbcEncrypt(data, e.aesKey, e.aesIv)
	if err != nil {
		return nil, err
	}
	bsData := base64.StdEncoding.EncodeToString(aesData)
	return []byte(bsData), nil
}

func (e *Encipherer) GetAesSecretKey() (string, string) {
	return string(e.aesKey), string(e.aesIv)
}
