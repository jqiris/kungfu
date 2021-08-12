package logger

import "time"

type Option func(l *Logger)

func WithLogLevel(level string) Option {
	return func(l *Logger) {
		if logLevel, ok := DescLevelMap[level]; ok {
			l.logLevel = logLevel
		}
	}
}

func WithOutType(typ string) Option {
	return func(l *Logger) {
		if outType, ok := DescOutTypeMap[typ]; ok {
			l.outType = outType
		}
	}
}

func WithLogName(name string) Option {
	return func(l *Logger) {
		l.logName = name
	}
}

func WithLogDir(dir string) Option {
	return func(l *Logger) {
		l.logDir = dir
	}
}

func WithLogDump(dump bool) Option {
	return func(l *Logger) {
		l.logDump = dump
	}
}
func WithLogRuntime(record bool) Option {
	return func(l *Logger) {
		l.logRuntime = record
	}
}

func WithLogFullPath(fullPath bool) Option {
	return func(l *Logger) {
		l.logFullPath = fullPath
	}
}

func WithTimeFormat(format string) Option {
	return func(l *Logger) {
		l.timeFormat = format
	}
}

func WithStdColor(show bool) Option {
	return func(l *Logger) {
		l.stdColor = show
	}
}

func WithZipDuration(d time.Duration) Option {
	return func(l *Logger) {
		if d < defDayDuration {
			d = defDayDuration
		}
		l.zipDuration = d
		l.zipEnd = time.Now()
		l.zipStart = l.zipEnd.Add(-d)
	}
}

func WithTickTime(t time.Duration) Option {
	return func(l *Logger) {
		l.tickTime = t
	}
}
