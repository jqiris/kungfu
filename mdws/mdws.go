package mdws

import (
	"fmt"
	"runtime"
)

func stack() string {
	txt := ""
	i := 0
	funcName, file, line, ok := runtime.Caller(i)
	for ok {
		txt += fmt.Sprintf("gin panic frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
		i++
		funcName, file, line, ok = runtime.Caller(i)
	}
	return txt
}
