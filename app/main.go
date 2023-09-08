package main

import (
	"fmt"

	"http-proxy/pkg/proxy"
	"http-proxy/pkg/server"
)

func main() {
	handler, err := proxy.NewHandler()
	if err != nil {
		fmt.Println(err)
		return
	}

	server.Run(8080, handler.Handle)
}
