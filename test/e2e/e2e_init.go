// +build e2e

package e2e

import (
	"os"
	"os/exec"
	"strings"
)

var fluxBinary string

func init() {
	fluxBinary = os.Getenv("FLUX_E2E_BINARY")
	if fluxBinary == "" {
		fluxBinary = "../../bin/flux"
	}
}

func run(name string, args ...string) ([]string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	return lines, nil
}
