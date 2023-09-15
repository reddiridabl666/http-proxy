package api

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"time"
)

func getTlsConfig() (*tls.Config, error) {
	cert, err := os.ReadFile("https/ca.crt")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(cert)
	if !ok {
		return nil, errors.New("error appending certificate")
	}

	return &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}, nil
}

func getTlsTransport() (*http.Transport, error) {
	cfg, err := getTlsConfig()
	if err != nil {
		return nil, err
	}

	return &http.Transport{
		TLSHandshakeTimeout: time.Second * 5,
		TLSClientConfig:     cfg,
	}, nil
}
