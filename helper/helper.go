package helper

import (
	"crypto/md5"
	"fmt"
	"runtime"
	"strconv"

	"github.com/jqiris/kungfu/treaty"
	"github.com/sirupsen/logrus"
	"reflect"
)

var (
	logger = logrus.WithField("package", "helper")
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

func SafeRun(f func()) {
	defer func() {
		if x := recover(); x != nil {
			logger.Errorf("SafeRun panic recover stack : %+v", x)
			i := 0
			funcName, file, line, ok := runtime.Caller(i)
			for ok {
				logger.Errorf("SafeRun frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
				i++
				funcName, file, line, ok = runtime.Caller(i)
			}
		}
	}()

	if f != nil {
		f()
	}
}

func IsNil(i interface{}) bool {
	defer func() {
		recover()
	}()
	vi := reflect.ValueOf(i)
	return vi.IsNil()
}
