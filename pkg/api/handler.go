package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"http-proxy/pkg/repo"
	"http-proxy/pkg/utils"

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
		client:    http.DefaultClient,
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

	resp, err := h.client.Do(req)
	if err != nil {
		utils.HttpError(errors.New("Error resending request: "+err.Error()), w)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.HttpError(errors.New("Error resending request: "+err.Error()), w)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(bytes)
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
