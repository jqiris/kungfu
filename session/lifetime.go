package session

import (
	"github.com/jqiris/kungfu/logger"
	"runtime"
	"strings"
)

type (
	// LifetimeHandler represents a callback
	// that will be called when a session close or
	// session low-level connection broken.
	LifetimeHandler func(*Session)

	lifetime struct {
		// callbacks that emitted on session closed
		onClosed []LifetimeHandler
	}
)

var Lifetime = &lifetime{}

// OnClosed set the Callback which will be called
// when session is closed Waring: session has closed.
func (lt *lifetime) OnClosed(h LifetimeHandler) {
	lt.onClosed = append(lt.onClosed, h)
}

func (lt *lifetime) Close(s *Session) {
	if len(lt.onClosed) < 1 {
		return
	}

	for _, h := range lt.onClosed {
		h(s)
	}
}
func OnSessionClosed(s *Session) {
	defer func() {
		if err := recover(); err != nil {
			logger.Infof("onSessionClosed: %v", err)
			println(stack())
		}
	}()

	Lifetime.Close(s)
}

func stack() string {
	buf := make([]byte, 10000)
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	s := string(buf)

	// skip nano frames lines
	const skip = 7
	count := 0
	index := strings.IndexFunc(s, func(c rune) bool {
		if c != '\n' {
			return false
		}
		count++
		return count == skip
	})
	return s[index+1:]
}
