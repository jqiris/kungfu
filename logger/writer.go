/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

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
	outType  OutType
	logDir   string
	logName  string
	logDump  bool   //是否转储
	dumpDate string //转储日期
	logFile  *os.File
	stdColor bool //是否标准输出显示彩色
	fileLock *sync.Mutex
	logChan  chan *LogItem
	dumpChan <-chan time.Time //dump通道
}

func newWriter() *Writer {
	nowTime := time.Now()
	w := &Writer{
		outType:  OutStd,
		logDump:  false,
		logFile:  &os.File{},
		stdColor: true,
		fileLock: new(sync.Mutex),
		logChan:  make(chan *LogItem, 300),
	}
	w.nextDumpTime(nowTime)
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

func (w *Writer) nextDumpTime(now time.Time) {
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 1, next.Location())
	w.dumpChan = time.After(next.Sub(now))
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
	logFile := w.getLogFile()
	w.logFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
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
		case now := <-w.dumpChan:
			w.dumpFile(now) //转储文件
			w.nextDumpTime(now)
		case <-ctx.Done():
			fmt.Println("关闭日志写入器")
			return
		}
	}
}

func (w *Writer) dumpFile(nowTime time.Time) {
	if w.outType <= OutStd || !w.logDump {
		return
	}
	nowDate := nowTime.Format("20060102")
	if w.dumpDate != nowDate {
		w.OpenFile()
		w.dumpDate = nowDate
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
