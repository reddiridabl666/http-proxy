package main

import (
	"fmt"

	"http-proxy/pkg/proxy"
	"http-proxy/pkg/repo"
	"http-proxy/pkg/server"
)

func main() {
	mongoConn, err := mongoConnect("root", "example", "mongo", 27017)
	if err != nil {
		fmt.Println(err)
		return
	}

	handler, err := proxy.NewHandler(repo.NewMongoRequestSaver(mongoConn), repo.NewMongoResponseSaver(mongoConn))
	if err != nil {
		fmt.Println(err)
		return
	}

	server.Run(8080, handler.Handle)
}
