package main

import (
	"fmt"
	"github.com/caiofilipini/encurtador/url"
	"net/http"
	"strings"
	"os"
)

var (
	dominio string
	porta   string
)

func init() {
	dominio = lerConfig("DOMINIO", "localhost")
	porta = lerConfig("PORTA", "8888")
}

type Headers map[string]string

func Encurtador(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		responderCom(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})
		return
	}

	rawBody := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(rawBody)
	body := string(rawBody)

	url := url.NovaUrl(body)

	responderCom(w, http.StatusCreated, Headers{
		"Location": fmt.Sprintf("http://%s:%s/r/%s", dominio, porta, url.Id),
	})
}

func Redirecionador(w http.ResponseWriter, r *http.Request) {
	caminho := strings.Split(r.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if url := url.Buscar(id); url != nil {
		http.Redirect(w, r, url.Destino, http.StatusMovedPermanently)
	} else {
		http.NotFound(w, r)
	}
}

func responderCom(w http.ResponseWriter, status int, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
}

func lerConfig(config string, valorPadrao string) string {
	if d := os.Getenv(config); d != "" {
		return d
	}
	return valorPadrao
}

func main() {
	http.HandleFunc("/r/", Redirecionador)
	http.HandleFunc("/api/encurtar", Encurtador)
	http.ListenAndServe(":"+porta, nil)
}
