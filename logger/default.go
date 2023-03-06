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

var (
	defLogger = NewLogger()
)

func SetLogger(l *Logger) {
	defLogger = l
}

func WithSuffix(suffix string) *Logger {
	return defLogger.WithSuffix(suffix)
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

func Report(txt ...any) {
	defLogger.Report(txt...)
}

func Reportf(txt string, args ...any) {
	defLogger.Reportf(txt, args...)
}
