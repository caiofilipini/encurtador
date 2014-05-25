package url

import (
	"math/rand"
	"net/url"
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
	RegistrarClick(id string)
	BuscarClicks(id string) int
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

func RegistrarClick(id string) {
	repo.RegistrarClick(id)
}

func BuscarClicks(id string) int {
	return repo.BuscarClicks(id)
}

func BuscarOuCriarNovaUrl(destino string) (*Url, error) {
	if u := repo.BuscarPorUrl(destino); u != nil {
		return u, nil
	}

	if _, err := url.ParseRequestURI(destino); err != nil {
		return nil, err
	}

	url := Url{gerarId(), time.Now(), destino}
	repo.Salvar(url)
	return &url, nil
}

func Buscar(id string) *Url {
	return repo.BuscarPorId(id)
}

func criarRepositorio() {
	if repo == nil {
		repo = NovoRepositorioMemoria()
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
