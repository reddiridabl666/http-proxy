package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

func WriteError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Proxy error: " + err.Error() + "\n"))
}

func PrintRequest(r *http.Request) {
	bytes, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(bytes))
}

func WriteResponse(resp *http.Response, w http.ResponseWriter) error {
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
