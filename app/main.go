package main

import (
	"fmt"
	"net/http"

	"http-proxy/pkg/proxy"
)

func main() {
	http.HandleFunc("/", proxy.Handler)

	fmt.Println("Listening at port 8080...")
	http.ListenAndServe(":8080", nil)
}
