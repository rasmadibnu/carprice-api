package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Unable to load .env file")
	}
}
