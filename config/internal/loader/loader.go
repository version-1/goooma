package loader

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Loader struct{}

func New() *Loader {
	loadEnvIfExist()
	return &Loader{}
}

func (l Loader) Connstr() (string, error) {
	v := os.Getenv("DATABASE_CONNSTR")
	if v == "" {
		return "", fmt.Errorf("DATABASE_CONNSTR is required")
	}

	return v, nil
}

func (l Loader) FilePath() (string, error) {
	v := os.Getenv("MIGRATION_FILES_PATH")

	if v == "" {
		return "", fmt.Errorf("MIGRATION_FILES_PATH is required")
	}

	return v, nil
}

func loadEnvIfExist() {
	var err error
	var envfile string
	env := os.Getenv("ENV")
	if "" == env {
		err = godotenv.Load()
		envfile = ".env"
	} else {
		envfile = ".env." + env
		err = godotenv.Load(envfile)
	}

	if err == nil {
		log.Printf("loading %s\n", envfile)
	}
}
