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

package zinx

import "fmt"

type Packet struct {
	Length int    //消息的长度
	Data   []byte //消息内容
}

//NewPacket create a Packet instance.
func NewPacket() *Packet {
	return &Packet{}
}

//String represents the Packet's in text mode.
func (p *Packet) String() string {
	return fmt.Sprintf("Length: %d, Data: %s", p.Length, string(p.Data))
}
