package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	// load env from .env file
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
}

func GetManagementUrl() string {
	return os.Getenv("RABBITMQ_MANAGEMENT_URL")
}

func GetManagementBasicAuth() (string, string) {
	return os.Getenv("USERNAME"), os.Getenv("PASSWORD")
}

func IsDeleteMode() bool {
	return os.Getenv("DELETE_MODE") == "ON"
}
