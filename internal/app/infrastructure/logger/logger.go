package logger

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"time"
)

type logger struct {
	*zerolog.Logger
}

var Logger *logger

func Init(cfg *Config) {
	var writers []io.Writer

	if cfg.LogToConsole {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	if cfg.LogToFile {
		writers = append(writers, newRollingFile(cfg))
	}

	mw := io.MultiWriter(writers...)

	l := zerolog.New(mw).With().Timestamp().Logger()

	Logger = &logger{
		Logger: &l,
	}
}

func newRollingFile(cfg *Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(cfg.Directory, cfg.Filename),
		MaxBackups: cfg.MaxBackups,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
	}
}
