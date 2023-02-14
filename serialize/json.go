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
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// JsonSerializer implements the serialize.Serializer interface
type JsonSerializer struct{}

// NewJsonSerializer returns a new Serializer.
func NewJsonSerializer() *JsonSerializer {
	return &JsonSerializer{}
}

// Marshal returns the JSON encoding of v.
func (s *JsonSerializer) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v.
func (s *JsonSerializer) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
