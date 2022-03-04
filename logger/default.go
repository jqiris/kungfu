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

func Fatal(txt ...any) {
	defLogger.Fatal(txt...)
}

func Fatalf(txt string, args ...any) {
	defLogger.Fatalf(txt, args...)
}

func Error(txt ...any) {
	defLogger.Error(txt...)
}

func Errorf(txt string, args ...any) {
	defLogger.Errorf(txt, args...)
}

func Warn(txt ...any) {
	defLogger.Warn(txt...)
}

func Warnf(txt string, args ...any) {
	defLogger.Warnf(txt, args...)
}

func Info(txt ...any) {
	defLogger.Info(txt...)
}

func Infof(txt string, args ...any) {
	defLogger.Infof(txt, args...)
}

func Debug(txt ...any) {
	defLogger.Debug(txt...)
}

func Debugf(txt string, args ...any) {
	defLogger.Debugf(txt, args...)
}
