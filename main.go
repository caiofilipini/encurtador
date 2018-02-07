package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/caiofilipini/encurtador/handlers"
	"github.com/caiofilipini/encurtador/url"
)

func main() {
	dominio := flag.String("d", "localhost", "domínio")
	porta := flag.Int("p", 8888, "porta")
	logLigado := flag.Bool("l", true, "log ligado/desligado")

	flag.Parse()

	urlBase := fmt.Sprintf("http://%s:%d", *dominio, *porta)

	url.ConfigurarRepositorio(url.NovoRepositorioMemoria())

	handler := handlers.NovoHandler(urlBase, *logLigado)
	http.HandleFunc("/r/", handler.Redirecionar)
	http.HandleFunc("/api/encurtar", handler.Encurtar)
	http.HandleFunc("/api/stats/", handler.Visualizar)

	if *logLigado {
		log.Printf("Iniciando servidor na porta %d...\n", *porta)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *porta), nil))
}
