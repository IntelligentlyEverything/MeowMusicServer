package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ListFiles function: Traverse all files in the specified directory and return a slice of the file path
func ListFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// Get Content function: Read the content of a specified file and return it
func GetFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get File Size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// Read File Content
	fileContent := make([]byte, fileSize)
	_, err = file.Read(fileContent)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

// fileHandler function: Handle file requests
func fileHandler(w http.ResponseWriter, r *http.Request) {
	// Obtain the path of the request
	filePath := r.URL.Path

	// Remove the prefix '/file/'
	if strings.HasPrefix(filePath, "/file/") {
		filePath = filePath[len("/file/"):]
	} else {
		NotFoundHandler(w, r)
		return
	}

	// Construct the complete file path
	fullFilePath := filepath.Join("./music-uploads", filePath)

	// Get file content
	fileContent, err := GetFileContent(fullFilePath)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	// Set appropriate Content-Type based on file extension
	ext := filepath.Ext(filePath)
	switch ext {
	case ".mp3":
		w.Header().Set("Content-Type", "audio/mpeg")
	case ".wav":
		w.Header().Set("Content-Type", "audio/wav")
	case ".flac":
		w.Header().Set("Content-Type", "audio/flac")
	case ".aac":
		w.Header().Set("Content-Type", "audio/aac")
	case ".ogg":
		w.Header().Set("Content-Type", "audio/ogg")
	case ".m4a":
		w.Header().Set("Content-Type", "audio/mp4")
	case ".amr":
		w.Header().Set("Content-Type", "audio/amr")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
		return
	}

	// Write file content to response
	w.Write(fileContent)
}
