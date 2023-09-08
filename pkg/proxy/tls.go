package proxy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func generateCertificate(host string) ([]byte, error) {
	cmd := exec.Command("./gen.sh", host)
	var out strings.Builder
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic("Certificate creation failure: " + err.Error())
		// return nil, err
	}

	res := []byte(out.String())

	cert, err := os.Create(fmt.Sprintf("certs/%s.crt", host))
	if err != nil {
		return nil, err
	}

	defer cert.Close()

	if len(res) == 0 {
		panic("Certificate creation failure: nothing returned from script")
	}

	i := 0
	var written int64 = 0

	for i < 5 {
		written, err = io.Copy(cert, bytes.NewReader(res))
		if err != nil {
			panic("Certificate creation failure: " + err.Error())
			// return nil, e
		}

		if written > 0 {
			break
		}

	}

	if written == 0 {
		panic("Certificate creation failure: nothing written")
	}

	return res, nil
}
