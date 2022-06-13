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

type digPayload struct {
	Name       string `json:"name"`
	Nameserver string `json:"nameserver"`
	Types      struct {
		A      bool `json:"a"`
		AAAA   bool `json:"aaaa"`
		ANY    bool `json:"any"`
		CAA    bool `json:"caa"`
		CNAME  bool `json:"cname"`
		DNSKEY bool `json:"dnskey"`
		DS     bool `json:"ds"`
		MX     bool `json:"mx"`
		NS     bool `json:"ns"`
		PTR    bool `json:"ptr"`
		SOA    bool `json:"soa"`
		SRV    bool `json:"srv"`
		TLSA   bool `json:"tlsa"`
		TSIG   bool `json:"tsig"`
		TXT    bool `json:"txt"`
	} `json:"types"`
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
	r.HandleFunc("/tools/dig", digHandler).Methods(http.MethodPost)
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

func digHandler(w http.ResponseWriter, r *http.Request) {
	request := digPayload{}
	decodeRequestPayload(w, r, &request)

	nameservers := map[string]string{
		"cloudflare": "1.1.1.1",
		"google":     "8.8.8.8",
		"quad9":      "9.9.9.9",
		"opendns":    "208.67.222.222",
		"comodo":     "8.26.56.26",
	}

	cmdopts := []string{
		"+yaml",
		"+notcp",
		"+recurse",
		"+qr",
		"+time=5",
		"+tries=3",
		"+retry=2",
	}

	cmdargs := []string{
		fmt.Sprintf("@%s", nameservers[request.Nameserver]),
	}

	// @fixme: Improve - repetitious
	if request.Types.A {
		cmdargs = append(cmdargs, "A")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.AAAA {
		cmdargs = append(cmdargs, "AAAA")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.CAA {
		cmdargs = append(cmdargs, "CAA")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.CNAME {
		cmdargs = append(cmdargs, "CNAME")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.DNSKEY {
		cmdargs = append(cmdargs, "DNSKEY")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.DS {
		cmdargs = append(cmdargs, "DS")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.MX {
		cmdargs = append(cmdargs, "MX")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.NS {
		cmdargs = append(cmdargs, "NS")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.PTR {
		cmdargs = append(cmdargs, "PTR")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.SOA {
		cmdargs = append(cmdargs, "SOA")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.SRV {
		cmdargs = append(cmdargs, "SRV")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.TLSA {
		cmdargs = append(cmdargs, "TLSA")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.TSIG {
		cmdargs = append(cmdargs, "TSIG")
		cmdargs = append(cmdargs, request.Name)
	}
	if request.Types.TXT {
		cmdargs = append(cmdargs, "TXT")
		cmdargs = append(cmdargs, request.Name)
	}

	// Run the command
	tool := exec.Command("/usr/bin/dig", append(cmdopts, cmdargs...)...)
	stdout, err := tool.Output()

	responseJSON(w, 200, Tool{
		Success: err == nil,
		Command: string(strings.Join(tool.Args, " ")),
		Output:  fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	})
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
