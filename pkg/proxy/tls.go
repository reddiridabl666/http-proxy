package proxy

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func generateCertificate(r *http.Request) ([]byte, error) {
	cmd := exec.Command("./gen_cert.sh", r.Host, strconv.Itoa(rand.Intn(1000000000000)))
	var out strings.Builder
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	res := []byte(out.String())

	cert, err := os.Create(fmt.Sprintf("certs/%s.crt", r.Host))
	if err != nil {
		return nil, err
	}

	defer cert.Close()

	_, err = io.Copy(cert, strings.NewReader(out.String()))
	if err != nil {
		cert.Close()
		return nil, err
	}

	return res, nil
}
