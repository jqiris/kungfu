package tests

import (
	"github.com/jqiris/kungfu/coder"
	"testing"
)

type CoderTest struct {
	Name   string `jason:"name"`
	Header string `json:"header"`
}

func TestCoder(t *testing.T) {
	//m := &CoderTest{
	//	Name:   "hello",
	//	Header: "welcome",
	//}
	//res, err := coder.Marshal(m)
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
	//logger.Infof("json encode:%+v", string(res))
	//ms := &CoderTest{}
	//err = coder.Unmarshal(res, ms)
	//if err != nil {
	//	logger.Error(err)
	//	return
	//}
	//logger.Infof("json decode:%+v", ms)
	a := &RequestBettingThrow{
		Roomid: 100011,
	}
	res, err := coder.Marshal(a)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("proto encode:%+v", string(res))
	ms := &RequestBettingThrow{}
	err = coder.Unmarshal(res, ms)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("proto decode:%+v", ms)
}
