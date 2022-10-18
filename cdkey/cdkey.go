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
	IsCodeValid(relationId int64, code string) bool
	CodeExchange(relationId int64, code string) error
}

type CdkeyProducer struct {
	relationId int64
	mgr        CdkeyManager
	lock       sync.Mutex
}

func NewCdkeyProducer(relationId int64, mgr CdkeyManager) *CdkeyProducer {
	return &CdkeyProducer{
		relationId: relationId,
		mgr:        mgr,
	}
}

func (m *CdkeyProducer) GenCode(max, cur int) (string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	code := m.mgr.MakeCode(m.relationId)
	if m.mgr.IsCodeExist(code) {
		if max > 0 && cur >= max {
			return "", ErrReachMaxTimes
		}
		return m.GenCode(max, cur+1)
	}
	if err := m.mgr.CodeStore(m.relationId, code); err != nil {
		return "", err
	}
	return code, nil
}

func (m *CdkeyProducer) GenCodes(num, max, cur int) ([]string, error) {
	if num < 1 {
		return nil, ErrMinNumOne
	}
	var codes []string
	for i := 0; i < num; i++ {
		code, err := m.GenCode(max, cur)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}

func (m *CdkeyProducer) ExchangeCode(code string) error {
	if !m.mgr.IsCodeValid(m.relationId, code) {
		return ErrInvalidCode
	}
	return m.mgr.CodeExchange(m.relationId, code)
}
