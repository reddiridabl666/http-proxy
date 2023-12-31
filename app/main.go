package main

import (
	"fmt"
	"net/http"

	"http-proxy/pkg/api"
	"http-proxy/pkg/proxy"
	"http-proxy/pkg/repo"
	"http-proxy/pkg/server"

	"github.com/gorilla/mux"
)

func main() {
	mongoConn, err := mongoConnect("root", "example", "mongo", 27017)
	if err != nil {
		fmt.Println(err)
		return
	}

	requests := repo.NewMongoRequestSaver(mongoConn)
	responses := repo.NewMongoResponseSaver(mongoConn)

	handler, err := proxy.NewHandler(requests, responses)
	if err != nil {
		fmt.Println(err)
		return
	}

	go startHttpApi(requests, responses)

	server.Run(8080, handler.Handle)
}

func startHttpApi(req repo.RequestSaver, resp repo.ResponseSaver) {
	router := mux.NewRouter()

	handler, err := api.NewHandler(req, resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	router.HandleFunc("/requests", handler.ListRequests)
	router.HandleFunc("/requests/{id}", handler.GetRequest)
	router.HandleFunc("/repeat/{id}", handler.RepeatRequest)
	router.HandleFunc("/scan/{id}", handler.ScanRequest)
	router.HandleFunc("/requests/{id}/dump", handler.DumpRequest)

	router.HandleFunc("/responses", handler.ListResponses)
	router.HandleFunc("/responses/{id}", handler.GetResponse)
	router.HandleFunc("/requests/{id}/response", handler.GetRequestResponse)

	fmt.Println("Api listening at port 8000...")

	http.ListenAndServe(":8000", router)
}
