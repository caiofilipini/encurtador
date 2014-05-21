package url

type storageMemoria struct {
	urls map[string]*Url
}

type predicado func(u *Url) bool

func (s *storageMemoria) IdExiste(id string) bool {
	_, existe := s.urls[id]
	return existe
}

func (s *storageMemoria) BuscarPorId(id string) *Url {
	return s.buscar(func(u *Url) bool { return u.Id == id })
}

func (s *storageMemoria) BuscarPorUrl(url string) *Url {
	return s.buscar(func(u *Url) bool { return u.destino == url })
}

func (s *storageMemoria) Salvar(url Url) error {
	s.urls[url.Id] = &url
	return nil
}

func (s *storageMemoria) buscar(p predicado) *Url {
	for _, u := range s.urls {
		if p(u) {
			return u
		}
	}
	return nil
}
