package main

import (
	"fmt"
	"log"
	"switchyard/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Env file not found")
	}

	config, err := config.LoadConfig()

	if err != nil {
		log.Fatal("Failed to load required configurations", err)
	}

	fmt.Println(config)
}
