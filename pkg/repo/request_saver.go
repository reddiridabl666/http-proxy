package repo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const kRequests = "requests"

type MongoRequestSaver struct {
	requests *mongo.Collection
}

func NewMongoRequestSaver(conn *mongo.Client) RequestSaver {
	return &MongoRequestSaver{
		requests: conn.Database(kDatabase).Collection(kRequests),
	}
}

func (s *MongoRequestSaver) Save(req *http.Request) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	req.Body = io.NopCloser(bytes.NewReader(body))

	value := bson.M{
		"method":     req.Method,
		"host":       req.Host,
		"path":       req.URL.Path,
		"get_params": parseQuery(req.URL),
		"headers":    parseHeaders(req.Header),
		"cookies":    parseCookies(req.Cookies()),
	}

	postParams, err := parsePostParams(req)
	if err != nil {
		return "", err
	}
	if postParams != nil {
		value["post_params"] = postParams
		req.Body = io.NopCloser(bytes.NewReader(body))
	} else {
		if req.Body != nil {
			value["body"] = string(body)
		}
	}

	res, err := s.requests.InsertOne(context.Background(), value)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *MongoRequestSaver) Get(id int64) (*http.Request, error) {
	return nil, nil
}

func (s *MongoRequestSaver) List() ([]*http.Request, error) {
	return nil, nil
}

func parseHeaders(headers http.Header) bson.M {
	res := toBson(headers)
	delete(res, "Cookie")
	return res
}

func parseQuery(input *url.URL) bson.M {
	input.RawQuery = strings.ReplaceAll(input.RawQuery, ";", "&")
	res := toBson(input.Query())
	delete(res, "Cookie")
	return res
}

func toBson(values map[string][]string) bson.M {
	res := bson.M{}
	for key, value := range values {
		if len(value) == 1 {
			res[key] = value[0]
		} else {
			res[key] = value
		}
	}
	return res
}

func parseCookies(cookies []*http.Cookie) map[string]string {
	res := make(map[string]string, 4)

	for _, cookie := range cookies {
		res[cookie.Name] = cookie.Value
	}

	return res
}

const maxMemory = 5 * 1024 * 1024 * 1024

func parsePostParams(req *http.Request) (bson.M, error) {
	if req.Body == nil {
		return nil, nil
	}

	err := func() error {
		if strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data") {
			return req.ParseMultipartForm(maxMemory)
		}
		return req.ParseForm()
	}()
	if err != nil {
		fmt.Println()
		return nil, err
	}

	return toBson(req.PostForm), nil
}
