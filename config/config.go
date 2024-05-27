package config

import (
	"github.com/version-1/goooma/config/internal/loader"
	"github.com/version-1/goooma/logger"
)

type Logger interface {
	Verbose() bool
	Printf(format string, v ...any)

	Errorf(format string, v ...any)
	Warnf(format string, v ...any)

	Fatal(v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
}

type Config struct {
	connstr  string
	filePath string
	logger   Logger
}

func (c Config) Logger() Logger {
	return c.logger
}

func (c Config) Connstr() string {
	return c.connstr
}

func (c Config) FilePath() string {
	return c.filePath
}

func New() (*Config, error) {
	l := loader.New()

	connstr, err := l.Connstr()
	if err != nil {
		return nil, err
	}

	filePath, err := l.FilePath()
	if err != nil {
		return nil, err
	}

	return &Config{
		connstr:  connstr,
		filePath: filePath,
		logger:   logger.DefaultLogger{},
	}, nil
}
