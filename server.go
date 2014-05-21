package main

import (
	"fmt"
	"github.com/caiofilipini/encurtador/url"
	"net/http"
	"strings"
)

func Encurtador(w http.ResponseWriter, r *http.Request) {
	rawBody := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(rawBody)
	body := string(rawBody)

	url := url.NovaUrl(body)

	w.Header().Set("Location", fmt.Sprintf("http://127.0.0.1:8888/r/%s", url.Id))
	w.WriteHeader(http.StatusCreated)
}

func Redirecionar(w http.ResponseWriter, r *http.Request) {
	caminho := strings.Split(r.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if url := url.Buscar(id); url != nil {
		http.Redirect(w, r, url.Destino, 301)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/r/", Redirecionar)
	http.HandleFunc("/api/encurtar", Encurtador)
	http.ListenAndServe(":8888", nil)
}
