package utils

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func WriteError(err error, conn net.Conn) {
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("Proxy error:" + err.Error())),
	}

	WriteResponse(resp, conn)
}

func HttpError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func PrintRequest(r *http.Request) {
	bytes, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(bytes))
}

func PrintResponse(r *http.Response, body bool) {
	bytes, _ := httputil.DumpResponse(r, body)
	fmt.Println(string(bytes))
}

func WriteResponse(resp *http.Response, conn net.Conn) error {
	bytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}

	_, err = conn.Write(bytes)
	return err
}

func GetPort(url *url.URL) string {
	port := url.Port()

	if port == "" {
		switch url.Scheme {
		case "https":
			port = "443"
		default:
			port = "80"
		}
	}

	return port
}

func SendRequest(conn net.Conn, req *http.Request) (*http.Response, error) {
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

const DefaultTimeout = time.Second * 10

func TcpConnect(host, port string) (net.Conn, error) {
	return net.DialTimeout("tcp", host+":"+port, DefaultTimeout)
}

func TlsConnect(host, port string) (net.Conn, error) {
	dialer := tls.Dialer{}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	conn, err := dialer.DialContext(ctx, "tcp", host+":"+port)

	return conn, err
}
