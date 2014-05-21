package url

type repositorioMemoria struct {
	urls map[string]*Url
}

type predicado func(u *Url) bool

func (r *repositorioMemoria) IdExiste(id string) bool {
	_, existe := r.urls[id]
	return existe
}

func (r *repositorioMemoria) BuscarPorId(id string) *Url {
	return r.buscar(func(u *Url) bool { return u.Id == id })
}

func (r *repositorioMemoria) BuscarPorUrl(url string) *Url {
	return r.buscar(func(u *Url) bool { return u.Destino == url })
}

func (r *repositorioMemoria) Salvar(url Url) error {
	r.urls[url.Id] = &url
	return nil
}

func (r *repositorioMemoria) buscar(p predicado) *Url {
	for _, u := range r.urls {
		if p(u) {
			return u
		}
	}
	return nil
}
