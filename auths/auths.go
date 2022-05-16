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

func (e *Encipherer) RsaPrikeyEncrypt(data []byte) (string, error) {
	rsaData, err := gorsa.RSA.PriKeyENCTYPT(data)
	if err != nil {
		return "", err
	}
	bsData := base64.StdEncoding.EncodeToString(rsaData)
	return bsData, nil
}

func (e *Encipherer) RsaPrikeyDecrypt(data []byte) ([]byte, error) {
	return gorsa.RSA.PriKeyDECRYPT(data)
}

func (e *Encipherer) RsaPubkeyEncrypt(data []byte) (string, error) {
	rsaData, err := gorsa.RSA.PubKeyENCTYPT(data)
	if err != nil {
		return "", err
	}
	bsData := base64.StdEncoding.EncodeToString(rsaData)
	return bsData, nil
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

func (e *Encipherer) AesCbcEncrypt(data []byte) (string, error) {
	aesData, err := goEncrypt.AesCbcEncrypt(data, e.aesKey, e.aesIv)
	if err != nil {
		return "", err
	}
	bsData := base64.StdEncoding.EncodeToString(aesData)
	return bsData, nil
}

func (e *Encipherer) GetAesSecretKey() (string, string) {
	return string(e.aesKey), string(e.aesIv)
}
