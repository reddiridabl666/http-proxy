package proxy

import (
	"bufio"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

func prepareRequest(r *http.Request) {
	r.URL.Host = ""
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Accept-Encoding")
}

func sendRequest(conn net.Conn, req *http.Request) (*http.Response, error) {
	bytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(bytes)
	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(conn), req)
}

const defaultTimeout = time.Second * 20

func tcpConnect(host, port string) (net.Conn, error) {
	return net.DialTimeout("tcp", host+":"+port, defaultTimeout)
}

func tlsConnect(host, port string) (net.Conn, error) {
	dialer := tls.Dialer{}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	conn, err := dialer.DialContext(ctx, "tcp", host+":"+port)

	return conn, err
}
