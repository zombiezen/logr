// Package logr defines abstract interfaces for logging.  Packages can depend on
// these interfaces and callers can implement logging in whatever way is
// appropriate.
//
// This design derives from Dave Cheney's blog:
//     http://dave.cheney.net/2015/11/05/lets-talk-about-logging
//
// This is a BETA grade API.  Until there is a significant 2nd implementation,
// I don't really know how it will change.
package logr

import "fmt"

// TODO: consider structured logging, a la uber-go/zap
// TODO: consider other bits of glog functionality like Flush, InfoDepth, OutputStats

// InfoLogger represents the ability to log non-error messages.
type InfoLogger interface {
	LogInfo(level int, fields map[string]interface{}, msg string)
}

// ErrorLogger represents the ability to log error messages.
type ErrorLogger interface {
	LogError(fields map[string]interface{}, msg string)
}

// Info logs non-error messages.
type Info struct {
	logger InfoLogger
	level  int
	fields map[string]interface{}
	prefix string
}

// NewInfo creates a new Info that logs to the given logger.
func NewInfo(logger InfoLogger) *Info {
	return &Info{logger: logger}
}

// Info calls LogInfo to its logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Info) Info(args ...interface{}) {
	s := fmt.Sprintln(args...)
	l.logger.LogInfo(l.level, l.fields, l.prefix+s[:len(s)-1])
}

// Infof calls LogInfo to its logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Info) Infof(format string, args ...interface{}) {
	l.logger.LogInfo(l.level, l.fields, l.prefix+fmt.Sprintf(format, args...))
}

// V returns a new log at the specific verbosity level.
// A higher verbosity level means a log message is less important.
func (l *Info) V(level int) *Info {
	return &Info{
		logger: l.logger,
		fields: l.fields,
		level:  level,
		prefix: l.prefix,
	}
}

// WithFields returns a new log with the given fields.
func (l *Info) WithFields(fields map[string]interface{}) *Info {
	f := make(map[string]interface{}, len(fields)+len(l.fields))
	for k, v := range l.fields {
		f[k] = v
	}
	for k, v := range fields {
		f[k] = v
	}
	return &Info{
		logger: l.logger,
		fields: f,
		level:  l.level,
		prefix: l.prefix,
	}
}

// WithPrefix returns a new log that prefixes all messages with a given string.
func (l *Info) WithPrefix(prefix string) *Info {
	return &Info{
		logger: l.logger,
		fields: l.fields,
		level:  l.level,
		prefix: l.prefix + prefix,
	}
}

// Log logs both error and non-error messages.
type Log struct {
	info Info
	err  ErrorLogger
}

// New creates a new log that sends output to the given loggers.
func New(info InfoLogger, err ErrorLogger) *Log {
	return &Log{info: Info{logger: info}, err: err}
}

// AsInfo returns the log as an info log.
func (l *Log) AsInfo() *Info {
	return &l.info
}

// Info calls LogInfo to its logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Log) Info(args ...interface{}) {
	l.info.Info(args...)
}

// Infof calls LogInfo to its logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Log) Infof(format string, args ...interface{}) {
	l.info.Infof(format, args...)
}

// Error calls LogError to its logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Log) Error(args ...interface{}) {
	s := fmt.Sprintln(args...)
	l.err.LogError(l.info.fields, l.info.prefix+s[:len(s)-1])
}

// Errorf calls LogError to its logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Log) Errorf(format string, args ...interface{}) {
	l.err.LogError(l.info.fields, l.info.prefix+fmt.Sprintf(format, args...))
}

// V returns a new log at the specific verbosity level.
// A higher verbosity level means a log message is less important.
func (l *Log) V(level int) *Info {
	return l.info.V(level)
}

// WithFields returns a new log with the given fields.
func (l *Log) WithFields(fields map[string]interface{}) *Log {
	return &Log{info: *l.info.WithFields(fields), err: l.err}
}

// WithPrefix returns a new log that prefixes all messages with a given string.
func (l *Log) WithPrefix(prefix string) *Log {
	return &Log{info: *l.info.WithPrefix(prefix), err: l.err}
}
