package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func loadCertificates() (map[string][]byte, error) {
	entries, err := os.ReadDir("certs")
	if err != nil {
		return nil, err
	}

	res := make(map[string][]byte, len(entries))

	for _, entry := range entries {
		host := strings.TrimSuffix(entry.Name(), ".crt")

		res[host], err = os.ReadFile("certs/" + entry.Name())
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func generateCertificate(host string) ([]byte, error) {
	cmd := exec.Command("./https/gen.sh", host)
	var out strings.Builder
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	res := []byte(out.String())

	cert, err := os.Create(fmt.Sprintf("certs/%s.crt", host))
	if err != nil {
		return nil, err
	}

	defer cert.Close()

	var written int64 = 0

	written, err = io.Copy(cert, bytes.NewReader(res))
	if err != nil {
		return nil, err
	}

	if written == 0 {
		return nil, errors.New("0 bytes written during certificate creation")
	}

	return res, nil
}
