package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/caiofilipini/encurtador/url"
)

type headers map[string]string

type Handler struct {
	urlBase   string
	logLigado bool
	stats     chan<- string
}

func NovoHandler(urlBase string, logLigado bool) Handler {
	stats := make(chan string)
	handler := Handler{
		urlBase:   urlBase,
		logLigado: logLigado,
		stats:     stats,
	}

	go handler.registrarEstatisticas(stats)

	return handler
}

func (h Handler) Redirecionar(w http.ResponseWriter, req *http.Request) {
	buscarUrlEExecutar(w, req, func(url *url.Url) {
		http.Redirect(w, req, url.Destino, http.StatusMovedPermanently)
		h.stats <- url.Id
	})
}

func (h Handler) Encurtar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		responderCom(w, http.StatusMethodNotAllowed, headers{"Allow": "POST"})
		return
	}

	url, nova, err := url.BuscarOuCriarNovaUrl(extrairUrl(r))

	if err != nil {
		responderCom(w, http.StatusBadRequest, nil)
		return
	}

	var status int
	if nova {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	urlCurta := fmt.Sprintf("%s/r/%s", h.urlBase, url.Id)

	responderCom(w, status, headers{
		"Location": urlCurta,
		"Link":     fmt.Sprintf("<%s/api/stats/%s>; rel=\"stats\"", h.urlBase, url.Id),
	})

	h.logar("URL %s encurtada com sucesso para %s", url.Destino, urlCurta)
}

func (h Handler) Visualizar(w http.ResponseWriter, r *http.Request) {
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

func responderCom(w http.ResponseWriter, status int, headers headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
}

func responderComJSON(w http.ResponseWriter, resposta string) {
	responderCom(w, http.StatusOK, headers{"Content-Type": "application/json"})
	fmt.Fprintf(w, resposta)
}

func extrairUrl(r *http.Request) string {
	defer r.Body.Close()
	rawBody := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(rawBody)
	return string(rawBody)
}

func (h Handler) registrarEstatisticas(stats <-chan string) {
	for id := range stats {
		url.RegistrarClick(id)
		h.logar("Click registrado com sucesso para %s", id)
	}
}

func (h Handler) logar(formato string, valores ...interface{}) {
	if h.logLigado {
		log.Printf(fmt.Sprintf("%s\n", formato), valores...)
	}
}
