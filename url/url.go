package url

import (
	"math/rand"
	"time"
)

const (
	tamanho    = 5
	caracteres = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_-+"
)

type Storage interface {
	IdExiste(id string) bool
	BuscarPorId(id string) *Url
	BuscarPorUrl(url string) *Url
	Salvar(url Url) error
}

type Url struct {
	Id      string
	Criacao time.Time
	destino string
}

var storage Storage

func init() {
	rand.Seed(time.Now().UnixNano())

	criarStorage()
}

func NovaUrl(destino string) *Url {
	if u := storage.BuscarPorUrl(destino); u != nil {
		return u
	}

	url := Url{gerarId(), time.Now(), destino}
	storage.Salvar(url)
	return &url
}

func Buscar(id string) *Url {
	return storage.BuscarPorId(id)
}

func criarStorage() {
	if storage == nil {
		storage = &storageMemoria{make(map[string]*Url)}
	}
}

func gerarId() string {
	novoId := func() string {
		id := make([]byte, tamanho, tamanho)
		for i := range id {
			id[i] = caracteres[rand.Intn(len(caracteres))]
		}
		return string(id)
	}

	for {
		if id := novoId(); !storage.IdExiste(id) {
			return id
		}
	}
}
