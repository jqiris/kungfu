package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	writer = newWriter()
)

type Writer struct {
	outType     OutType
	logDir      string
	logName     string
	logDump     bool   //是否转储
	dumpDate    string //转储日期
	logFile     *os.File
	stdColor    bool //是否标准输出显示彩色
	fileLock    *sync.Mutex
	zipDuration time.Duration //zip压缩时长
	zipStart    time.Time     //zip压缩开始
	zipEnd      time.Time     //zip压缩结束
	tickTime    time.Duration //检查间隔
	logChan     chan *LogItem
}

func newWriter() *Writer {
	nowTime := time.Now()
	w := &Writer{
		outType:     OutStd,
		logDump:     false,
		logFile:     &os.File{},
		stdColor:    false,
		zipDuration: defZipDuration,
		zipStart:    nowTime,
		zipEnd:      nowTime.Add(defZipDuration), //zip时间
		tickTime:    10 * time.Minute,
		fileLock:    new(sync.Mutex),
		logChan:     make(chan *LogItem, 100),
	}
	go w.logWriting()
	return w
}

func (w *Writer) initLogger() {
	switch w.outType {
	case OutStd:
	case OutFile:
		fallthrough
	case OutAll:
		w.OpenFile()
	}
}

func (w *Writer) OpenFile() {
	//进行文件转储
	w.fileLock.Lock()
	defer w.fileLock.Unlock()
	if w.logFile != nil {
		w.logFile.Close()
	}
	if _, err := os.Stat(w.logDir); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(w.logDir, 0766); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
	w.logFile, err = os.OpenFile(w.getLogFile(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatal(err)
	}
}

func (w *Writer) getLogFile() string {
	nowDate := time.Now().Format("20060102")
	logName := strings.TrimSuffix(w.logName, logSuffix)
	if w.outType > OutStd && w.logDump {
		logName = logName + "_" + nowDate
	}
	logName = logName + logSuffix
	return path.Join(w.logDir, logName)
}

func (w *Writer) getLogFileByTime(dt time.Time) (*os.File, error) {
	nowDate := dt.Format("20060102")
	logName := strings.TrimSuffix(w.logName, logSuffix)
	if w.outType > OutStd && w.logDump {
		logName = logName + "_" + nowDate
	}
	logName = logName + logSuffix
	file := path.Join(w.logDir, logName)
	return os.OpenFile(file, os.O_RDWR, 7)
}

func (w *Writer) getZipFileName(start, end time.Time) string {
	s, e := start.Format("20060102"), end.Format("20060102")
	logName := strings.TrimSuffix(w.logName, logSuffix)
	file := fmt.Sprintf("%s_%s_%s%s", logName, s, e, zipSuffix)
	return path.Join(w.logDir, file)
}

func (w *Writer) OutStd(level LogLevel, txt string) {
	if w.stdColor {
		printer := LevelColorMap[level]
		if _, err := printer.Println(txt); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(txt)
	}
}
func (w *Writer) OutFile(txt string) {
	_, err := w.logFile.Write([]byte(txt + "\n"))
	if err != nil {
		w.OutStd(ERROR, err.Error())
	}
}

func (w *Writer) logWriting() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case item := <-w.logChan:
			if item.allowLevel < item.logLevel {
				continue
			}
			txt := item.logFormat()
			if w.outType == OutStd || w.outType == OutAll {
				w.OutStd(item.logLevel, txt)
			}
			if w.outType == OutFile || w.outType == OutAll {
				w.OutFile(txt)
			}
			if item.logLevel == FATAL {
				os.Exit(1)
			}
		case <-time.After(w.tickTime):
			w.checkDump() //每隔10分钟检查下转储
		case <-ctx.Done():
			fmt.Println("关闭日志写入器")
			return
		}
	}
}

func (w *Writer) checkDump() {
	if w.outType > OutStd && w.logDump {
		nowTime := time.Now()
		nowDate := nowTime.Format("20060102")
		if w.dumpDate != nowDate {
			w.OpenFile()
			w.dumpDate = nowDate
		}
		if nowTime.After(w.zipEnd) {
			//过了zip压缩时间,进行压缩
			start, end := w.zipStart, w.zipEnd
			w.zipStart, w.zipEnd = nowTime, nowTime.Add(w.zipDuration)
			zipFiles := make([]*os.File, 0)
			for s := start; s.Before(end); s = s.Add(defDayDuration) {
				if s.Format("20060102") == nowDate {
					continue
				}
				if file, err := w.getLogFileByTime(s); err == nil {
					zipFiles = append(zipFiles, file)
				}
			}
			if len(zipFiles) > 0 {
				dest := w.getZipFileName(start, end)
				if exist, _ := PathExists(dest); exist {
					dest = strings.TrimSuffix(dest, zipSuffix) + "_" + nowTime.Format("20060102150405") + zipSuffix
				}
				if err := Compress(zipFiles, dest); err != nil {
					fmt.Println(err)
				} else {
					//删除压缩文件
					for _, file := range zipFiles {
						_ = file.Close()
						if err = os.Remove(file.Name()); err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
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
