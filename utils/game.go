package utils

import (
	"fmt"
	"time"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/stores"
)

//产生分布式唯一房号
func NextRoomCode(key string, low, high int) int {
	lockKey := key + "Lock"
	if stores.Lock(lockKey) {
		defer stores.Unlock(lockKey)
		if code, err := stores.Incr(key); err == nil {
			roomCode := int(code)
			isValid := true
			if roomCode > high {
				isValid = false
				roomCode = low + roomCode%(high+1)
			} else if roomCode < low {
				isValid = false
				roomCode = low + roomCode%(high-low)
			}
			if roomCode > high {
				panic(fmt.Sprintf("big than max low:%v,high:%v,code:%v", low, high, code))
			}
			if isValid {
				return roomCode
			}
			if err := stores.Set(key, roomCode, -1); err == nil {
				return roomCode
			} else {
				logger.Error(err)
				goto Finally
			}
		} else {
			logger.Error(err)
			goto Finally
		}
	}
Finally:
	time.Sleep(3 * time.Second)
	return NextRoomCode(key, low, high)
}
