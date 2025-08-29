package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestNotFoundHandler Test the NotFoundHandler function
func TestNotFoundHandler(t *testing.T) {
	// Set the environment variable HOME_URL
	os.Setenv("HOME_URL", "http://example.com/home")
	defer os.Unsetenv("HOME_URL")

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder that will record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotFoundHandler)

	// Call NotFoundHandler
	handler.ServeHTTP(rr, req)

	// Check if the status code of the response is 404
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned status code %d, expected %d", status, http.StatusNotFound)
	}

	// Check the headers of the response
	expectedServer := "MeowMusicServer"
	if server := rr.Header().Get("Server"); server != expectedServer {
		t.Errorf("Server header returned %s, expected %s", server, expectedServer)
	}

	expectedContentType := "text/html; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Content-Type header returned %s, expected %s", contentType, expectedContentType)
	}

	// Check the body of the response is not empty
	if rr.Body.String() == "" {
		t.Errorf("The response body is empty, expected some data.")
	}
}
