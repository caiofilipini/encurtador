package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/caiofilipini/encurtador/url"
	"github.com/garyburd/redigo/redis"
)

var (
	logLigado *bool
	porta     *int
	urlBase   string
)

func init() {
	dominio := flag.String("d", "localhost", "domínio")
	porta = flag.Int("p", 8888, "porta")
	logLigado = flag.Bool("l", true, "log ligado/desligado")

	flag.Parse()

	urlBase = fmt.Sprintf("http://%s:%d", *dominio, *porta)
}

type Headers map[string]string

type Redirecionador struct {
	stats chan string
}

func (r *Redirecionador) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	buscarUrlEExecutar(w, req, func(url *url.Url) {
		http.Redirect(w, req, url.Destino, http.StatusMovedPermanently)
		r.stats <- url.Id
	})
}

func Encurtador(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		responderCom(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})
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

	urlCurta := fmt.Sprintf("%s/r/%s", urlBase, url.Id)

	responderCom(w, status, Headers{
		"Location": urlCurta,
		"Link":     fmt.Sprintf("<%s/api/stats/%s>; rel=\"stats\"", urlBase, url.Id),
	})

	logar("URL %s encurtada com sucesso para %s.", url.Destino, urlCurta)
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

func extrairUrl(r *http.Request) string {
	rawBody := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(rawBody)
	return string(rawBody)
}

func registrarEstatisticas(stats <-chan string) {
	for id := range stats {
		url.RegistrarClick(id)
		logar("Click registrado com sucesso para %s.", id)
	}
}

func logar(formato string, valores ...interface{}) {
	if *logLigado {
		log.Printf(fmt.Sprintf("%s\n", formato), valores...)
	}
}

func main() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	url.ConfigurarRepositorio(url.NovoRepositorioRedis(conn))

	stats := make(chan string)
	defer close(stats)
	go registrarEstatisticas(stats)

	http.Handle("/r/", &Redirecionador{stats})
	http.HandleFunc("/api/encurtar", Encurtador)
	http.HandleFunc("/api/stats/", Visualizador)

	logar("Iniciando servidor na porta %d...", *porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *porta), nil))
}
