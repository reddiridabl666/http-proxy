package proxy

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"http-proxy/pkg/utils"
)

type Handler struct {
	certs     map[string][]byte
	tlsConfig *tls.Config
	key       []byte
}

func NewHandler() (*Handler, error) {
	keyBytes, err := os.ReadFile("cert.key")
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair("ca.crt", "ca.key")
	if err != nil {
		return nil, err
	}

	return &Handler{
		certs:     make(map[string][]byte, 4),
		tlsConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		key:       keyBytes,
	}, nil
}

func (h *Handler) Handle(conn net.Conn) error {
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return err
	}

	return h.handleRequest(conn, req)
}

func (h *Handler) handleRequest(clientConn net.Conn, toProxy *http.Request) error {
	var hostConn net.Conn
	var err error

	if toProxy.Method == http.MethodConnect {
		err := h.handleHTTPS(clientConn, toProxy)
		if err != nil {
			return err
		}

		clientConn = tls.Server(clientConn, h.tlsConfig)

		fmt.Println("Reading the actual request")
		readBytes := []byte{}
		clientConn.SetReadDeadline(time.Now().Add(defaultTimeout))

		// _, err = clientConn.Read(readBytes)
		toProxy, err = http.ReadRequest(bufio.NewReader(clientConn))
		if err != nil {
			return err
		}

		fmt.Println(string(readBytes))

		hostConn, err = h.tlsConnect(getHost(toProxy.URL))
		if err != nil {
			return err
		}
	} else {
		hostConn, err = tcpConnect(getHost(toProxy.URL))
		if err != nil {
			return err
		}
	}

	req, err := NewRequest(toProxy)
	if err != nil {
		return err
	}

	resp, err := h.sendRequest(hostConn, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return utils.WriteResponse(resp, clientConn)
}

const defaultTimeout = time.Second * 5

func tcpConnect(host string) (net.Conn, error) {
	return net.DialTimeout("tcp", host, defaultTimeout)
}

func (h *Handler) tlsConnect(host string) (net.Conn, error) {
	cert, err := tls.X509KeyPair(h.certs[host], h.key)
	if err != nil {
		return nil, err
	}

	dialer := tls.Dialer{
		Config: &tls.Config{Certificates: []tls.Certificate{cert}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return dialer.DialContext(ctx, "tcp", host)
}

func (h *Handler) handleHTTPS(conn net.Conn, r *http.Request) error {
	_, err := conn.Write([]byte("HTTP/1.0 200 Connection established\n"))
	if err != nil {
		return err
	}

	_, exists := h.certs[r.Host]
	if !exists {
		cert, err := generateCertificate(r)
		if err != nil {
			return fmt.Errorf("error generating certificate: %v", err)
		}
		h.certs[r.Host] = cert
	}

	return nil
}

func NewRequest(r *http.Request) (*http.Request, error) {
	res, err := http.NewRequest(r.Method, r.URL.Path, nil)
	if err != nil {
		return nil, err
	}

	res.Header = r.Header
	res.Host = r.Host
	res.Header.Del("Proxy-Connection")
	res.Body = r.Body

	return res, nil
}

func (h *Handler) sendRequest(conn net.Conn, req *http.Request) (*http.Response, error) {
	bytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bytes))

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
