package helpers

import (
	"nft-raffle/logger"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "server"

var DotEnvHelper IDotEnvHelper = NewDotEnvHelper()

type IDotEnvHelper interface {
	GetEnvVariable(key string) string
}

type dotEnvHelperStruct struct{}

func NewDotEnvHelper() IDotEnvHelper {
	return &dotEnvHelperStruct{}
}

func (d *dotEnvHelperStruct) GetEnvVariable(key string) string {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		logger.Logger.Fatal("Error loading .env file in dotEnvHelper.go " + err.Error())
	}

	return os.Getenv(key)
}
