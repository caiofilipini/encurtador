package main

import (
	"fmt"
	"github.com/caiofilipini/encurtador/url"
	"net/http"
	"strings"
)

var (
	dominio string
	porta   string
)

func init() {
	dominio = "localhost"
	porta = "8888"
}

func Encurtador(w http.ResponseWriter, r *http.Request) {
	rawBody := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(rawBody)
	body := string(rawBody)

	url := url.NovaUrl(body)

	w.Header().Set("Location", fmt.Sprintf("http://%s:%s/r/%s", dominio, porta, url.Id))
	w.WriteHeader(http.StatusCreated)
}

func Redirecionar(w http.ResponseWriter, r *http.Request) {
	caminho := strings.Split(r.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if url := url.Buscar(id); url != nil {
		http.Redirect(w, r, url.Destino, http.StatusMovedPermanently)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	http.HandleFunc("/r/", Redirecionar)
	http.HandleFunc("/api/encurtar", Encurtador)
	http.ListenAndServe(":"+porta, nil)
}
