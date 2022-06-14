package logger

import (
	"time"

	"github.com/fatih/color"
)

type LogLevel int

const (
	FATAL LogLevel = iota //0-致命错误
	ERROR                 //1-错误
	WARN                  //2-警告
	INFO                  //3-通知
	DEBUG                 //4-调试
)

type OutType int

const (
	OutStd  OutType = iota //0-标准输出
	OutFile                //1-文件输出
	OutAll                 //2-文件和标准都输出
)

type ChangChanType int

const (
	ChangChanZip  ChangChanType = iota //改变zip chan
	ChangChanDump                      // 改变dump chan
)

const (
	logSuffix      = ".log"
	zipSuffix      = ".zip"
	lockSuffix     = ".lock"
	defZipDuration = 744 * time.Hour //31天
	defDayDuration = 24 * time.Hour  //1天

)

var (
	LevelDescMap = map[LogLevel]string{
		FATAL: "FATAL",
		ERROR: "ERROR",
		WARN:  "WARN",
		INFO:  "INFO",
		DEBUG: "DEBUG",
	}
	DescLevelMap = map[string]LogLevel{
		"fatal": FATAL,
		"error": ERROR,
		"warn":  WARN,
		"info":  INFO,
		"debug": DEBUG,
	}
	DescOutTypeMap = map[string]OutType{
		"out_std":  OutStd,
		"out_file": OutFile,
		"out_all":  OutAll,
	}
	LevelColorMap = map[LogLevel]*color.Color{
		FATAL: color.New(color.FgHiMagenta), //紫红色
		ERROR: color.New(color.FgHiRed),     //红色
		WARN:  color.New(color.FgHiYellow),  //黄色
		INFO:  color.New(color.FgHiBlue),    //绿色
		DEBUG: color.New(color.FgHiCyan),    //青蓝色
	}
)

type Reporter func(ctx string)
