package repo

import "net/http"

type RequestSaver interface {
	Save(*http.Request) (string, error)
	Get(id int64) (*http.Request, error)
	List() ([]*http.Request, error)
}

type ResponseSaver interface {
	Save(*http.Response) (string, error)
	Get(id int64) (*http.Response, error)
	List() ([]*http.Response, error)
}

const kDatabase = "http-proxy"
