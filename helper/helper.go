package helper

import (
	"crypto/md5"
	"fmt"
	"strconv"

	"github.com/jqiris/kungfu/treaty"
)

//FindServerConfig 查找服务器配置
func FindServerConfig(servers map[string]*treaty.Server, serverId string) *treaty.Server {
	if server, ok := servers[serverId]; ok {
		return server
	}
	return nil
}

//Md5 md5加密
func Md5(str string) string {
	data := []byte(str)
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash) //将[]byte转成16进制
}

//IntToString 整数转字符串
func IntToString(val int) string {
	return strconv.Itoa(val)
}
