package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestListFiles Test ListFiles function
func TestListFiles(t *testing.T) {
	// Create temporary directories and files
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Cannot create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Create some test files
	testFiles := []string{"file1.txt", "subdir/file2.txt", "file3.mp3"}
	for _, f := range testFiles {
		if err := os.MkdirAll(filepath.Join(tempDir, filepath.Dir(f)), 0777); err != nil {
			t.Fatalf("Cannot create subdirectory: %s", err)
		}
		if err := ioutil.WriteFile(filepath.Join(tempDir, f), []byte(""), 0666); err != nil {
			t.Fatalf("Cannot create file: %s", err)
		}
	}

	// Call ListFiles function
	files, err := ListFiles(tempDir)
	if err != nil {
		t.Fatalf("ListFiles function returns error: %s", err)
	}

	// Check if the returned file list is correct
	if len(files) != len(testFiles) {
		t.Errorf("Expected file count is %d, but actual is %d", len(testFiles), len(files))
	}

	for _, f := range testFiles {
		expectedPath := filepath.Join(tempDir, f)
		found := false
		for _, file := range files {
			if file == expectedPath {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found", expectedPath)
		}
	}
}

// TestGetFileContent Test GetFileContent function
func TestGetFileContent(t *testing.T) {
	// Create temporary directories and files
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Cannot create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	testFilePath := filepath.Join(tempDir, "testfile.txt")
	if err := ioutil.WriteFile(testFilePath, []byte("100% Lovely"), 0666); err != nil {
		t.Fatalf("Cannot create file: %s", err)
	}

	// Call GetFileContent function
	content, err := GetFileContent(testFilePath)
	if err != nil {
		t.Fatalf("GetFileContent function returns error: %s", err)
	}

	// Check if the returned content is correct
	expectedContent := "100% Lovely"
	if !strings.EqualFold(string(content), expectedContent) {
		t.Errorf("Expected file content is %s, but actual is %s", expectedContent, string(content))
	}
}

// TestFileHandlerAudio Test fileHandler function for audio files
func TestFileHandlerAudio(t *testing.T) {
	// Create temporary directories and files
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Cannot create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	testFilePath := filepath.Join(tempDir, "testfile.mp3")
	if err := ioutil.WriteFile(testFilePath, []byte("binary audio data"), 0666); err != nil {
		t.Fatalf("Cannot create file: %s", err)
	}

	// Set HTTP request
	req, err := http.NewRequest("GET", "/file/testfile.mp3", nil)
	if err != nil {
		t.Fatalf("Cannot create HTTP request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fileHandler)

	// Execute request
	handler.ServeHTTP(rr, req)

	// Check response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code error: got %v, want %v", status, http.StatusOK)
	}

	// Check Content-Type header
	expectedContentType := "audio/mpeg"
	actualContentType := rr.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Content-Type error: got %v, want %v", actualContentType, expectedContentType)
	}

	// Check response body
	expectedBody := "binary audio data"
	if rr.Body.String() != expectedBody {
		t.Errorf("Response body error: got %v, want %v", rr.Body.String(), expectedBody)
	}
}

// TestFileHandlerImage Test fileHandler function for image files
func TestFileHandlerImage(t *testing.T) {
	// Create temporary directories and files
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Cannot create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	testFilePath := filepath.Join(tempDir, "testfile.jpg")
	if err := ioutil.WriteFile(testFilePath, []byte("binary image data"), 0666); err != nil {
		t.Fatalf("Cannot create file: %s", err)
	}

	// Set HTTP request
	req, err := http.NewRequest("GET", "/file/testfile.jpg", nil)
	if err != nil {
		t.Fatalf("Cannot create HTTP request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fileHandler)

	// Execute request
	handler.ServeHTTP(rr, req)

	// Check response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code error: got %v, want %v", status, http.StatusOK)
	}

	// Check Content-Type header
	expectedContentType := "image/jpeg"
	actualContentType := rr.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Content-Type error: got %v, want %v", actualContentType, expectedContentType)
	}

	// Check response body
	expectedBody := "binary image data"
	if rr.Body.String() != expectedBody {
		t.Errorf("Response body error: got %v, want %v", rr.Body.String(), expectedBody)
	}
}

// TestFileHandlerNotFound Test fileHandler function for nonexistent files
func TestFileHandlerNotFound(t *testing.T) {
	// Create temporary directories and files
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Cannot create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Set HTTP request
	req, err := http.NewRequest("GET", "/file/nonexistentfile.mp3", nil)
	if err != nil {
		t.Fatalf("Cannot create HTTP request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fileHandler)

	// executive request
	handler.ServeHTTP(rr, req)

	// Check response status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Status code error: got %v, want %v", status, http.StatusNotFound)
	}
}

// TestFileHandlerDefaultContentType Test fileHandler function for default content type
func TestFileHandlerDefaultContentType(t *testing.T) {
	// Create temporary directories and files
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	testFilePath := filepath.Join(tempDir, "testfile.txt")
	if err := ioutil.WriteFile(testFilePath, []byte("text data"), 0666); err != nil {
		t.Fatalf("Cannot create file: %s", err)
	}

	// Set HTTP request
	req, err := http.NewRequest("GET", "/file/testfile.txt", nil)
	if err != nil {
		t.Fatalf("Cannot create HTTP request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fileHandler)

	// Execute request
	handler.ServeHTTP(rr, req)

	// Check response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code error: got %v, want %v", status, http.StatusOK)
	}

	// Check Content-Type header
	expectedContentType := "application/octet-stream"
	actualContentType := rr.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Content-Type error: got %v, want %v", actualContentType, expectedContentType)
	}

	// Check response body
	expectedBody := "text data"
	if rr.Body.String() != expectedBody {
		t.Errorf("Response body error: got %v, want %v", rr.Body.String(), expectedBody)
	}
}
