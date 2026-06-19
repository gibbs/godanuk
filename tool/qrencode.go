package tool

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
)

type qrencodeToolResponse struct {
	Success bool   `json:"success"`
	Command string `json:"command"`
	Output  string `json:"output"`
}

type qrencodeToolPayload struct {
	Input         string `json:"input"`
	Background    string `json:"background"`
	CaseSensitive bool   `json:"casesensitive"`
	Dpi           uint16 `json:"dpi"`
	Foreground    string `json:"foreground"`
	IgnoreCase    bool   `json:"ignorecase"`
	Level         string `json:"level"`
	Margin        uint16 `json:"margin"`
	Size          uint16 `json:"size"`
	Type          string `json:"type"`
}

func QrencodePayloadInterface() *qrencodeToolPayload {
	return &qrencodeToolPayload{}
}

func QrencodeHandler(request *qrencodeToolPayload) *qrencodeToolResponse {
	cmdargs := []string{
		fmt.Sprintf("--type=%s", request.Type),
	}

	if strings.EqualFold(request.Type, "svg") {
		cmdargs = append(cmdargs, fmt.Sprintf("--inline"))
	}

	if request.Margin > 0 {
		cmdargs = append(cmdargs, fmt.Sprintf("--margin=%d", request.Margin))
	}

	if request.Size > 0 {
		cmdargs = append(cmdargs, fmt.Sprintf("--size=%d", request.Size))
	}

	if request.Dpi > 0 {
		cmdargs = append(cmdargs, fmt.Sprintf("--dpi=%d", request.Dpi))
	}

	if request.Level != "" {
		cmdargs = append(cmdargs, fmt.Sprintf("--level=%s", request.Level))
	}

	if request.CaseSensitive {
		cmdargs = append(cmdargs, "--casesensitive")
	}

	if request.IgnoreCase {
		cmdargs = append(cmdargs, "--ignorecase")
	}

	if request.Foreground != "" {
		cmdargs = append(cmdargs, fmt.Sprintf("--foreground=%s", request.Foreground))
	}

	if request.Background != "" {
		cmdargs = append(cmdargs, fmt.Sprintf("--background=%s", request.Background))
	}

	cmdargs = append(cmdargs, "-o", "-", "--", request.Input)

	// Run the command
	tool := exec.Command("/usr/bin/qrencode", cmdargs...)
	stdout, err := tool.Output()
	trimmed := strings.TrimSuffix(string(stdout), "\n")

	return &qrencodeToolResponse{
		Success: err == nil,
		Command: strings.Join(tool.Args, " "),
		Output:  base64.StdEncoding.EncodeToString([]byte(trimmed)),
	}
}
