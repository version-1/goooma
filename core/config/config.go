package config

import "github.com/version-1/goooma/core/config/loader"

type Config struct {
	connstr  string
	filePath string
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
	}, nil
}
