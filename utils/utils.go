package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
	jsoniter "github.com/json-iterator/go"
)

var (
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	QuickCrash = false
)

func SetQuickCrash(crash bool) {
	QuickCrash = crash
}

func StringToInt(s string) int {
	if res, err := strconv.Atoi(s); err != nil {
		logger.Debug(err)
		return 0
	} else {
		return res
	}
}

func StringToInt32(s string) int32 {
	if res, err := strconv.Atoi(s); err != nil {
		logger.Debug(err)
		return 0
	} else {
		return int32(res)
	}
}

func StringToInt64(s string) int64 {
	if res, err := strconv.Atoi(s); err != nil {
		logger.Debug(err)
		return 0
	} else {
		return int64(res)
	}
}
func StringToUint(s string) uint {
	if res, err := strconv.Atoi(s); err != nil {
		logger.Debug(err)
		return 0
	} else {
		return uint(res)
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

//Int32ToString 整数转字符串
func Int32ToString(val int32) string {
	return strconv.Itoa(int(val))
}

//Int64ToString 整数转字符串
func Int64ToString(val int64) string {
	return strconv.Itoa(int(val))
}

func SafeRun(f func()) {
	defer func() {
		if QuickCrash {
			return
		}
		if x := recover(); x != nil {
			txt := fmt.Sprintf("SafeRun panic recover stack : %+v\n", x)
			i := 0
			funcName, file, line, ok := runtime.Caller(i)
			for ok {
				txt += fmt.Sprintf("SafeRun panic frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
				i++
				funcName, file, line, ok = runtime.Caller(i)
			}
			logger.Report(txt)
		}
	}()

	if f != nil {
		f()
	}
}

func Recovery() {
	if QuickCrash {
		return
	}
	if x := recover(); x != nil {
		txt := fmt.Sprintf("service panic stop stack : %+v\n", x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			txt += fmt.Sprintf("service panic frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}
		logger.Report(txt)
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

// 生成区间[-m, n]的安全随机数
func RangeRand(min, max int) int {
	if min > max || min < 0 {
		panic("param is wrong!")
	}
	return rand.Intn(max-min+1) + min
}

func Stringify(data any) string {
	bs, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
	}
	return string(bs)
}

func JsonMarshal(data any) ([]byte, error) {
	return json.Marshal(data)
}

func JsonUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// 获取环境变量信息
func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}
