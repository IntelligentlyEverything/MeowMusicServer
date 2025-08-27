package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// getApiConfig gets the API configuration from environment variables.
func getApiConfig() ([]string, []string) {
	urls := make([]string, 10)
	types := make([]string, 10)

	for i := 0; i < 10; i++ {
		urlKey := fmt.Sprintf("API_URL_%d", i)
		typeKey := fmt.Sprintf("API_TYPE_%d", i)
		if i == 0 {
			urlKey = "API_URL"
			typeKey = "API_TYPE"
		}
		urls[i] = os.Getenv(urlKey)
		types[i] = os.Getenv(typeKey)
	}

	return urls, types
}

// fetchDataFromApi is the response from an API.
func fetchDataFromApi(url string, apiType string) {
	log.Printf(url)
	log.Printf(apiType)
	return
}

// pollApis polls APIs and fetches data from them.
func pollApis() {
	urls, types := getApiConfig()

	for i, url := range urls {
		if url == "" {
			log.Printf("Skipping API %d as URL is not provided", i)
			continue
		}

		apiType := types[i]
		if apiType == "" {
			log.Printf("Skipping API %d as Type is not provided", i)
			continue
		}

		apiResponse, err := fetchDataFromApi(url, apiType)
		if err != nil {
			log.Printf("Error fetching data from API %d (%s): %v", i, apiType, err)
			continue
		}

		fmt.Printf("Successfully fetched data from API %d (%s): %s\n", i, apiType, apiResponse.Data)
	}
}

// Aggregation API: Get API on other servers and send API response to apiHandler.
func aggregationAPI(w http.ResponseWriter, r *http.Request) {}

// apiHandler is the handler function for API requests.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "MeowMusicServer")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	queryParams := r.URL.Query()
	//key := queryParams.Get("key")
	msg := queryParams.Get("msg")
	//num := queryParams.Get("num")
	//quality := queryParams.Get("quality")
	if msg == "" {
		fmt.Fprintf(w, `{"code": "0", "msg": "API Operation successful."}`)
	}
}
