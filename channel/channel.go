package channel

import (
	"github.com/jqiris/kungfu/stores"
	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
)

var (
	logger     = logrus.WithField("package", "channel")
	channelKey = "kungfu:channel"
)

func GetChannel(server *treaty.Server, uid int32) *treaty.GameChannel {
	if stores.HExists(channelKey, server.ServerId) {
		res := make(map[int32]*treaty.GameChannel)
		if err := stores.HGet(channelKey, server.ServerId, &res); err != nil {
			logger.Error(err)
			return nil
		}
		return res[uid]
	}
	return nil
}

func SaveChannel(ch *treaty.GameChannel) error {
	chMap := make(map[int32]*treaty.GameChannel)
	if stores.HExists(channelKey, ch.Backend.ServerId) {
		if err := stores.HGet(channelKey, ch.Backend.ServerId, &chMap); err != nil {
			logger.Error(err)
		}
	}
	chMap[ch.Uid] = ch
	return stores.HSet(channelKey, ch.Backend.ServerId, chMap)
}

func DestroyChannel(backend *treaty.Server, uid int32) error {
	chMap := make(map[int32]*treaty.GameChannel)
	if stores.HExists(channelKey, backend.ServerId) {
		if err := stores.HGet(channelKey, backend.ServerId, &chMap); err != nil {
			logger.Error(err)
		}
	}
	delete(chMap, uid)
	return stores.HSet(channelKey, backend.ServerId, chMap)
}
