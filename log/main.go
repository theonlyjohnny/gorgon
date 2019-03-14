package log

import (
	"fmt"
	"os"

	"github.com/theonlyjohnny/jogger-go/logger"
)

//Log is a shared logger across modules
var Log *logger.Logger

func setupLogger() error {
	opts := logger.Config{
		AppName:    "manticore",
		LogLevel:   "debug",
		LogConsole: true,
		LogSyslog:  nil,
	}
	var loggerErr error
	Log, loggerErr = logger.CreateLogger(opts)
	return loggerErr
}

func init() {
	if err := setupLogger(); err != nil {
		fmt.Printf("Failed to setup logger ? %s \n", err)
		os.Exit(1)
	}
}
