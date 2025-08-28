package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "MeowMusicServer")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path != "/" {
		NotFoundHandler(w, r)
		return
	}
	fmt.Fprintf(w, "<h1>音乐服务器</h1>")
}
