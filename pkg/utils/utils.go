package utils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func WriteError(err error, conn net.Conn) {
	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("Proxy error:" + err.Error())),
	}

	WriteResponse(resp, conn)
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
