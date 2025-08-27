package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Api Response.
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Song represents a song information.
type Song struct {
	Num              int    `json:"num"`
	Song             string `json:"song"`
	Singer           string `json:"singer"`
	Cover            string `json:"cover"`
	Url_audition     string `json:"url_audition"`
	Url_standard     string `json:"url_standard"`
	Url_highquality  string `json:"url_highquality"`
	Url_superquality string `json:"url_superquality"`
	Url_lossless     string `json:"url_lossless"`
	Url_hires        string `json:"url_hires"`
	Url_lyric        string `json:"url_lyric"`
}

// API Song list response.
type SongList struct {
	Num    int    `json:"num"`
	Song   string `json:"song"`
	Singer string `json:"singer"`
	Album  string `json:"album"`
	Pay    int    `json:"pay"`
}

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
		response := Response{
			Code: 0,
			Msg:  "API Operation successful.",
			Data: []interface{}{},
		}
		json.NewEncoder(w).Encode(response)
	} else {
		songs_list := pollApis(msg)
		if songs_list != nil {
			response := Response{
				Code: 0,
				Msg:  "API Operation successful.",
				Data: songs_list,
			}
			json.NewEncoder(w).Encode(response)
		} else {
			response := Response{
				Code: 1,
				Msg:  "No resources available.",
				Data: []interface{}{},
			}
			json.NewEncoder(w).Encode(response)
		}
	}
}

// Aggregation API: Get API on other servers and send API response to apiHandler.
//func aggregationAPI(w http.ResponseWriter, r *http.Request) {}

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

// pollApis polls APIs and fetches data from them.
func pollApis(msg string) []SongList {
	urls, types := getApiConfig()
	var songs []SongList
	var num int

	for i, url := range urls {
		if url == "" {
			//log.Printf("Skipping API %d as URL is not provided", i)
			continue
		}

		apiType := types[i]
		if apiType == "" || apiType == "NETEASE" || apiType == "QQ" || apiType == "KUWO" || apiType == "KUGOU" || apiType == "XIAMI" {
			//log.Printf("Skipping API %d as Type is not provided or not supported", i)
			continue
		} else if apiType == "YAOHU" {
			res, err := fetchDataFromYaohuApi(url, msg)
			if err != nil {
				log.Printf("Error fetching data from API %d (%s): %v", i, apiType, err)
				continue
			}
			for _, yaoHuSong := range res {
				pay := 0
				if yaoHuSong.Pay == "[收费]" {
					pay = 1
				}
				song := SongList{
					Num:    num,
					Song:   yaoHuSong.Name,
					Singer: yaoHuSong.Singer,
					Album:  yaoHuSong.Album,
					Pay:    pay,
				}
				songs = append(songs, song)
				num++
			}
		}
	}
	return songs
}

// YaoHuSong represents a song.
type YaoHuSongList struct {
	N      int    `json:"n"`
	Name   string `json:"name"`
	Singer string `json:"singer"`
	Album  string `json:"album"`
	Pay    string `json:"pay"`
}

// YaoHuResponseData represents the response data.
type YaoHuResponseData struct {
	Code int `json:"code"`
	Data struct {
		Songs []YaoHuSongList `json:"songs"`
	} `json:"data"`
}

// fetchDataFromYaohuApi is the response from api.yaohud.cn.
func fetchDataFromYaohuApi(apiURL string, songName string) ([]YaoHuSongList, error) {
	url := fmt.Sprintf("%s%s", apiURL, songName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed, status: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var responseData YaoHuResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}
	if responseData.Code != 200 || len(responseData.Data.Songs) == 0 {
		return nil, fmt.Errorf("invalid response: %d", responseData.Code)
	}
	//log.Printf("Successfully fetched data from API: %+v", responseData.Data.Songs)
	return responseData.Data.Songs, nil
}
