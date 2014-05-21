package main

import (
	"fmt"
	"github.com/caiofilipini/encurtador/url"
)

func main() {
	u1 := url.NovaUrl("http://google.com")
	u2 := url.NovaUrl("http://google.com/1")
	u3 := url.NovaUrl("http://google.com/2")
	fmt.Println(u1.Id)
	fmt.Println(u2.Id)
	fmt.Println(u3.Id)
	fmt.Println(url.Buscar(u1.Id))
	fmt.Println(url.Buscar(u2.Id))
	fmt.Println(url.Buscar(u3.Id))
}
