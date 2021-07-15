package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/treaty"
	"os"
	"runtime"
	"strconv"
)

func StringToInt(s string) int {
	if res, err := strconv.Atoi(s); err != nil {
		logger.Error(err)
		return 0
	} else {
		return res
	}
}

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

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
