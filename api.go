package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Api Response.
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Tips string      `json:"tips"`
	Ip   string      `json:"ip"`
}

// API Song response.
type Song struct {
	Num      int         `json:"num"`
	Song     string      `json:"song"`
	Singer   string      `json:"singer"`
	Album    string      `json:"album"`
	Cover    string      `json:"cover"`
	MusicURL interface{} `json:"music_url"`
	Lyric    interface{} `json:"lyric"`
}

type MusicURL struct {
	Audition     string `json:"audition"`
	Standard     string `json:"standard"`
	Highquality  string `json:"highquality"`
	Superquality string `json:"superquality"`
	Lossless     string `json:"lossless"`
	Hires        string `json:"hires"`
}

type Lyric struct {
	Mrc string `json:"mrc"`
	Lrc string `json:"lrc"`
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
	ip, err := IPhandler(r)
	if err != nil {
		ip = "0.0.0.0"
	}
	if msg == "" {
		response := Response{
			Code: 1,
			Msg:  "API Operation successful but no request provided.",
			Data: []interface{}{},
			Tips: "Provide by " + os.Getenv("WEBSITE_NAME"),
			Ip:   ip,
		}
		json.NewEncoder(w).Encode(response)
	}
}

// API response.

// Processing requests.
func IPhandler(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip, nil
	}
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0]), nil
	}
	ip = r.RemoteAddr
	if ip != "" {
		return strings.Split(ip, ":")[0], nil
	}

	return "", fmt.Errorf("unable to obtain IP address information")
}
