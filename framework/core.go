package framework

import "net/http"

type Core struct {
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// TODO
	return
}

func NewCore() *Core {
	return &Core{}
}
