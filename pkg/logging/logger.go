package logging

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

type Config struct {
	Debug     bool
	InfoFile  string
	ErrorFile string
}

func InitLogger(c Config) {
	initZap(c)
	Logger = zapLog
}
