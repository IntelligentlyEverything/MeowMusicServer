package main

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestLoadEnvFile(t *testing.T) {
	// Create a temporary. env file for testing purposes
	err := os.WriteFile(".env", []byte("WEBSITE_NAME=MeowRippleMusic\nHOME_URL=http://127.0.0.1:2233\nPORT=2233\n"), 0644)
	if err != nil {
		log.Fatalf("Failed to create .env file for testing: %v\n", err)
	}

	// Load the .env file
	err = godotenv.Load()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Check if the PORT environment variables are loaded correctly
	port := os.Getenv("PORT")
	if port != "2233" {
		t.Fatalf("Expected PORT to be set to '2233', but got '%s'", port)
	}

	// After the test is completed, delete the temporary. env file
	err = os.Remove(".env")
	if err != nil {
		log.Fatalf("Failed to remove .env file after testing: %v\n", err)
	}
}
