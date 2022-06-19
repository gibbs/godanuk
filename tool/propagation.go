package tool

import (
	"fmt"
	"os/exec"
	"strings"
)

type propagationToolResponse struct {
	Success  bool   `json:"success"`
	Provider string `json:"provider"`
	Output   string `json:"output"`
}

type propagationToolPayload struct {
	Name       string `json:"name"`
	Nameserver string `json:"nameserver"`
	Types      struct {
		A     bool `json:"a"`
		AAAA  bool `json:"aaaa"`
		CAA   bool `json:"caa"`
		CNAME bool `json:"cname"`
		MX    bool `json:"mx"`
		NS    bool `json:"ns"`
		PTR   bool `json:"ptr"`
		SOA   bool `json:"soa"`
		SRV   bool `json:"srv"`
		TXT   bool `json:"txt"`
	} `json:"types"`
}

func PropagationPayloadInterface() *propagationToolPayload {
	return &propagationToolPayload{}
}

func PropagationHandler(request *propagationToolPayload) *propagationToolResponse {
	nameservers := map[string]map[string]string{
		// Global
		"cloudflare": {
			"name":   "Cloudflare",
			"server": "1.1.1.1",
		},
		"google": {
			"name":   "Google",
			"server": "8.8.8.8",
		},
		"opendns": {
			"name":   "OpenDNS",
			"server": "208.67.222.222",
		},
		// UK
		"bt": {
			"name":   "BT",
			"server": "62.6.40.178",
		},
		"plusnet": {
			"name":   "PlusNet",
			"server": "212.159.13.49",
		},
		"sky": {
			"name":   "Sky Broadband",
			"server": "90.207.238.97",
		},
		"ee": {
			"name":   "EE",
			"server": "87.237.17.198",
		},
		"virgin": {
			"name":   "Virgin Media",
			"server": "194.168.4.100",
		},
		"talktalk": {
			"name":   "TalkTalk",
			"server": "62.24.134.1",
		},
		"vodafone": {
			"name":   "Vodafone",
			"server": "90.255.255.90",
		},
		"zen": {
			"name":   "Zen Internet",
			"server": "212.23.3.100",
		},
		"hypernotic": {
			"name":   "Hyperoptic",
			"server": "141.0.144.64",
		},
		"kcom": {
			"name":   "KCOM",
			"server": "212.50.160.38",
		},
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
		fmt.Sprintf("@%s", nameservers[request.Nameserver]["server"]),
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
	if request.Types.TXT {
		cmdargs = append(cmdargs, "TXT")
		cmdargs = append(cmdargs, request.Name)
	}

	// Run the command
	tool := exec.Command("/usr/bin/dig", append(cmdopts, cmdargs...)...)
	stdout, err := tool.Output()

	return &propagationToolResponse{
		Success:  err == nil,
		Provider: nameservers[request.Nameserver]["name"],
		Output:   fmt.Sprint(strings.TrimSuffix(string(stdout), "\n")),
	}
}
