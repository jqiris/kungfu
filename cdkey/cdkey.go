package cdkey

import (
	"errors"
	"sync"
)

var (
	ErrReachMaxTimes = errors.New("reach max create times")
	ErrMinNumOne     = errors.New("min create num one")
	ErrInvalidCode   = errors.New("invalid code")
)

type CdkeyManager interface {
	MakeCode(relationId int64) string
	IsCodeExist(code string) bool
	CodeStore(relationId int64, code string) error
	IsCodeValid(code string) (bool, int64)
	CodeExchange(relationId int64, code string) error
}

type CdkeyProducer struct {
	mgr  CdkeyManager
	lock sync.Mutex
}

func NewCdkeyProducer(mgr CdkeyManager) *CdkeyProducer {
	return &CdkeyProducer{
		mgr: mgr,
	}
}

func (m *CdkeyProducer) GenCode(relationId int64, max, cur int) (string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	code := m.mgr.MakeCode(relationId)
	if m.mgr.IsCodeExist(code) {
		if max > 0 && cur >= max {
			return "", ErrReachMaxTimes
		}
		return m.GenCode(relationId, max, cur+1)
	}
	if err := m.mgr.CodeStore(relationId, code); err != nil {
		return "", err
	}
	return code, nil
}

func (m *CdkeyProducer) GenCodes(relationId int64, num, max int) ([]string, error) {
	if num < 1 {
		return nil, ErrMinNumOne
	}
	var codes []string
	for i := 0; i < num; i++ {
		code, err := m.GenCode(relationId, max, 0)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}

func (m *CdkeyProducer) ExchangeCode(code string) error {
	valid, relationId := m.mgr.IsCodeValid(code)
	if !valid {
		return ErrInvalidCode
	}
	return m.mgr.CodeExchange(relationId, code)
}
