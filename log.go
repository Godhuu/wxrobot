package wxrobot

import "log"

// Logger
type Logger interface {
	Debug(v ...interface{})
	Warn(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}

var std = new(stdLogger)

type stdLogger struct {
}

// Debug
func (std *stdLogger) Debug(v ...interface{}) {
	log.Default().Println("[debug] ", v)
}

// Warn
func (std *stdLogger) Warn(v ...interface{}) {
	log.Default().Println("[warn] ", v)
}

// Info
func (std *stdLogger) Info(v ...interface{}) {
	log.Default().Println("[info] ", v)
}

// Error
func (std *stdLogger) Error(v ...interface{}) {
	log.Default().Println("[error] ", v)
}
