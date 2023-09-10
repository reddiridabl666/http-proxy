package repo

import "net/http"

type RequestSaver interface {
	Save(*http.Request) (string, error)
	Get(id string) (*http.Request, error)
	List(limit int64) ([]*http.Request, error)
}

type ResponseSaver interface {
	Save(*http.Response) (string, error)
	Get(id string) (*http.Response, error)
	List(limit int64) ([]*http.Response, error)
}

const kDatabase = "http-proxy"
