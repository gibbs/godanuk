package tool

import (
	"fmt"
	"os/exec"
	"strings"
)

type uuidgenToolResponse struct {
	Success bool   `json:"success"`
	Command string `json:"command"`
	Output  string `json:"output"`
}

type uuidgenToolPayload struct {
	Random bool `json:"random"`
	Time   bool `json:"time"`
}

func UuidgenPayloadInterface() *uuidgenToolPayload {
	return &uuidgenToolPayload{}
}

func UuidgenHandler(request *uuidgenToolPayload) *uuidgenToolResponse {
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

	return &uuidgenToolResponse{
		Success: err == nil,
		Command: string(strings.Join(tool.Args, " ")),
		Output:  fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	}
}
