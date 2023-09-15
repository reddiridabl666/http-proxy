package repo

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestSaver interface {
	Save(*http.Request) (string, error)
	Get(id string) (*RequestData, error)
	GetEncoded(id string) (*http.Request, error)
	List(limit int64) ([]*RequestData, error)
}

type ResponseSaver interface {
	Save(requestId string, resp *http.Response) (string, error)
	Get(id string) (*ResponseData, error)
	GetByRequest(requestId string) (*ResponseData, error)
	List(limit int64) ([]*ResponseData, error)
}

const kDatabase = "http-proxy"

type RequestData struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Scheme     string             `json:"scheme"`
	Host       string             `json:"host"`
	Method     string             `json:"method"`
	Path       string             `json:"path"`
	Cookies    map[string]string  `json:"cookies"`
	Body       string             `json:"body,omitempty" bson:"body,omitempty"`
	Headers    bson.M             `json:"headers"`
	GetParams  bson.M             `json:"get_params" bson:"get_params"`
	PostParams bson.M             `json:"post_params" bson:"post_params"`
}

type ResponseData struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	RequestId primitive.ObjectID `json:"request_id" bson:"request_id"`
	Code      int                `json:"code"`
	Message   string             `json:"message"`
	Body      string             `json:"body,omitempty" bson:"body,omitempty"`
	Headers   bson.M             `json:"headers"`
}
