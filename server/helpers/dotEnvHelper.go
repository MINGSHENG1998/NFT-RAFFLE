package helpers

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type DotEnvHelper interface {
	GetEnvVariable(key string) string
}

type dotEnvHelperStruct struct{}

const projectDirName = "server"

var dotEnvHelperImpl DotEnvHelper = NewDotEnvHelper()

func NewDotEnvHelper() DotEnvHelper {
	return &dotEnvHelperStruct{}
}

func (d *dotEnvHelperStruct) GetEnvVariable(key string) string {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatal("Error loading .env file in dotEnvHelper.go ", err.Error())
	}

	return os.Getenv(key)
}
