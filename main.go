package main

import (
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Setup a logger to a file
	f, err := os.OpenFile("mock-http-server-go.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	log.Println("--- Starting Mock HTTP Server ---")

	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		log.Println("Could not load .env file")
	}

	SetupServerHandlers()
	StartServer()
}
