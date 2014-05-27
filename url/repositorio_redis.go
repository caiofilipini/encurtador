package url

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type RepositorioRedis struct {
	Conn redis.Conn
}

func NovoRepositorioRedis(c redis.Conn) *RepositorioRedis {
	return &RepositorioRedis{c}
}

func (r *RepositorioRedis) IdExiste(id string) bool {
	existe, err := redis.Bool(r.Conn.Do("EXISTS", chave(id)))
	tratar(err)
	return existe
}

func (r *RepositorioRedis) BuscarPorId(id string) *Url {
	if res := r.get(chave(id)); res != nil {
		b, err := redis.Bytes(res, nil)
		tratar(err)

		url := &Url{}
		err = json.Unmarshal(b, url)
		tratar(err)

		return url
	}
	return nil
}

func (r *RepositorioRedis) BuscarPorUrl(url string) *Url {
	if s := r.get(urls(url)); s != nil {
		id, err := redis.String(s, nil)
		tratar(err)

		return r.BuscarPorId(id)
	}
	return nil
}

func (r *RepositorioRedis) Salvar(url Url) error {
	json, err := json.Marshal(url)
	tratar(err)

	r.Conn.Send("MULTI")
	r.Conn.Send("SET", chave(url.Id), string(json))
	r.Conn.Send("SET", urls(url.Destino), url.Id)
	_, err = r.Conn.Do("EXEC")
	return err
}

func (r *RepositorioRedis) RegistrarClick(id string) {
	_, err := r.Conn.Do("INCR", clicks(id))
	tratar(err)
}

func (r *RepositorioRedis) BuscarClicks(id string) int {
	c, err := redis.Int(r.Conn.Do("GET", clicks(id)))
	tratar(err)
	return c
}

func (r *RepositorioRedis) get(chave string) interface{} {
	res, err := r.Conn.Do("GET", chave)
	tratar(err)
	return res
}

func chave(id string) string {
	return fmt.Sprintf("encurtador:urls:%s", id)
}

func clicks(id string) string {
	return fmt.Sprintf("encurtador:clicks:%s", id)
}

func urls(url string) string {
	return fmt.Sprintf("encurtador:urls:ids:%s", hash(url))
}

func hash(url string) string {
	h := md5.New()
	h.Write([]byte(url))
	return hex.EncodeToString(h.Sum(nil))
}

func tratar(err error) {
	if err != nil {
		panic(err)
	}
}
