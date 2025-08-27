package main

import (
	"net/http"
	"os"
	"path/filepath"
)

// ListFiles 函数：遍历指定目录中的所有文件，并返回文件路径的切片
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

// GetFileContent 函数：读取指定文件的内容并返回
func GetFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// 读取文件内容
	fileContent := make([]byte, fileSize)
	_, err = file.Read(fileContent)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

// fileHandler 函数：处理文件请求
func fileHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	fileName := queryParams.Get("file")
	if fileName == "" {
		NotFoundHandler(w, r)
		return
	}

	// 构造完整的文件路径
	filePath := filepath.Join("./music-uploads", fileName)

	// 获取文件内容
	fileContent, err := GetFileContent(filePath)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	// 设置适当的Content-Type根据文件扩展名
	ext := filepath.Ext(fileName)
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
		http.Error(w, "Unsupported file type", http.StatusUnsupportedMediaType)
		return
	}

	// 写入文件内容到响应
	w.Write(fileContent)
}
