package utils

import (
	"nft-raffle-cron/logger"
	"os"
	"regexp"
	"sync"

	"github.com/joho/godotenv"
)

const projectDirName = "cron"

var (
	dotEnvUtil     *DotEnvUtil
	dotEnvUtilOnce sync.Once
)

type DotEnvUtil struct{}

func GetDotEnvUtil() *DotEnvUtil {
	if dotEnvUtil == nil {
		dotEnvUtilOnce.Do(func() {
			dotEnvUtil = &DotEnvUtil{}
		})
	}
	return dotEnvUtil
}

func (d *DotEnvUtil) GetEnvVariable(key string) string {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		logger.Logger.Fatal("Error loading .env file in dotEnvUtil.go " + err.Error())
	}

	return os.Getenv(key)
}
