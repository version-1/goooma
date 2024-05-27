package logger

import (
	"log"
)

type DefaultLogger struct{}

func (l DefaultLogger) Verbose() bool {
	return true
}

func (l DefaultLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l DefaultLogger) Infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l DefaultLogger) Errorf(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func (l DefaultLogger) Warnf(format string, v ...interface{}) {
	log.Printf("[WARN] "+format, v...)
}

func (l DefaultLogger) Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func (l DefaultLogger) Info(v ...interface{}) {
	log.Println(v...)
}
