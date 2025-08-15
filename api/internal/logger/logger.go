package logger

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	globalLogger *log.Logger
	once         sync.Once
)

func Init() {
	once.Do(func() {
		// https://specifications.freedesktop.org/basedir-spec/latest/index.html#variables
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get home directory: %v", err)
		}

		logDir := filepath.Join(homeDir, ".local", "state")

		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create log dir: %v", err)
		}

		outfile, err := os.Create(filepath.Join(logDir, "yt-dlp.log"))
		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}

		globalLogger = log.New(outfile, "", log.LstdFlags|log.Lshortfile)
	})
}

func Get() *log.Logger {
	if globalLogger == nil {
		Init()
	}
	return globalLogger
}

func Println(v ...interface{}) {
	Get().Println(v...)
}

func Printf(format string, v ...interface{}) {
	Get().Printf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	Get().Fatalf(format, v...)
}
