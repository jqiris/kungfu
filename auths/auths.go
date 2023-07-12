/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package auths

import (
	"encoding/base64"
	"errors"

	"github.com/farmerx/gorsa"
	"github.com/jqiris/kungfu/v2/serialize"
	"github.com/wumansgy/goEncrypt"
)

type Encipherer struct {
	rsaPriKey  string
	rsaPubKey  string
	aesKey     []byte
	aesIv      []byte
	serializer serialize.Serializer
	unencrypt  bool
}

func NewEncipherer(options ...Option) *Encipherer {
	encipherer := &Encipherer{
		serializer: serialize.NewJsonSerializer(),
		unencrypt:  false,
	}
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

func (e *Encipherer) RsaPrikeyDecrypt(data string) ([]byte, error) {
	bs, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return gorsa.RSA.PriKeyDECRYPT(bs)
}

func (e *Encipherer) RsaPubkeyEncrypt(data []byte) (string, error) {
	rsaData, err := gorsa.RSA.PubKeyENCTYPT(data)
	if err != nil {
		return "", err
	}
	bsData := base64.StdEncoding.EncodeToString(rsaData)
	return bsData, nil
}

func (e *Encipherer) RsaPubkeyDecrypt(data string) ([]byte, error) {
	bs, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return gorsa.RSA.PubKeyDECRYPT(bs)
}

func (e *Encipherer) AesCbcDecrypt(data string) ([]byte, error) {
	bs, err := base64.StdEncoding.DecodeString(data)
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

type DecryptType int

const (
	_                    DecryptType = iota
	DecryptTypeAes                   //1-aes解密
	DecryptTypeRsaPrikey             //2-rsa私钥解密
	DecryptTypeRsaPubkey             //2-rsa公钥解密
	DecryptTypeNone                  //2-none
)

type EncryptType int

const (
	_                    EncryptType = iota
	EncryptTypeAes                   //1-aes加密
	EncryptTypeRsaPrikey             //2-rsa私钥加密
	EncryptTypeRsaPubkey             //2-rsag公钥加密
	EncryptTypeNone                  //-none
)

func DecryptData(ep *Encipherer, typ DecryptType, src string, res any) error {
	if ep.unencrypt {
		typ = DecryptTypeNone
	}
	var data []byte
	var err error
	switch typ {
	case DecryptTypeAes:
		data, err = ep.AesCbcDecrypt(src)
		if err != nil {
			return err
		}
	case DecryptTypeRsaPrikey:
		data, err = ep.RsaPrikeyDecrypt(src)
		if err != nil {
			return err
		}
	case DecryptTypeRsaPubkey:
		data, err = ep.RsaPubkeyDecrypt(src)
		if err != nil {
			return err
		}
	case DecryptTypeNone:
		data = []byte(src)
	}
	err = ep.serializer.Unmarshal(data, res)
	if err != nil {
		return err
	}
	return nil
}

func EncryptData(ep *Encipherer, typ EncryptType, src any) (string, error) {
	if ep.unencrypt {
		typ = EncryptTypeNone
	}
	data, err := ep.serializer.Marshal(src)
	if err != nil {
		return "", err
	}
	switch typ {
	case EncryptTypeAes:
		return ep.AesCbcEncrypt(data)
	case EncryptTypeRsaPrikey:
		return ep.RsaPrikeyEncrypt(data)
	case EncryptTypeRsaPubkey:
		return ep.RsaPubkeyEncrypt(data)
	case EncryptTypeNone:
		return string(data), nil
	}
	return "", errors.New("no suit type")
}

func EncryptDataRaw(ep *Encipherer, typ EncryptType, data []byte) (string, error) {
	if ep.unencrypt {
		typ = EncryptTypeNone
	}
	switch typ {
	case EncryptTypeAes:
		return ep.AesCbcEncrypt(data)
	case EncryptTypeRsaPrikey:
		return ep.RsaPrikeyEncrypt(data)
	case EncryptTypeRsaPubkey:
		return ep.RsaPubkeyEncrypt(data)
	case EncryptTypeNone:
		return string(data), nil
	}
	return "", errors.New("no suit type")
}
