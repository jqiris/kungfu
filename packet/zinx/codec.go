package zinx

import (
	"bytes"
	"errors"
	"github.com/jqiris/kungfu/v2/utils"
	"io"
	"net"
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
	size int // last packet length
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
	c.size = utils.LittleBytesToInt(header)

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
		p := &Packet{Length: c.size, Data: c.buf.Next(c.size)}
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

func ReadMsg(conn net.Conn) (*Message, error) {
	buf := make([]byte, HeadLength)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	p := NewPacket()
	p.Length = utils.LittleBytesToInt(buf)
	buf = make([]byte, p.Length)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	p.Data = buf
	return MsgDecode(p.Data)
}

func Encode(msg *Message) ([]byte, error) {
	data, err := MsgEncode(msg)
	if err != nil {
		return nil, err
	}
	p := &Packet{Length: len(data)}
	buf := make([]byte, p.Length+HeadLength)
	copy(buf[:HeadLength], utils.LittleIntToBytes(p.Length))
	copy(buf[HeadLength:], data)
	return buf, nil
}
