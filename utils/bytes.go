/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

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

func LittleBytesToInt32(b []byte) int32 {
	var result int32 = 0
	for k, v := range b {
		result = result + int32(v)<<(k*8)
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

func LittleInt32ToBytes(n int32) []byte {
	buf := make([]byte, 4)
	buf[0] = byte(n & 0xFF)
	buf[1] = byte((n >> 8) & 0xFF)
	buf[2] = byte((n >> 16) & 0xFF)
	buf[3] = byte((n >> 24) & 0xFF)
	return buf
}
