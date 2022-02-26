package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

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

type mkpasswdPayload struct {
	Input  string `json:"input"`
	Salt   string `json:"salt"`
	Rounds uint32 `json:"rounds"`
	Method string `json:"method"`
}

type pwgenPayload struct {
	NoNumerals   bool   `json:"no-numerals"`
	NoCapitalize bool   `json:"no-capitalize"`
	Ambiguous    bool   `json:"ambiguous"`
	Capitalize   bool   `json:"capitalize"`
	NumPasswords uint16 `json:"num-passwords"`
	Numerals     bool   `json:"numerals"`
	RemoveChars  string `json:"remove-chars"`
	Secure       bool   `json:"secure"`
	NoVowels     bool   `json:"no-vowels"`
	Symbols      bool   `json:"symbols"`
	Length       uint16 `json:"length"`
}

type uuidgenPayload struct {
	Random bool `json:"random"`
	Time   bool `json:"time"`
}

func main() {
	r := mux.NewRouter()

	// Healthcheck route
	r.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	// Route handler
	r.HandleFunc("/tools/mkpasswd", mkpasswdHandler).Methods(http.MethodPost)
	r.HandleFunc("/tools/pwgen", pwgenHandler).Methods(http.MethodPost)
	r.HandleFunc("/tools/uuidgen", uuidHandler).Methods(http.MethodPost)

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

func mkpasswdHandler(w http.ResponseWriter, r *http.Request) {
	request := mkpasswdPayload{}
	decodeRequestPayload(w, r, &request)
	cmdargs := []string{
		request.Input,
		fmt.Sprintf("--method=%s", request.Method),
		fmt.Sprintf("--rounds=%d", request.Rounds),
	}

	if len(request.Salt) > 0 {
		cmdargs = append(cmdargs, fmt.Sprintf("--salt=%s", request.Salt))
	}

	// Run the command
	tool := exec.Command("/usr/bin/mkpasswd", cmdargs...)
	stdout, err := tool.Output()

	responseJSON(w, 200, Tool{
		Success: err == nil,
		Command: string(strings.Join(tool.Args, " ")),
		Output:  fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	})
}

func pwgenHandler(w http.ResponseWriter, r *http.Request) {
	request := pwgenPayload{}
	decodeRequestPayload(w, r, &request)
	cmdargs := []string{
		"-1",
		fmt.Sprintf("--num-passwords=%d", request.NumPasswords),
	}

	if request.RemoveChars != "" {
		cmdargs = append(cmdargs, fmt.Sprintf("--remove-chars=\"%s\"", request.RemoveChars))
	}

	if request.NoNumerals {
		cmdargs = append(cmdargs, "--no-numerals")
	}

	if request.NoCapitalize {
		cmdargs = append(cmdargs, "--no-capitalize")
	}

	if request.Ambiguous {
		cmdargs = append(cmdargs, "--ambiguous")
	}

	if request.Capitalize {
		cmdargs = append(cmdargs, "--capitalize")
	}

	if request.Numerals {
		cmdargs = append(cmdargs, "--numerals")
	}

	if request.Secure {
		cmdargs = append(cmdargs, "--secure")
	}

	if request.NoVowels {
		cmdargs = append(cmdargs, "--no-vowels")
	}

	if request.Symbols {
		cmdargs = append(cmdargs, "--symbols")
	}

	// Password length
	cmdargs = append(cmdargs, fmt.Sprintf("%d", request.Length))

	// Run the command
	tool := exec.Command("/usr/bin/pwgen", cmdargs...)
	stdout, err := tool.Output()

	responseJSON(w, 200, Tool{
		Success: err == nil,
		Command: string(strings.Join(tool.Args, " ")),
		Output:  fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	})
}

func uuidHandler(w http.ResponseWriter, r *http.Request) {
	request := uuidgenPayload{}
	decodeRequestPayload(w, r, &request)
	cmdargs := []string{}

	if request.Random {
		cmdargs = append(cmdargs, "--random")
	}

	if request.Time {
		cmdargs = append(cmdargs, "--time")
	}

	// Run the command
	tool := exec.Command("/usr/bin/uuidgen", cmdargs...)
	stdout, err := tool.Output()

	responseJSON(w, 200, Tool{
		Success: err == nil,
		Command: string(strings.Join(tool.Args, " ")),
		Output:  fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	})
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
