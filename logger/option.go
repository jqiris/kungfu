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

func WithReporter(reporter Reporter) Option {
	return func(l *Logger) {
		l.reporter = reporter
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
