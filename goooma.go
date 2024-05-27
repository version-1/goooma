package goooma

import (
	"context"

	"github.com/version-1/goooma/config"
	"github.com/version-1/goooma/internal/command"
)

type Config interface {
	Connstr() string
	FilePath() string
	Logger() config.Logger
}

type Goooma struct {
	config Config
}

func New() (*Goooma, error) {
	c, err := config.New()
	if err != nil {
		return nil, err
	}

	return NewWith(c)
}

func NewWith(config Config) (*Goooma, error) {
	return &Goooma{
		config: config,
	}, nil
}

const version = "v0.2.0"

func (g Goooma) Run(ctx context.Context, args ...string) {
	if len(args) <= 0 {
		g.Logger().Fatal(command.Useage())
		return
	}

	cmd := args[0]
	if cmd == "version" {
		g.Logger().Info(version)
		return
	}

	operation := command.New(g, cmd)

	done, err := operation.ExecWithoutDB(ctx, args...)
	if err != nil {
		g.Logger().Fatal(err)
		return
	}

	if done {
		return
	}

	err = operation.Exec(ctx, args...)
	if err != nil {
		g.Logger().Fatal(err)
	}
}

func (g Goooma) Up(ctx context.Context) {
	g.Run(ctx, "up")
}

func (g Goooma) Down(ctx context.Context) {
	g.Run(ctx, "down")
}

func (g Goooma) Drop(ctx context.Context) {
	g.Run(ctx, "drop")
}

func (g Goooma) Reset(ctx context.Context) {
	g.Run(ctx, "reset")
}

func (g Goooma) FilePath() string {
	return g.config.FilePath()
}

func (g Goooma) Logger() command.Logger {
	return g.config.Logger()
}

func (g Goooma) Connstr() string {
	return g.config.Connstr()
}
