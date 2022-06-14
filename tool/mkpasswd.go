package tool

import (
	"fmt"
	"os/exec"
	"strings"
)

type mkpasswdToolResponse struct {
	Success bool   `json:"success"`
	Command string `json:"command"`
	Output  string `json:"output"`
}

type mkpasswdToolPayload struct {
	Input  string `json:"input"`
	Salt   string `json:"salt"`
	Rounds uint32 `json:"rounds"`
	Method string `json:"method"`
}

func MkpasswdPayloadInterface() *mkpasswdToolPayload {
	return &mkpasswdToolPayload{}
}

func MkpasswdHandler(request *mkpasswdToolPayload) *mkpasswdToolResponse {
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

	return &mkpasswdToolResponse{
		Success: err == nil,
		Command: string(strings.Join(tool.Args, " ")),
		Output:  fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	}
}
