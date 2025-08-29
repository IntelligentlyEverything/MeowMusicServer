package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	TAG = "MeowMusicServer"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("%s Loading .env file failed: %v\n", TAG, err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("%s PORT environment variable not set\n", TAG)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/file/", fileHandler)
	fmt.Printf("%s Started.\n喵波音律-音乐家园QQ交流群:865754861\n", TAG)
	fmt.Printf("Starting music server at port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("%s Failed to start server: %v\n", TAG, err)
	}
}
