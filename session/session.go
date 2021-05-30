package session

import (
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/stores"
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

var (
	logger     = logrus.WithField("package", "session")
	sessionKey = "kungfu:session"
)

func GetSession(uid int32) *treaty.Session {
	uField := helper.IntToString(int(uid))
	if stores.HExists(sessionKey, uField) {
		res := &treaty.Session{}
		if err := stores.HGet(sessionKey, uField, res); err != nil {
			logger.Error(err)
			return nil
		}
		return res
	}
	return nil
}

func SaveSession(uid int32, sess *treaty.Session) error {
	uField := helper.IntToString(int(uid))
	return stores.HSet(sessionKey, uField, sess)
}

func DestorySession(uid int32) error {
	uField := helper.IntToString(int(uid))
	return stores.HDel(sessionKey, uField)
}
