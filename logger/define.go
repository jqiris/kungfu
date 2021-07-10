package logger

import (
	"github.com/fatih/color"
	"time"
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

const (
	logSuffix      = ".log"
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
