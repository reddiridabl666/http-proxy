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

// var client http.Client

func Handler(w http.ResponseWriter, toProxy *http.Request) {
	req, err := NewRequest(toProxy)
	if err != nil {
		utils.WriteError(err, w)
		return
	}

	bytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		utils.WriteError(err, w)
		return
	}
	fmt.Println(string(bytes))

	conn, err := net.DialTimeout("tcp", getHost(toProxy.URL), time.Second*5)
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

	defer resp.Body.Close()

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
