<reference>
<name>#fitten_ref_1#</name>
<type>file</type>
<path>main_test.go</path>
<content>
package main

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

// TestMain Test the logic of the main function
func TestMain(t *testing.T) {
	// Create a temporary. env file for testing purposes
	godotenv.Write(map[string]string{
		"WEBSITE_NAME": "MeowRippleMusic",
		"HOME_URL": "http://127.0.0.1:2233",
		"PORT":         "2233",
	}, ".env.test")

	// Set environment variables to point to temporary. env files
	os.Setenv("GO_ENV", "test")
	os.Setenv("GOPATH", os.Getenv("GOPATH")+"/src/test")

	// Using httptest to simulate HTTP requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Test server is running!")
	}))
	defer server.Close()

	// Redirects standard output to capture log information
	var output strings.Builder
	log.SetOutput(&output)

	// Capture the output of fmt. Prin
	fmtOutput := ""
	fmtBackup := fmt.Printf
	fmt.Printf = func(format string, a ...interface{}) {
		fmtOutput = fmt.Sprintf(format, a...)
	}

	// Call the main function
	os.Args = []string{"main.go"}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main function panicked: %v", r)
		}
	}()
	main()

	// Restore standard output
	log.SetOutput(os.Stderr)
	fmt.Printf = fmtBackup

	// Check if the output of the main function is correct
	expectedLog := "MeowMusicServer Loading .env file failed: open .env: no such file or directory\n"
	expectedOutput := "MeowMusicServer Started.\n喵波音律-音乐家园QQ交流群:865754861\nStarting music server at port 8080\n"

	if output.String() == expectedLog {
		t.Errorf("Expected log output, got: %v", output.String())
	}

	if !strings.Contains(fmtOutput, expectedOutput) {
		t.Errorf("Expected fmt output, got: %v", fmtOutput)
	}

	// Delete temporary. env files
	os.Remove(".env.test")
}

// TestHandler tests various HTTP processing functions
func TestHandlers(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			indexHandler(w, r)
		case "/api":
			apiHandler(w, r)
		case "/file":
			fileHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Test root path
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Errorf("Error getting /: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200 for /, got: %d", resp.StatusCode)
	}

	// Test /api path
	resp, err = http.Get(server.URL + "/api")
	if err != nil {
		t.Errorf("Error getting /api: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200 for /api, got: %d", resp.StatusCode)
	}

	// Test /file path
	resp, err = http.Get(server.URL + "/file")
	if err != nil {
		t.Errorf("Error getting /file: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200 for /file, got: %d", resp.StatusCode)
	}
}

// Simulate the indexHandler function
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the index page!")
}

// Simulate the apiHandler function
func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the API page!")
}

// Simulate the fileHandler function
func fileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the file page!")
}
</content>
</reference>
ExpandCopyInsert
