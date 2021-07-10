package logger

import "context"

var (
	defLogger, defCancel = NewLogger()
)

func SetLogger(l *Logger, cancel context.CancelFunc) {
	defLogger = l
	defCancel()
	defCancel = cancel
}

func Fatal(txt ...interface{}) {
	defLogger.Fatal(txt...)
}

func Fatalf(txt string, args ...interface{}) {
	defLogger.Fatalf(txt, args...)
}

func Error(txt ...interface{}) {
	defLogger.Error(txt...)
}

func Errorf(txt string, args ...interface{}) {
	defLogger.Errorf(txt, args...)
}

func Warn(txt ...interface{}) {
	defLogger.Warn(txt...)
}

func Warnf(txt string, args ...interface{}) {
	defLogger.Warnf(txt, args...)
}

func Info(txt ...interface{}) {
	defLogger.Info(txt...)
}

func Infof(txt string, args ...interface{}) {
	defLogger.Infof(txt, args...)
}

func Debug(txt ...interface{}) {
	defLogger.Debug(txt...)
}

func Debugf(txt string, args ...interface{}) {
	defLogger.Debugf(txt, args...)
}
