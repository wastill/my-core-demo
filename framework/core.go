package framework

import (
	"log"
	"net/http"
)

type Core struct {
	router map[string]ControllerHandler
}

func NewCore() *Core {
	return &Core{router: make(map[string]ControllerHandler)}
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// TODO
	log.Printf("core.serveHTTP")
	ctx := NewContext(request, response)
	if router, ok := c.router["foo"]; ok {
		router(ctx)
	}

}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}
