package proxy

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"http-proxy/pkg/utils"
)

func Handle(conn net.Conn) error {
	reader := bufio.NewReader(conn)

	req, err := http.ReadRequest(reader)
	if err != nil {
		return err
	}

	return handleRequest(conn, req)
}

func handleRequest(conn net.Conn, toProxy *http.Request) error {
	req, err := NewRequest(toProxy)
	if err != nil {
		return err
	}

	resp, err := sendRequest(req, toProxy.URL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return utils.WriteResponse(resp, conn)
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

func sendRequest(req *http.Request, originalURL *url.URL) (*http.Response, error) {
	bytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bytes))

	conn, err := net.DialTimeout("tcp", getHost(originalURL), time.Second*5)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(bytes)
	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(conn), req)
}

func getHost(url *url.URL) string {
	return url.Hostname() + ":" + getPort(url)
}

func getPort(url *url.URL) string {
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
