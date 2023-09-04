package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func Load() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
}

func Get(key string) string {
	return os.Getenv(key)
}

func GetEnvToDuration(e string) (d time.Duration, err error) {
	var envValue int
	envValue, err = strconv.Atoi(os.Getenv(e))
	d = time.Duration(envValue) * time.Second
	return
}

func GetEnvToInt(e string) int {
	result, err := strconv.Atoi(os.Getenv(e))
	if err != nil {
		panic(err)
	}
	return result
}
