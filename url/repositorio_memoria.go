package url

type repositorioMemoria struct {
	urls map[string]*Url
	clicks map[string]int
}

func NovoRepositorioMemoria() *repositorioMemoria {
	return &repositorioMemoria{
		make(map[string]*Url),
		make(map[string]int),
	}
}

func (r *repositorioMemoria) IdExiste(id string) bool {
	_, existe := r.urls[id]
	return existe
}

func (r *repositorioMemoria) BuscarPorId(id string) *Url {
	return r.urls[id]
}

func (r *repositorioMemoria) BuscarPorUrl(url string) *Url {
	for _, u := range r.urls {
		if u.Destino == url {
			return u
		}
	}
	return nil
}

func (r *repositorioMemoria) Salvar(url Url) error {
	r.urls[url.Id] = &url
	return nil
}

func (r *repositorioMemoria) RegistrarClick(id string) {
	atual := 0

	c, existe := r.clicks[id]
	if existe {
		atual = c + 1
	} else {
		atual = 1
	}

	r.clicks[id] = atual
}

func (r *repositorioMemoria) BuscarClicks(id string) int {
	c, encontrado := r.clicks[id]

	if !encontrado {
		return -1
	}

	return c
}
