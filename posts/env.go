package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const envPath = "./build/.env"

func goDotEnvVariable(key string) string {

	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
