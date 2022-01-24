package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate"
	"github.com/joho/godotenv"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) <= 0 {
		useage()
		return
	}
	command := args[0]

	mustLoadEnv()

	fmt.Println(command)

	m, err := migrate.New(
		filePath(),
		connstr(),
	)
	defer m.Close()
	if err != nil {
		log.Fatal(err)
	}

	m.Log = Logger{}

	exec(m, command, args...)
}

func exec(m *migrate.Migrate, command string, args ...string) {
	switch command {
	case "up":
		up(m)
	case "down":
		down(m)
	case "steps":
		count := args[1]
		c, err := strconv.Atoi(count)
		if err != nil {
			fmt.Println("migrate steps [count]")
			os.Exit(1)
			return
		}
		steps(m, c)
	case "drop":
		drop(m)
	default:
		useage()
	}

}

func filePath() string {
	v := os.Getenv("MIGRATION_FILES_PATH")

	if v == "" {
		log.Fatal("MIGRATION_FILES_PATH is not set!!!")
	}

	return v
}

func connstr() string {
	v := os.Getenv("DATABASE_CONNSTR")
	if v == "" {
		log.Fatal("DATABASE_CONNSTR is not set!!!")
	}

	return v
}

func mustLoadEnv() {
	var err error
	var envfile string
	env := os.Getenv("GOOOMA_ENV")
	if "" == env {
		err = godotenv.Load()
		envfile = ".env"
	} else {
		envfile = ".env" + env
		err = godotenv.Load(envfile)
	}

	if err != nil {
		log.Fatalf("Error loading %s\n", envfile)
	}
}

func useage() {
	fmt.Println("migrate [up|down|steps|drop]")
}

func up(m *migrate.Migrate) {
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

func down(m *migrate.Migrate) {
	if err := m.Down(); err != nil {
		log.Fatal(err)
	}
}

func steps(m *migrate.Migrate, count int) {
	if err := m.Steps(count); err != nil {
		log.Fatal(err)
	}
}

func drop(m *migrate.Migrate) {
	if err := m.Drop(); err != nil {
		log.Fatal(err)
	}
}

type Logger struct{}

func (l Logger) Verbose() bool {
	return true
}

func (l Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
