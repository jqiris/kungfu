package session

import (
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/stores"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
)

var (
	sessionKey = "kungfu:session"
)

func GetSession(uid int32) *treaty.Session {
	uField := utils.IntToString(int(uid))
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
	uField := utils.IntToString(int(uid))
	return stores.HSet(sessionKey, uField, sess)
}

func DestorySession(uid int32) error {
	uField := utils.IntToString(int(uid))
	return stores.HDel(sessionKey, uField)
}
