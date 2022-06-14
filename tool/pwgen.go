package tool

import (
	"fmt"
	"os/exec"
	"strings"
)

type pwgenToolResponse struct {
	Success bool   `json:"success"`
	Command string `json:"command"`
	Output  string `json:"output"`
}

type pwgenToolPayload struct {
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

func PwgenPayloadInterface() *pwgenToolPayload {
	return &pwgenToolPayload{}
}

func PwgenHandler(request *pwgenToolPayload) *pwgenToolResponse {
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

	return &pwgenToolResponse{
		Success: err == nil,
		Command: string(strings.Join(tool.Args, " ")),
		Output:  fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	}
}
