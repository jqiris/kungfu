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
	"errors"
	"fmt"
)

// PacketType represents the network packet's type such as: handshake and so on.
type PacketType byte

const (
	_ PacketType = iota
	// Handshake represents a handshake: request(client) <====> handshake response(server)
	Handshake = 0x01

	// HandshakeAck represents a handshake ack from client to server
	HandshakeAck = 0x02

	// Heartbeat represents a heartbeat
	Heartbeat = 0x03

	// Data represents a common data packet
	Data = 0x04

	// Kick represents a kick off packet
	Kick = 0x05 // disconnect message from server
)

// ErrWrongPacketType represents a wrong packet type.
var ErrWrongPacketType = errors.New("wrong packet type")

// Packet represents a network packet.
type Packet struct {
	Type   PacketType
	Length int
	Data   []byte
}

//NewPacket create a Packet instance.
func NewPacket() *Packet {
	return &Packet{}
}

//String represents the Packet's in text mode.
func (p *Packet) String() string {
	return fmt.Sprintf("MsgType: %d, Length: %d, Data: %s", p.Type, p.Length, string(p.Data))
}
