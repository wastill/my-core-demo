package main

import (
	"github.com/wastill/my-core-demo/framework"
	"log"
	"net/http"
)

func main() {

	core := framework.NewCore()
	RegisterRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	//http.ServeFile(w, r, "index.html")
}
