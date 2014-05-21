package url

import (
	"math/rand"
	"time"
)

const (
	tamanho  = 5
	simbolos = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-+"
)

type Repositorio interface {
	IdExiste(id string) bool
	BuscarPorId(id string) *Url
	BuscarPorUrl(url string) *Url
	Salvar(url Url) error
}

type Url struct {
	Id      string
	Criacao time.Time
	Destino string
}

var repo Repositorio

func init() {
	rand.Seed(time.Now().UnixNano())

	criarRepositorio()
}

func NovaUrl(destino string) *Url {
	if u := repo.BuscarPorUrl(destino); u != nil {
		return u
	}

	url := Url{gerarId(), time.Now(), destino}
	repo.Salvar(url)
	return &url
}

func Buscar(id string) *Url {
	return repo.BuscarPorId(id)
}

func criarRepositorio() {
	if repo == nil {
		repo = &repositorioMemoria{make(map[string]*Url)}
	}
}

func gerarId() string {
	novoId := func() string {
		id := make([]byte, tamanho, tamanho)
		for i := range id {
			id[i] = simbolos[rand.Intn(len(simbolos))]
		}
		return string(id)
	}

	for {
		if id := novoId(); !repo.IdExiste(id) {
			return id
		}
	}
}
