package goooma

import (
	"context"
	"flag"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/version-1/goooma/core/command"
	"github.com/version-1/goooma/core/config"
	"github.com/version-1/goooma/core/logger"
)

type Config interface {
	Connstr() string
	FilePath() string
}

type Goooma struct {
	logger Logger
	config Config
	args   []string
}

type Logger interface {
	Verbose() bool
	Printf(format string, v ...any)

	Errorf(format string, v ...any)
	Warnf(format string, v ...any)

	Fatal(v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
}

func New() (*Goooma, error) {
	c, err := config.New()
	if err != nil {
		return nil, err
	}

	return NewWith(c)
}

func NewWith(config Config) (*Goooma, error) {
	flag.Parse()
	args := flag.Args()
	if len(args) <= 0 {
		return nil, command.Useage()
	}

	return &Goooma{
		args:   args,
		logger: logger.DefaultLogger{},
		config: config,
	}, nil
}

const version = "v0.1.0"

func (g Goooma) Run(ctx context.Context) {
	cmd := g.args[0]
	if cmd == "version" {
		g.logger.Info(version)
		return
	}

	operation := command.New(g, cmd)

	done, err := operation.ExecWithoutDB(ctx, g.args...)
	if err != nil {
		g.logger.Fatal(err)
		return
	}

	if done {
		return
	}

	m, err := g.prepareMigration()
	defer m.Close()
	if err != nil {
		g.logger.Fatal(err)
	}

	operation.SetMigration(m)
	err = operation.Exec(ctx, g.args...)
	if err != nil {
		g.logger.Fatal(err)
	}
}

func (g Goooma) FilePath() string {
	return g.config.FilePath()
}

func (g Goooma) Logger() command.Logger {
	return g.logger
}

func (g Goooma) prepareMigration() (*migrate.Migrate, error) {
	g.logger.Info("Preparing migration")
	g.logger.Infof("file location: %s", g.config.FilePath())
	m, err := migrate.New(
		fmt.Sprintf("file://%s", g.config.FilePath()),
		g.config.Connstr(),
	)
	if err != nil {
		return nil, err
	}
	m.Log = g.logger

	return m, nil
}
