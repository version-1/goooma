package command

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Logger interface {
	Fatal(v ...any)
	Infof(format string, v ...any)
	Info(v ...any)

	Verbose() bool
	Printf(format string, v ...any)
}

type App interface {
	Logger() Logger
	FilePath() string
	Connstr() string
}

type Command struct {
	cmd string
	app App
	m   *migrate.Migrate
}

func New(app App, cmd string) *Command {
	return &Command{
		app: app,
		cmd: cmd,
	}
}

func (c Command) logger() Logger {
	return c.app.Logger()
}

func (c Command) ExecWithoutDB(ctx context.Context, args ...string) (bool, error) {
	switch c.cmd {
	case "gen":
		err := c.gen(args[1])
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (c *Command) refreshConnection() error {
	if c.m != nil {
		c.m.Close()
		c.m = nil
	}

	if err := c.prepareMigration(); err != nil {
		return err
	}

	return nil
}

func (c Command) Exec(ctx context.Context, args ...string) error {
	if c.m == nil {
		var err error
		c.app.Logger().Info("Preparing migration")
		c.app.Logger().Infof("file location: %s", c.app.FilePath())
		if err = c.prepareMigration(); err != nil {
			return err
		}
	}
	defer c.m.Close()

	c.logger().Infof("command: %s", c.cmd)

	switch c.cmd {
	case "up":
		return c.m.Up()
	case "down":
		return c.m.Down()
	case "steps":
		countStr := args[1]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return err
		}
		return c.m.Steps(count)
	case "drop":
		return c.m.Drop()
	case "reset":
		err := c.m.Drop()
		if err != nil {
			return err
		}
		c.refreshConnection()

		return c.m.Up()
	case "gen":
		return c.gen(args[1])
	default:
		return Useage()
	}
}

func Useage() error {
	return fmt.Errorf("migrate [up|down|steps|drop|reset|gen]")
}

func (c Command) gen(filename string) error {
	path := c.app.FilePath()
	timestr := time.Now().Format("20060102150405")

	up := fmt.Sprintf("%s/%s_%s.up.sql", path, timestr, filename)
	down := fmt.Sprintf("%s/%s_%s.down.sql", path, timestr, filename)

	f1, err := os.Create(up)
	if err != nil {
		return err
	}

	f2, err := os.Create(down)
	if err != nil {
		return err
	}
	defer f1.Close()
	defer f2.Close()

	c.logger().Info("Generate Succsess!!")
	c.logger().Info(up)
	c.logger().Info(down)

	return nil
}

func (c *Command) prepareMigration() error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", c.app.FilePath()),
		c.app.Connstr(),
	)
	if err != nil {
		return err
	}
	m.Log = c.app.Logger()
	c.m = m

	return nil
}
