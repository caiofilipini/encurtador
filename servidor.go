package main

import (
	"fmt"
	"github.com/caiofilipini/encurtador/url"
	"net/http"
	"os"
	"strings"
)

var (
	dominio string
	porta   string
	ids     chan string
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

	url, err := url.BuscarOuCriarNovaUrl(extrairUrl(r))

	if err != nil {
		responderCom(w, http.StatusBadRequest, nil)
		return
	}

	responderCom(w, http.StatusCreated, Headers{
		"Location": fmt.Sprintf("http://%s:%s/r/%s", dominio, porta, url.Id),
		"Link": fmt.Sprintf("<http://%s:%s/api/stats/%s>; rel=\"stats\"", dominio, porta, url.Id),
	})
}

func Redirecionador(w http.ResponseWriter, r *http.Request) {
	id := extrairId(r)

	if url := url.Buscar(id); url != nil {
		http.Redirect(w, r, url.Destino, http.StatusMovedPermanently)
		ids <- id
	} else {
		http.NotFound(w, r)
	}
}

func Coletor(w http.ResponseWriter, r *http.Request) {
	id := extrairId(r)

	if clicks := url.BuscarClicks(id); clicks > -1 {
		fmt.Fprintf(w, "Clicks: %d", clicks)
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

func extrairUrl(r *http.Request) string {
	rawBody := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(rawBody)
	return string(rawBody)
}

func extrairId(r *http.Request) string {
	caminho := strings.Split(r.URL.Path, "/")
	return caminho[len(caminho)-1]
}

func registrarEstatisticas(ids chan string) {
	for id := range ids {
		url.RegistrarClick(id)
		fmt.Printf("Click registrado com sucesso para %s.\n", id)
	}
}

func main() {
	ids = make(chan string)
	defer close(ids)
	go registrarEstatisticas(ids)

	http.HandleFunc("/r/", Redirecionador)
	http.HandleFunc("/api/encurtar", Encurtador)
	http.HandleFunc("/api/stats/", Coletor)

	http.ListenAndServe(":"+porta, nil)
}
