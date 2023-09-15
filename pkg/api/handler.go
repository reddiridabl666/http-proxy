package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"

	"http-proxy/pkg/repo"
	"http-proxy/pkg/utils"
	"http-proxy/pkg/xxe"

	"github.com/gorilla/mux"
)

type Handler struct {
	requests  repo.RequestSaver
	responses repo.ResponseSaver
	client    *http.Client
}

func NewHandler(req repo.RequestSaver, resp repo.ResponseSaver) *Handler {
	return &Handler{
		requests:  req,
		responses: resp,
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request) {
	req, err := h.requests.Get(mux.Vars(r)["id"])
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	w.Write(bytes)
}

const kDefaultListSize = 5

func (h *Handler) ListRequests(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = kDefaultListSize
	}

	requests, err := h.requests.List(limit)
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	bytes, err := json.Marshal(requests)
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	w.Write(bytes)
}

func (h *Handler) RepeatRequest(w http.ResponseWriter, r *http.Request) {
	req, err := h.requests.GetEncoded(mux.Vars(r)["id"])
	if err != nil {
		utils.HttpError(errors.New("Error getting request: "+err.Error()), w)
		return
	}

	// conn, err := utils.TcpConnect(req.Host, utils.GetPort(req.URL))
	resp, err := h.client.Do(req)
	if err != nil {
		utils.HttpError(errors.New("Error resending request: "+err.Error()), w)
		return
	}

	// resp, err := utils.SendRequest(conn, req)
	// if err != nil {
	// 	utils.HttpError(errors.New("Error resending request: "+err.Error()), w)
	// 	return
	// }
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.HttpError(errors.New("Error resending request: "+err.Error()), w)
		return
	}

	for key, values := range resp.Header {
		for _, elem := range values {
			w.Header().Add(key, elem)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(bytes)
}

func (h *Handler) DumpRequest(w http.ResponseWriter, r *http.Request) {
	req, err := h.requests.GetEncoded(mux.Vars(r)["id"])
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	bytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	w.Write(bytes)
}

func (h *Handler) ScanRequest(w http.ResponseWriter, r *http.Request) {
	req, err := h.requests.GetEncoded(mux.Vars(r)["id"])
	if err != nil {
		utils.HttpError(errors.New("Error getting request: "+err.Error()), w)
		return
	}

	hadXML, err := xxe.AddVulnerability(req)
	if err != nil {
		utils.HttpError(errors.New("Error adding vulnerability to request: "+err.Error()), w)
		return
	}

	resp, err := h.client.Do(req)
	if err != nil {
		utils.HttpError(errors.New("Error resending request: "+err.Error()), w)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.HttpError(errors.New("Error reading response: "+err.Error()), w)
		return
	}

	if xxe.IsVulnerable(body) {
		body = []byte("Request vulnerable, response:\n" + string(body))
	} else {
		body = []byte("Request is not vulnerable\n")
	}

	w.Write(body)
}

func (h *Handler) GetResponse(w http.ResponseWriter, r *http.Request) {
	req, err := h.responses.Get(mux.Vars(r)["id"])
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	w.Write(bytes)
}

func (h *Handler) GetRequestResponse(w http.ResponseWriter, r *http.Request) {
	req, err := h.responses.GetByRequest(mux.Vars(r)["id"])
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	w.Write(bytes)
}

func (h *Handler) ListResponses(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = kDefaultListSize
	}

	requests, err := h.responses.List(limit)
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	bytes, err := json.Marshal(requests)
	if err != nil {
		utils.HttpError(err, w)
		return
	}

	w.Write(bytes)
}
