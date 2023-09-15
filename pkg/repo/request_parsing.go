package repo

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func toRequest(data *RequestData) (*http.Request, error) {
	res, err := http.NewRequest(data.Method, "http://"+data.Host+data.Path, nil)
	if err != nil {
		return nil, err
	}

	res.Host = data.Host

	res.Header = fromBson(data.Headers)
	addCookies(res, data.Cookies)
	// res.Header.Set("Cookie", encodeCookies(data.Cookies))
	res.URL.RawQuery = makeQuery(fromBson(data.GetParams))
	res.Body = getBody(res, data)

	return res, nil
}

func addCookies(req *http.Request, cookies map[string]string) {
	for key, value := range cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: value})
	}
}

// func encodeCookies(cookies map[string]string) string {
// 	if len(cookies) == 0 {
// 		return ""
// 	}

// 	cookieString := ""

// 	for key, value := range cookies {
// 		cookieString += fmt.Sprintf("%s=%s; ", key, value)
// 	}
// 	return strings.TrimSuffix(cookieString, "; ")
// }

func makeQuery(values map[string][]string) string {
	return url.Values(values).Encode()
}

func getBody(req *http.Request, data *RequestData) io.ReadCloser {
	if req.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		var params url.Values = fromBson(data.PostParams)
		return io.NopCloser(strings.NewReader(params.Encode()))
	}

	return io.NopCloser(strings.NewReader(data.Body))
}

func parseHeaders(headers http.Header) bson.M {
	res := toBson(headers)
	delete(res, "Cookie")
	return res
}

func parseQuery(input *url.URL) bson.M {
	input.RawQuery = strings.ReplaceAll(input.RawQuery, ";", "&")
	res := toBson(input.Query())
	return res
}

func toBson(values map[string][]string) bson.M {
	res := make(bson.M, len(values))

	for key, value := range values {
		if len(value) == 1 {
			res[key] = value[0]
		} else {
			res[key] = value
		}
	}
	return res
}

func fromBson(values bson.M) map[string][]string {
	res := make(map[string][]string, len(values))

	for key, value := range values {
		str, ok := value.(string)
		if ok {
			res[key] = []string{str}
		}

		arr, ok := value.(bson.A)
		if !ok {
			continue
		}

		res[key] = make([]string, len(arr))
		for _, elem := range arr {
			str, ok := elem.(string)
			if !ok {
				continue
			}
			res[key] = append(res[key], str)
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

func parsePostParams(req *http.Request) (bson.M, error) {
	if req.Body == nil {
		return nil, nil
	}

	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	return toBson(req.PostForm), nil
}
