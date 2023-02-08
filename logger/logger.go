package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// 支持日志分级，支持按照天转存，支持并发写入
var err error

type LogItem struct {
	logLevel   LogLevel
	allowLevel LogLevel
	logRuntime bool   //是否记录运行时信息
	timeFormat string //日期显示格式
	logTime    time.Time
	logFile    string
	logLine    int
	logTxt     []any
	logSuffix  string //日志前缀
}

func (l *LogItem) logFormat() string {
	logTxt := ""
	if len(l.logSuffix) > 0 {
		if l.logRuntime || needRuntime(l.logLevel) {
			format := "[%s %s] %s [%s file:%s line:%d]"
			logTxt = fmt.Sprintf(format, LevelDescMap[l.logLevel], l.logTime.Format(l.timeFormat), fmt.Sprint(l.logTxt...), l.logSuffix, l.logFile, l.logLine)
		} else {
			format := "[%s %s] %s [%s]"
			logTxt = fmt.Sprintf(format, LevelDescMap[l.logLevel], l.logTime.Format(l.timeFormat), fmt.Sprint(l.logTxt...), l.logSuffix)
		}
	} else {
		if l.logRuntime || needRuntime(l.logLevel) {
			format := "[%s %s] %s [file:%s line:%d]"
			logTxt = fmt.Sprintf(format, LevelDescMap[l.logLevel], l.logTime.Format(l.timeFormat), fmt.Sprint(l.logTxt...), l.logFile, l.logLine)
		} else {
			format := "[%s %s] %s"
			logTxt = fmt.Sprintf(format, LevelDescMap[l.logLevel], l.logTime.Format(l.timeFormat), fmt.Sprint(l.logTxt...))
		}
	}
	return logTxt
}

type Logger struct {
	logLevel   LogLevel
	logRuntime bool     //是否记录运行时信息
	timeFormat string   //日期显示格式
	reporter   Reporter //上报者
	logSuffix  string   //日志后缀
}

func NewLogger(options ...Option) *Logger {
	l := &Logger{
		logLevel:   DEBUG,
		logRuntime: false,
		timeFormat: "2006-01-02 15:04:05",
		logSuffix:  "",
	}
	for _, option := range options {
		option(l)
	}
	writer.initLogger()
	return l
}

func (l *Logger) WithSuffix(suffix string) *Logger {
	return &Logger{
		logLevel:   l.logLevel,
		logRuntime: l.logRuntime,
		timeFormat: l.timeFormat,
		reporter:   l.reporter,
		logSuffix:  suffix,
	}
}

func (l *Logger) GetCallerPath(file string) string {
	fileArr := strings.Split(file, "/")
	fileLen := len(fileArr)
	if fileLen > 2 {
		return fileArr[fileLen-2] + "/" + fileArr[fileLen-1]
	}
	return file
}

func (l *Logger) NewLogItem(level LogLevel, txt ...any) *LogItem {
	item := &LogItem{
		logLevel:   level,
		allowLevel: l.logLevel,
		logTime:    time.Now(),
		logTxt:     txt,
		logRuntime: l.logRuntime,
		timeFormat: l.timeFormat,
		logSuffix:  l.logSuffix,
	}
	if l.logRuntime || needRuntime(level) {
		_, file, line, ok := runtime.Caller(4)
		if ok {
			item.logFile = l.GetCallerPath(file)
			item.logLine = line
		}
	}
	return item
}

func (l *Logger) Fatal(txt ...any) {
	item := l.NewLogItem(FATAL, txt...)
	writer.logChan <- item
}

func (l *Logger) Fatalf(tmp string, args ...any) {
	txt := fmt.Sprintf(tmp, args...)
	l.Fatal(txt)
}

func (l *Logger) Error(txt ...any) {
	item := l.NewLogItem(ERROR, txt...)
	writer.logChan <- item
}

func (l *Logger) Errorf(tmp string, args ...any) {
	txt := fmt.Sprintf(tmp, args...)
	l.Error(txt)
}

func (l *Logger) Warn(txt ...any) {
	item := l.NewLogItem(WARN, txt...)
	writer.logChan <- item
}

func (l *Logger) Warnf(tmp string, args ...any) {
	txt := fmt.Sprintf(tmp, args...)
	l.Warn(txt)
}

func (l *Logger) Info(txt ...any) {
	item := l.NewLogItem(INFO, txt...)
	writer.logChan <- item
}

func (l *Logger) Infof(tmp string, args ...any) {
	txt := fmt.Sprintf(tmp, args...)
	l.Info(txt)
}

func (l *Logger) Debug(txt ...any) {
	item := l.NewLogItem(DEBUG, txt...)
	writer.logChan <- item
}

func (l *Logger) Debugf(tmp string, args ...any) {
	txt := fmt.Sprintf(tmp, args...)
	l.Debug(txt)
}

func (l *Logger) Report(txt ...any) {
	item := l.NewLogItem(REPORT, txt...)
	l.reportItem(item)
	writer.logChan <- item
}

func (l *Logger) Reportf(tmp string, args ...any) {
	txt := fmt.Sprintf(tmp, args...)
	l.Report(txt)
}

func (l *Logger) reportItem(item *LogItem) {
	content := item.logFormat()
	if len(content) > 0 && l.reporter != nil {
		l.reporter(content)
	}
}

func needRuntime(logLevel LogLevel) bool {
	return logLevel <= ERROR
}
