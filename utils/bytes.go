package utils

//big-end 大端序

//BigBytesToInt 字节转整形
func BigBytesToInt(b []byte) int {
	result := 0
	for _, v := range b {
		result = result<<8 + int(v)
	}
	return result
}

// BigIntToBytes 整形转字节
func BigIntToBytes(n int) []byte {
	buf := make([]byte, 3)
	buf[0] = byte((n >> 16) & 0xFF)
	buf[1] = byte((n >> 8) & 0xFF)
	buf[2] = byte(n & 0xFF)
	return buf
}

//little-end

// LittleBytesToInt 字节转整形
func LittleBytesToInt(b []byte) int {
	result := 0
	for k, v := range b {
		result = result + int(v)<<(k*8)
	}
	return result
}

func LittleBytesToUint32(b []byte) uint32 {
	var result uint32 = 0
	for k, v := range b {
		result = result + int(v)<<(k*8)
	}
	return result
}

// LittleIntToBytes 整形转字节
func LittleIntToBytes(n int) []byte {
	buf := make([]byte, 4)
	buf[0] = byte(n & 0xFF)
	buf[1] = byte((n >> 8) & 0xFF)
	buf[2] = byte((n >> 16) & 0xFF)
	buf[3] = byte((n >> 24) & 0xFF)
	return buf
}

func LittleUInt32ToBytes(n uint32) []byte {
	buf := make([]byte, 4)
	buf[0] = byte(n & 0xFF)
	buf[1] = byte((n >> 8) & 0xFF)
	buf[2] = byte((n >> 16) & 0xFF)
	buf[3] = byte((n >> 24) & 0xFF)
	return buf
}
