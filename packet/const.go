package packet

const (
	_ int32 = iota
	StatusStart
	StatusHandshake
	StatusWorking
	StatusClosed
)
