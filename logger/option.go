package logger

import (
	"encoding/json"
	"time"
)

type Option func(l *Logger)

func WithLogLevel(level string) Option {
	return func(l *Logger) {
		if logLevel, ok := DescLevelMap[level]; ok {
			l.logLevel = logLevel
		}
	}
}

func WithLogRuntime(record bool) Option {
	return func(l *Logger) {
		l.logRuntime = record
	}
}

func WithTimeFormat(format string) Option {
	return func(l *Logger) {
		l.timeFormat = format
	}
}

func WithReportUrl(url string) Option {
	return func(l *Logger) {
		l.reportUrl = url
	}
}

func WithReportUser(user []string) Option {
	return func(l *Logger) {
		bs, _ := json.Marshal(user)
		l.reportUser = string(bs)
	}
}

/*------------------writer-------------*/
func WithOutType(typ string) Option {
	return func(l *Logger) {
		if outType, ok := DescOutTypeMap[typ]; ok {
			writer.outType = outType
		}
	}
}

func WithLogName(name string) Option {
	return func(l *Logger) {
		writer.logName = name
	}
}

func WithLogDir(dir string) Option {
	return func(l *Logger) {
		writer.logDir = dir
	}
}

func WithLogDump(dump bool) Option {
	return func(l *Logger) {
		writer.logDump = dump
	}
}

func WithStdColor(show bool) Option {
	return func(l *Logger) {
		writer.stdColor = show
	}
}

func WithZipDuration(d time.Duration) Option {
	return func(l *Logger) {
		if d < defDayDuration {
			d = defDayDuration
		}
		writer.zipDuration = d
		writer.zipEnd = time.Now()
		writer.zipStart = writer.zipEnd.Add(-d)
	}
}

func WithTickTime(t time.Duration) Option {
	return func(l *Logger) {
		writer.tickTime = t
	}
}
