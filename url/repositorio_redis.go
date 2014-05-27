package url

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type repositorioRedis struct {
	conn redis.Conn
}

func NovoRepositorioRedis(c redis.Conn) *repositorioRedis {
	return &repositorioRedis{c}
}

func (r *repositorioRedis) IdExiste(id string) bool {
	existe, err := redis.Bool(r.conn.Do("EXISTS", chave(id)))
	tratar(err)
	return existe
}

func (r *repositorioRedis) BuscarPorId(id string) *Url {
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

func (r *repositorioRedis) BuscarPorUrl(url string) *Url {
	if s := r.get(urls(url)); s != nil {
		id, err := redis.String(s, nil)
		tratar(err)

		return r.BuscarPorId(id)
	}
	return nil
}

func (r *repositorioRedis) Salvar(url Url) error {
	json, err := json.Marshal(url)
	tratar(err)

	r.conn.Send("MULTI")
	r.conn.Send("SET", chave(url.Id), string(json))
	r.conn.Send("SET", urls(url.Destino), url.Id)
	_, err = r.conn.Do("EXEC")
	return err
}

func (r *repositorioRedis) RegistrarClick(id string) {
	_, err := r.conn.Do("INCR", clicks(id))
	tratar(err)
}

func (r *repositorioRedis) BuscarClicks(id string) int {
	c, err := redis.Int(r.conn.Do("GET", clicks(id)))
	tratar(err)
	return c
}

func (r *repositorioRedis) get(chave string) interface{} {
	res, err := r.conn.Do("GET", chave)
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
