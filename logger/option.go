package logger

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
