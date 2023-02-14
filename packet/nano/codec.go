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

package nano

import (
	"bytes"
	"errors"
	"github.com/jqiris/kungfu/v2/utils"
)

// Codec constants.
const (
	HeadLength    = 4
	MaxPacketSize = 64 * 1024
)

// ErrPacketSizeExceed is the error used for encode/decode.
var ErrPacketSizeExceed = errors.New("codec: packet size exceed")

// A Decoder reads and decodes network data slice
type Decoder struct {
	buf  *bytes.Buffer
	size int  // last packet length
	typ  byte // last packet type
}

// NewDecoder returns a new decoder that used for decode network bytes slice.
func NewDecoder() *Decoder {
	return &Decoder{
		buf:  bytes.NewBuffer(nil),
		size: -1,
	}
}

func (c *Decoder) forward() error {
	header := c.buf.Next(HeadLength)
	c.typ = header[0]
	if c.typ < Handshake || c.typ > Kick {
		return ErrWrongPacketType
	}
	c.size = utils.BigBytesToInt(header[1:])

	// packet length limitation
	if c.size > MaxPacketSize {
		return ErrPacketSizeExceed
	}
	return nil
}

// Decode  decode the network bytes slice to packet.Packet(s)
func (c *Decoder) Decode(data []byte) ([]*Packet, error) {
	c.buf.Write(data)

	var (
		packets []*Packet
		err     error
	)
	// check length
	if c.buf.Len() < HeadLength {
		return nil, err
	}

	// first time
	if c.size < 0 {
		if err = c.forward(); err != nil {
			return nil, err
		}
	}

	for c.size <= c.buf.Len() {
		p := &Packet{Type: PacketType(c.typ), Length: c.size, Data: c.buf.Next(c.size)}
		packets = append(packets, p)

		// more packet
		if c.buf.Len() < HeadLength {
			c.size = -1
			break
		}

		if err = c.forward(); err != nil {
			return nil, err

		}

	}

	return packets, nil
}

// Encode create a packet.Packet from  the raw bytes slice and then encode to network bytes slice
// Protocol refs: https://github.com/NetEase/pomelo/wiki/Communication-Protocol
//
// -<type>-|--------<length>--------|-<data>-
// --------|------------------------|--------
// 1 byte packet type, 3 bytes packet data length(big end), and data segment
func Encode(typ PacketType, data []byte) ([]byte, error) {
	if typ < Handshake || typ > Kick {
		return nil, ErrWrongPacketType
	}

	p := &Packet{Type: typ, Length: len(data)}
	buf := make([]byte, p.Length+HeadLength)
	buf[0] = byte(p.Type)

	copy(buf[1:HeadLength], utils.BigIntToBytes(p.Length))
	copy(buf[HeadLength:], data)

	return buf, nil
}
