package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

//支持日志分级，支持按照天转存，支持并发写入
var err error

type LogItem struct {
	logLevel LogLevel
	logTime  time.Time
	logFile  string
	logLine  int
	logTxt   []interface{}
}

type Logger struct {
	logLevel    LogLevel
	outType     OutType
	logDir      string
	logName     string
	logDump     bool   //是否转储
	dumpDate    string //转储日期
	logFile     *os.File
	logChan     chan *LogItem
	fileLock    *sync.Mutex
	logRuntime  bool          //是否记录运行时信息
	timeFormat  string        //日期显示格式
	stdColor    bool          //是否标准输出显示彩色
	zipDuration time.Duration //zip压缩时长
	zipStart    time.Time     //zip压缩开始
	zipEnd      time.Time     //zip压缩结束
	tickTime    time.Duration //检查间隔
}

func NewLogger(options ...Option) (*Logger, context.CancelFunc) {
	nowTime := time.Now()
	l := &Logger{
		logLevel:    DEBUG,
		outType:     OutStd,
		logDump:     false,
		logRuntime:  false,
		fileLock:    new(sync.Mutex),
		logChan:     make(chan *LogItem, 1024),
		timeFormat:  "2006-01-02 15:04:05",
		zipDuration: defZipDuration,
		zipStart:    nowTime.Add(-defZipDuration),
		zipEnd:      nowTime, //默认启动做一次检查
		tickTime:    10 * time.Minute,
	}
	for _, option := range options {
		option(l)
	}
	l.initLogger()
	ctx, cancel := context.WithCancel(context.Background())
	go l.logWriting(ctx)
	return l, cancel
}

func (l *Logger) initLogger() {
	switch l.outType {
	case OutStd:
	case OutFile:
		fallthrough
	case OutAll:
		l.OpenFile()
	}
}

func (l *Logger) getLogFile() string {
	nowDate := time.Now().Format("20060102")
	logName := l.logName
	if strings.HasSuffix(logName, logSuffix) {
		logName = strings.TrimSuffix(logName, logSuffix)
	}
	if l.outType > OutStd && l.logDump {
		logName = logName + "_" + nowDate
	}
	logName = logName + logSuffix
	return path.Join(l.logDir, logName)
}

func (l *Logger) getLogFileByTime(dt time.Time) (*os.File, error) {
	nowDate := dt.Format("20060102")
	logName := l.logName
	if strings.HasSuffix(logName, logSuffix) {
		logName = strings.TrimSuffix(logName, logSuffix)
	}
	if l.outType > OutStd && l.logDump {
		logName = logName + "_" + nowDate
	}
	logName = logName + logSuffix
	file := path.Join(l.logDir, logName)
	return os.OpenFile(file, os.O_RDWR, 7)
}

func (l *Logger) getZipFileName(start, end time.Time) string {
	s, e := start.Format("20060102"), end.Format("20060102")
	logName := l.logName
	if strings.HasSuffix(logName, logSuffix) {
		logName = strings.TrimSuffix(logName, logSuffix)
	}
	file := fmt.Sprintf("%s_%s_%s%s", logName, s, e, zipSuffix)
	return path.Join(l.logDir, file)
}

