package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gibbs/godanuk/tool"
	"github.com/gorilla/mux"
)

var (
	listenAddr = flag.String("listen-address", "0.0.0.0:8084", "The address to listen on")
)

type Tool struct {
	Success bool   `json:"success"`
	Command string `json:"command"`
	Output  string `json:"output"`
}

func main() {
	r := mux.NewRouter()

	// Healthcheck route
	r.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	// Dig tool
	r.HandleFunc("/tools/dig", func(w http.ResponseWriter, r *http.Request) {
		request := tool.DigPayloadInterface()
		decodeRequestPayload(w, r, request)
		responseJSON(w, 200, tool.DigHandler(request))
	}).Methods(http.MethodPost)

	// Mkpasswd tool
	r.HandleFunc("/tools/mkpasswd", func(w http.ResponseWriter, r *http.Request) {
		request := tool.MkpasswdPayloadInterface()
		decodeRequestPayload(w, r, &request)
		responseJSON(w, 200, tool.MkpasswdHandler(request))
	}).Methods(http.MethodPost)

	// Pwgen
	r.HandleFunc("/tools/pwgen", func(w http.ResponseWriter, r *http.Request) {
		request := tool.PwgenPayloadInterface()
		decodeRequestPayload(w, r, &request)
		responseJSON(w, 200, tool.PwgenHandler(request))
	}).Methods(http.MethodPost)

	// UUID gen
	r.HandleFunc("/tools/uuidgen", func(w http.ResponseWriter, r *http.Request) {
		request := tool.UuidgenPayloadInterface()
		decodeRequestPayload(w, r, &request)
		responseJSON(w, 200, tool.UuidgenHandler(request))
	}).Methods(http.MethodPost)

	// Setup the server
	srv := &http.Server{
		Addr:         *listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Shutdown signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// 15 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Shutdown
	srv.Shutdown(ctx)
	os.Exit(0)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := json.NewEncoder(w).Encode(map[string]bool{"success": false})

	responseJSON(w, http.StatusBadGateway, response)
}

func decodeRequestPayload(w http.ResponseWriter, r *http.Request, payload interface{}) {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(payload); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}
}

func responseJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(response))
}

func responseError(w http.ResponseWriter, status int, message string) {
	responseJSON(w, status, map[string]string{"success": "false", "error": message})
}
