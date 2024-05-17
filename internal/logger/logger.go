package logger

import (
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"aim_testwork/internal/conf"

	"github.com/rs/zerolog"
)

// Global logger for project
var (
	once sync.Once
	log  zerolog.Logger
)

// Initialize logger
func Get() zerolog.Logger {
	once.Do(func() {
		cfg := conf.Get()

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		logLevel, err := zerolog.ParseLevel(cfg.Launch.LogLevel)
		if err != nil {
			panic(err)
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if cfg.Launch.AppMode == "release" {
			logFile, err := os.OpenFile(
				"./server.log",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0o664,
			)
			if err != nil {
				panic(err)
			}

			output = zerolog.MultiLevelWriter(os.Stderr, logFile)
		}
		buildInfo, _ := debug.ReadBuildInfo()

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("go_version", buildInfo.GoVersion).
			Logger()
	})

	return log
}
