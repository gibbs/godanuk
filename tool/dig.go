package tool

import (
	"fmt"
	"os/exec"
	"strings"
)

type digToolResponse struct {
	Success bool   `json:"success"`
	Command string `json:"command"`
	Output  string `json:"output"`
}

type digToolPayload struct {
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

func DigPayloadInterface() *digToolPayload {
	return &digToolPayload{}
}

func DigHandler(request *digToolPayload) *digToolResponse {
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
	if request.Types.ANY {
		cmdargs = append(cmdargs, "ANY")
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

	return &digToolResponse{
		Success: err == nil,
		Command: string(strings.Join(tool.Args, " ")),
		Output:  fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	}
}
