package utils

import (
	"time"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/stores"
)

//产生分布式唯一房号
func NextRoomCode(key string, low, high int) int {
	if stores.Lock(key) {
		defer stores.Unlock(key)
		roomCode := stores.GetInt(key)
		if roomCode == 0 || roomCode == high {
			roomCode = low
		} else {
			roomCode = roomCode + 1
		}
		if err := stores.Set(key, roomCode, -1); err == nil {
			return roomCode
		} else {
			logger.Error(err)
			goto Finally
		}
	}
Finally:
	time.Sleep(time.Second)
	return NextRoomCode(key, low, high)
}
