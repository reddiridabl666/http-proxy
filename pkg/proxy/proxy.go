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
	"sync"
	"time"

	"http-proxy/pkg/utils"
)

type Handler struct {
	certs map[string][]byte
	mutex sync.Mutex
	// tlsConfig *tls.Config
	key []byte
}

func NewHandler() (*Handler, error) {
	keyBytes, err := os.ReadFile("cert.key")
	if err != nil {
		return nil, err
	}

	return &Handler{
		certs: make(map[string][]byte, 4),
		key:   keyBytes,
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

	utils.PrintRequest(toProxy)

	host := toProxy.URL.Hostname()
	port := getPort(toProxy.URL)

	if toProxy.Method == http.MethodConnect {
		clientConn, err = h.tlsUpgrade(clientConn, host)
		if err != nil {
			return err
		}

		fmt.Println("Reading the actual request")
		toProxy, err = http.ReadRequest(bufio.NewReader(clientConn))
		if err != nil {
			return err
		}

		fmt.Println("Connecting to host: " + host)
		hostConn, err = h.tlsConnect(host, port)
		if err != nil {
			return err
		}
	} else {
		hostConn, err = tcpConnect(host, port)
		if err != nil {
			return err
		}
	}

	req, err := NewRequest(toProxy)
	if err != nil {
		return err
	}

	fmt.Println("Proxying request to host: " + host + "\n")
	resp, err := h.sendRequest(hostConn, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return utils.WriteResponse(resp, clientConn)
}

const defaultTimeout = time.Second * 5

func tcpConnect(host, port string) (net.Conn, error) {
	return net.DialTimeout("tcp", host+":"+port, defaultTimeout)
}

func (h *Handler) tlsConnect(host, port string) (net.Conn, error) {
	dialer := tls.Dialer{}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	conn, err := dialer.DialContext(ctx, "tcp", host+":"+port)

	return conn, err
}

func (h *Handler) getTlsConfig(host string) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(h.certs[host], h.key)
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

func (h *Handler) tlsUpgrade(clientConn net.Conn, host string) (net.Conn, error) {
	_, err := clientConn.Write([]byte("HTTP/1.0 200 Connection established\n\n"))
	if err != nil {
		return nil, err
	}

	err = h.generateCertificate(host)
	if err != nil {
		return nil, err
	}

	cfg, err := h.getTlsConfig(host)
	if err != nil {
		return nil, err
	}

	tlsConn := tls.Server(clientConn, cfg)
	clientConn.SetReadDeadline(time.Now().Add(defaultTimeout))

	return tlsConn, nil
}

func (h *Handler) generateCertificate(host string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, exists := h.certs[host]
	if !exists {
		fmt.Printf("Generating certificate for %s\n", host)
		cert, err := generateCertificate(host)
		if err != nil {
			return fmt.Errorf("error generating certificate: %v", err)
		}
		h.certs[host] = cert
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
	// fmt.Println(string(bytes))

	_, err = conn.Write(bytes)
	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(conn), req)
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
