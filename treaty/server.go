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

package treaty

import (
	"encoding/json"
	"path"
)

func RegSeverItem(x *Server) string {
	return path.Join(x.ServerType, x.ServerId)
}
func RegSeverQueue(serverType, queue string) string {
	return path.Join(serverType, queue)
}

func RegServerType(x *Server) string {
	return x.ServerType
}

func RegSerialize(x *Server) string {
	if res, err := json.Marshal(x); err == nil {
		return string(res)
	}
	return ""
}

func RegUnSerialize(s []byte) (*Server, error) {
	x := &Server{}
	if err := json.Unmarshal(s, x); err != nil {
		return nil, err
	}
	return x, nil
}