func (l *Logger) OpenFile() {
	//进行文件转储
	l.fileLock.Lock()
	defer l.fileLock.Unlock()
	if _, err := os.Stat(l.logDir); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(l.logDir, 0766); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
	l.logFile, err = os.OpenFile(l.getLogFile(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatal(err)
	}
}

func (l *Logger) checkDump() {
	if l.outType > OutStd && l.logDump {
		nowTime := time.Now()
		nowDate := nowTime.Format("20060102")
		if l.dumpDate != nowDate {
			l.OpenFile()
			l.dumpDate = nowDate
		}
		if nowTime.After(l.zipEnd) {
			//过了zip压缩时间,进行压缩
			start, end := l.zipStart, l.zipEnd
			l.zipStart, l.zipEnd = nowTime, nowTime.Add(l.zipDuration)
			zipFiles := make([]*os.File, 0)
			for s := start; s.Before(end); s = s.Add(defDayDuration) {
				if s.Format("20060102") == nowDate {
					continue
				}
				if file, err := l.getLogFileByTime(s); err == nil {
					zipFiles = append(zipFiles, file)
				} else {
					l.Debug(err)
				}
			}
			if len(zipFiles) > 0 {
				dest := l.getZipFileName(start, end)
				if exist, _ := PathExists(dest); exist {
					dest = strings.TrimSuffix(dest, zipSuffix) + "_" + nowTime.Format("20060102150405") + zipSuffix
				}
				if err := Compress(zipFiles, dest); err != nil {
					l.Error(err)
				} else {
					//删除压缩文件
					for _, file := range zipFiles {
						_ = file.Close()
						if err = os.Remove(file.Name()); err != nil {
							l.Error(err)
						}
					}
				}
			}
		}
	}
}

func (l *Logger) OutStd(level LogLevel, txt string) {
	if l.stdColor {
		printer := LevelColorMap[level]
		if _, err := printer.Println(txt); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(txt)
	}
}
func (l *Logger) OutFile(txt string) {
	_, err := l.logFile.Write([]byte(txt + "\n"))
	if err != nil {
		l.OutStd(ERROR, err.Error())
	}
}

func (l *Logger) logFormat(item *LogItem) string {
	logTxt := ""
	if l.logRuntime {
		format := "[%s %s] %s [file:%s line:%d]"
		logTxt = fmt.Sprintf(format, LevelDescMap[item.logLevel], item.logTime.Format(l.timeFormat), fmt.Sprint(item.logTxt...), item.logFile, item.logLine)
	} else {
		format := "[%s %s] %s"
		logTxt = fmt.Sprintf(format, LevelDescMap[item.logLevel], item.logTime.Format(l.timeFormat), fmt.Sprint(item.logTxt...))
	}
	return logTxt
}

func (l *Logger) logWriting(ctx context.Context) {
	tick := time.NewTicker(l.tickTime)
	for {
		select {
		case item := <-l.logChan:
			if l.logLevel < item.logLevel {
				continue
			}
			txt := l.logFormat(item)
			if l.outType == OutStd || l.outType == OutAll {
				l.OutStd(item.logLevel, txt)
			}
			if l.outType == OutFile || l.outType == OutAll {
				l.OutFile(txt)
			}
		case <-tick.C:
			l.checkDump() //每隔10分钟检查下转储
		case <-ctx.Done():
			return
		}
	}
}

func (l *Logger) NewLogItem(level LogLevel, txt ...interface{}) *LogItem {
	item := &LogItem{
		logLevel: level,
		logTime:  time.Now(),
		logTxt:   txt,
	}
	if l.logRuntime {
		_, file, line, ok := runtime.Caller(4)
		if ok {
			item.logFile = path.Base(file)
			item.logLine = line
		}
	}
	return item
}
func (l *Logger) Fatal(txt ...interface{}) {
	item := l.NewLogItem(FATAL, txt...)
	l.logChan <- item
	os.Exit(1)
}

func (l *Logger) Fatalf(tmp string, args ...interface{}) {
	txt := fmt.Sprintf(tmp, args...)
	l.Fatal(txt)
}

func (l *Logger) Error(txt ...interface{}) {
	item := l.NewLogItem(ERROR, txt...)
	l.logChan <- item
}

func (l *Logger) Errorf(tmp string, args ...interface{}) {
	txt := fmt.Sprintf(tmp, args...)
	l.Error(txt)
}

func (l *Logger) Warn(txt ...interface{}) {
	item := l.NewLogItem(WARN, txt...)
	l.logChan <- item
}

func (l *Logger) Warnf(tmp string, args ...interface{}) {
	txt := fmt.Sprintf(tmp, args...)
	l.Warn(txt)
}

func (l *Logger) Info(txt ...interface{}) {
	item := l.NewLogItem(INFO, txt...)
	l.logChan <- item
}

func (l *Logger) Infof(tmp string, args ...interface{}) {
	txt := fmt.Sprintf(tmp, args...)
	l.Info(txt)
}

func (l *Logger) Debug(txt ...interface{}) {
	item := l.NewLogItem(DEBUG, txt...)
	l.logChan <- item
}

func (l *Logger) Debugf(tmp string, args ...interface{}) {
	txt := fmt.Sprintf(tmp, args...)
	l.Debug(txt)
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
