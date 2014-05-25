package main

import (
	"encoding/json"
	"fmt"
	"github.com/caiofilipini/encurtador/url"
	"net/http"
	"os"
	"strings"
)

var (
	porta   string
	urlBase string
	ids     chan string
)

func init() {
	dominio := lerConfig("DOMINIO", "localhost")
	porta = lerConfig("PORTA", "8888")
	urlBase = fmt.Sprintf("http://%s:%s", dominio, porta)
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
		"Location": fmt.Sprintf("%s/r/%s", urlBase, url.Id),
		"Link":     fmt.Sprintf("<%s/api/stats/%s>; rel=\"stats\"", urlBase, url.Id),
	})
}

func Redirecionador(w http.ResponseWriter, r *http.Request) {
	buscarUrlEExecutar(w, r, func(url *url.Url) {
		http.Redirect(w, r, url.Destino, http.StatusMovedPermanently)
		ids <- url.Id
	})
}

func Visualizador(w http.ResponseWriter, r *http.Request) {
	buscarUrlEExecutar(w, r, func(url *url.Url) {
		json, err := json.Marshal(url.Stats())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responderComJSON(w, string(json))
	})
}

func buscarUrlEExecutar(w http.ResponseWriter, r *http.Request, executor func(*url.Url)) {
	caminho := strings.Split(r.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if url := url.Buscar(id); url != nil {
		executor(url)
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

func responderComJSON(w http.ResponseWriter, resposta string) {
	responderCom(w, http.StatusOK, Headers{"Content-Type": "application/json"})
	fmt.Fprintf(w, resposta)
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
	http.HandleFunc("/api/stats/", Visualizador)

	http.ListenAndServe(":"+porta, nil)
}
