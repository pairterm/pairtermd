package server

import (
	"io"
	"net/http"
)

func WebpackHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8080" + r.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}
