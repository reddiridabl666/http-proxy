package proxy

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"http-proxy/pkg/repo"
	"http-proxy/pkg/utils"
)

type Handler struct {
	certs         map[string][]byte
	mutex         sync.Mutex
	key           []byte
	requestSaver  repo.RequestSaver
	responseSaver repo.ResponseSaver
}

func NewHandler(req repo.RequestSaver, resp repo.ResponseSaver) (*Handler, error) {
	keyBytes, err := os.ReadFile("https/cert.key")
	if err != nil {
		return nil, err
	}

	certs, err := loadCertificates()
	if err != nil {
		return nil, err
	}

	return &Handler{
		certs:         certs,
		key:           keyBytes,
		requestSaver:  req,
		responseSaver: resp,
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

	// utils.PrintRequest(toProxy)

	host := toProxy.URL.Hostname()
	port := utils.GetPort(toProxy.URL)

	if toProxy.Method == http.MethodConnect {
		clientConn, err = h.tlsUpgrade(clientConn, host)
		if err != nil {
			return err
		}

		// fmt.Println("Reading the actual request")
		toProxy, err = http.ReadRequest(bufio.NewReader(clientConn))
		if err != nil {
			return err
		}

		// fmt.Println("Connecting to host: " + host)
		hostConn, err = tlsConnect(host, port)
		if err != nil {
			return err
		}
	} else {
		hostConn, err = tcpConnect(host, port)
		if err != nil {
			return err
		}
	}

	prepareRequest(toProxy)
	requestId, err := h.requestSaver.Save(toProxy)
	if err != nil {
		return err
	}

	// fmt.Println("Proxying request to host: " + host + "\n")
	resp, err := sendRequest(hostConn, toProxy)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	h.responseSaver.Save(requestId, resp)

	return utils.WriteResponse(resp, clientConn)
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
