package url

type repositorioMemoria struct {
	urls map[string]*Url
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
