package logging

import (
	"testing"
)

func init() {
	InitLogger(Config{
		Debug:     false,
		InfoFile:  "./storage/info.log",
		ErrorFile: "./storage/error.log",
	})
}

func TestInitLogger(t *testing.T) {
	zapLog.Warn("hhhh ")
}
