package packet

import "errors"

var (
	ErrSessionOnNotify    = errors.New("current session working on notify mode")
	ErrCloseClosedGroup   = errors.New("close closed group")
	ErrClosedGroup        = errors.New("group closed")
	ErrMemberNotFound     = errors.New("member not found in the group")
	ErrCloseClosedSession = errors.New("close closed session")
	ErrSessionDuplication = errors.New("session has existed in the current group")
)
