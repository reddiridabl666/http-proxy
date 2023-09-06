package main

import (
	"fmt"
	"net"

	"http-proxy/pkg/proxy"
)

func main() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: 8080,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening at port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = proxy.Handle(conn)
		if err != nil {
			fmt.Println(err)
		}
	}
}
