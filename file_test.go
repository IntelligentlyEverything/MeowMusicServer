package main

import (
	"io/ioutil"
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
	defer os.RemoveAll(tempDir) // Clean up after test

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
		t.Fatalf("ListFiles function returned error: %s", err)
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

// TestGetFileContent Test the GetFile Content function
func TestGetFileContent(t *testing.T) {
	// Create temporary file
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Unable to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "testfile.txt")
	if err := ioutil.WriteFile(testFilePath, []byte("100% Lovely"), 0666); err != nil {
		t.Fatalf("Unable to create file: %s", err)
	}

	// Call GetFileContent function
	content, err := GetFileContent(testFilePath)
	if err != nil {
		t.Fatalf("GetFileContent function returned error: %s", err)
	}

	// Check if the returned content is correct
	expectedContent := "100% Lovely"
	if !strings.EqualFold(string(content), expectedContent) {
		t.Errorf("Expected file content is %s, but actual is %s", expectedContent, string(content))
	}
}
