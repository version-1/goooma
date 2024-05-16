package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	done := preexec(command, args...)
	if done {
		return
	}

	m, err := migrate.New(
		"file://"+filePath(),
		connstr(),
	)
	defer m.Close()
	if err != nil {
		log.Fatal(err)
	}

	m.Log = Logger{}

	exec(m, command, args...)
}

func preexec(command string, args ...string) bool {
	switch command {
	case "gen":
		gen(filePath(), args[1])
		return true
	}
	return false
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
	case "gen":
		gen(filePath(), args[1])
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
		envfile = ".env." + env
		err = godotenv.Load(envfile)
	}

	if err != nil {
		log.Fatalf("Error loading %s\n", envfile)
	}
}

func useage() {
	fmt.Println("migrate [up|down|steps|drop|gen]")
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

func gen(path string, filename string) {
	timestr := time.Now().Format("20060102150405")
	up := fmt.Sprintf("%s/%s_%s.up.sql", path, timestr, filename)
	down := fmt.Sprintf("%s/%s_%s.down.sql", path, timestr, filename)

	f1, err := os.Create(up)
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Create(down)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()
	defer f2.Close()

	fmt.Println("Generate Succsess!!")
	fmt.Println(up)
	fmt.Println(down)
}

type Logger struct{}

func (l Logger) Verbose() bool {
	return true
}

func (l Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
