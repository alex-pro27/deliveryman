package logger

import (
	"git.samberi.com/dois/delivery_api/config"
	glogger "github.com/google/logger"
	"os"
)

var Logger *glogger.Logger
var LogFile *os.File

func Init() {
	lf, err := os.OpenFile(config.Config.System.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		glogger.Fatalf("Failed to open log file: %v", err)
	}
	Logger = glogger.Init("MainLog", true, true, lf)
	LogFile = lf
}

func HandleError(err error) {
	if err != nil {
		Logger.Error(err)
	}
}

func Close() {
	if err := LogFile.Close(); err != nil {
		glogger.Fatalf("Failed to close log file: %v", err)
	}
	Logger.Close()
}
