package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func goDotEnvVariable(key string) string {

	err := godotenv.Load("./build/.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
