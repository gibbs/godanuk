package tool

import (
	"fmt"
	"os/exec"
	"strings"
)

type rdapToolResponse struct {
	Success bool   `json:"success"`
	Command string `json:"command"`
	Output  string `json:"output"`
}

type rdapToolPayload struct {
	Query string `json:"query"`

	Verbose      bool   `json:"verbose"`
	Timeout      uint16 `json:"timeout"`
	Insecure     bool   `json:"insecure"`
	Experimental bool   `json:"experimental"`

	Format string `json:"format"`

	Server string   `json:"server"`
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}

func RdapPayloadInterface() *rdapToolPayload {
	return &rdapToolPayload{}
}

func RdapHandler(request *rdapToolPayload) *rdapToolResponse {
	cmdargs := []string{}

	if request.Verbose {
		cmdargs = append(cmdargs, "--verbose")
	}

	if request.Timeout > 0 {
		cmdargs = append(cmdargs, fmt.Sprintf("--timeout=%d", request.Timeout))
	}

	if request.Insecure {
		cmdargs = append(cmdargs, "--insecure")
	}

	if request.Experimental {
		cmdargs = append(cmdargs, "--experimental")
	}

	switch strings.ToLower(request.Format) {
	case "whois":
		cmdargs = append(cmdargs, "--whois")
	case "json":
		cmdargs = append(cmdargs, "--json")
	case "raw":
		cmdargs = append(cmdargs, "--raw")
	default:
		cmdargs = append(cmdargs, "--text")
	}

	if request.Server != "" {
		cmdargs = append(cmdargs, fmt.Sprintf("--server=%s", request.Server))
	}

	if request.Type != "" {
		cmdargs = append(cmdargs, fmt.Sprintf("--type=%s", request.Type))
	}

	for _, field := range request.Fields {
		if field != "" {
			cmdargs = append(cmdargs, "-f", field)
		}
	}

	cmdargs = append(cmdargs, "--", request.Query)

	// Run the command
	tool := exec.Command("/usr/bin/rdap", cmdargs...)
	stdout, err := tool.Output()

	return &rdapToolResponse{
		Success: err == nil,
		Command: strings.Join(tool.Args, " "),
		Output:  fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	}
}
