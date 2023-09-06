package proxy

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"

	"http-proxy/pkg/utils"
)

// var client http.Client

func Handler(w http.ResponseWriter, toProxy *http.Request) {
	req, err := NewRequest(toProxy)
	if err != nil {
		utils.WriteError(err, w)
		return
	}

	conn, err := net.Dial("tcp", req.Host+":80")
	if err != nil {
		utils.WriteError(err, w)
		return
	}

	bytes, err := httputil.DumpRequest(req, true)
	fmt.Println(string(bytes))

	if err != nil {
		utils.WriteError(err, w)
		return
	}

	_, err = conn.Write(bytes)
	if err != nil {
		utils.WriteError(err, w)
		return
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		utils.WriteError(err, w)
		return
	}

	err = utils.WriteResponse(resp, w)
	if err != nil {
		fmt.Printf("Error responding to client %s: %s\n", toProxy.RemoteAddr, err)
	}
}

func NewRequest(r *http.Request) (*http.Request, error) {
	res, err := http.NewRequest(r.Method, r.URL.Path, nil)
	if err != nil {
		return nil, err
	}

	res.Header = r.Header
	res.Host = r.Host
	res.Header.Del("Proxy-Connection")

	return res, nil
}
