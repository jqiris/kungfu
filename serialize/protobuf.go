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

package serialize

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

// ErrWrongValueType is the error used for marshal the value with protobuf encoding.
var ErrWrongValueType = errors.New("protobuf: convert on wrong type value")

// ProtoSerializer implements the serialize.Serializer interface
type ProtoSerializer struct{}

// NewProtoSerializer returns a new Serializer.
func NewProtoSerializer() *ProtoSerializer {
	return &ProtoSerializer{}
}

// Marshal returns the protobuf encoding of v.
func (s *ProtoSerializer) Marshal(v any) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrWrongValueType
	}
	return proto.Marshal(pb)
}

// Unmarshal parses the protobuf-encoded data and stores the result
// in the value pointed to by v.
func (s *ProtoSerializer) Unmarshal(data []byte, v any) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
