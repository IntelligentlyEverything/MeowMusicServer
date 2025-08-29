package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestMain Test the behavior of the main function
func TestMain(m *testing.M) {
	// Create a temporary. env file for testing purposes
	tmpEnvFile, err := os.CreateTemp("", "test.env")
	if err != nil {
		log.Fatalf("Failed to create temporary .env file: %v", err)
	}
	defer os.Remove(tmpEnvFile.Name()) // Delete the temporary file after the test is finished

	// Write environment variables for testing purposes
	_, err = tmpEnvFile.WriteString("PORT=8080\n")
	if err != nil {
		log.Fatalf("Failed to write to temporary .env file: %v", err)
	}
	tmpEnvFile.Close()

	// Set the .env file path
	os.Setenv("GOPATH", tmpEnvFile.Name())
	defer os.Unsetenv("GOPATH")

	// Run the tests
	exitCode := m.Run()

	// Ensure that the correct exit code can be obtained when the program exits
	os.Exit(exitCode)
}

// TestIndexHandler Test the behavior of the indexHandler function
func TestIndexHandler(t *testing.T) {
	// Create a test server
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)

	// Call the processing function and pass in the test request
	handler.ServeHTTP(rr, req)

	// Check if the response status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response content meets expectations
	expected := "MeowMusicServer Started.\n喵波音律-音乐家园QQ交流群:865754861\nStarting music server at port 8080\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestAPIHandler Test the behavior of the apiHandler function
func TestAPIHandler(t *testing.T) {
	// Create a test server
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apiHandler)

	// Call the processing function and pass in the test request
	handler.ServeHTTP(rr, req)

	// Check if the response status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestFileHandler Test the behavior of the fileHandler function
func TestFileHandler(t *testing.T) {
	// Create a test server
	req, err := http.NewRequest("GET", "/file", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fileHandler)

	// Call the processing function and pass in the test request
	handler.ServeHTTP(rr, req)

	// Check if the response status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
