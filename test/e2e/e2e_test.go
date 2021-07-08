// +build e2e

package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFluxCheckPre(t *testing.T) {
	var lines []string
	var err error

	lines, err = run("kubectl", "version", "--short")
	assert.NoError(t, err)

	clientVersion := strings.SplitN(lines[0], ": v", 2)[1]
	serverVersion := strings.SplitN(lines[1], ": v", 2)[1]

	lines, err = run(fluxBinary, "check", "--pre")
	assert.NoError(t, err)

	assert.Equal(t, "► checking prerequisites", lines[0])
	assert.Equal(t, fmt.Sprintf("✔ kubectl %s >=1.18.0-0", clientVersion), lines[1])
	assert.Equal(t, fmt.Sprintf("✔ Kubernetes %s >=1.16.0-0", serverVersion), lines[2])
	assert.Equal(t, "✔ prerequisites checks passed", lines[3])
}

func TestFluxInstallManifests(t *testing.T) {
	var lines []string
	var err error

	lines, err = run(fluxBinary, "install", "--manifests", "../../manifests/install/")
	assert.NoError(t, err)

	assert.Equal(t, "✚ generating manifests", lines[0])
	assert.Equal(t, "✔ manifests build completed", lines[1])
	assert.Equal(t, "► installing components in flux-system namespace", lines[2])
	assert.Equal(t, "◎ verifying installation", lines[3])
	// line 4 to 7 are deployment ready for each controller
	for i := 4; i <= 7; i++ {
		assert.Contains(t, lines[i], "deployment ready")
	}
	assert.Equal(t, "✔ install finished", lines[8])
}

func TestFluxCreateSecret(t *testing.T) {
	var lines []string
	var err error

	lines, err = run(fluxBinary,
		"create", "secret", "git", "git-ssh-test",
		"--url", "ssh://git@github.com/stefanprodan/podinfo")
	assert.NoError(t, err)
	assert.Contains(t, lines[0], "✚ deploy key:")
	assert.Equal(t, "► secret 'git-ssh-test' created in 'flux-system' namespace", lines[2])

	lines, err = run(fluxBinary,
		"create", "secret", "git", "git-https-test",
		"--url", "https://github.com/stefanprodan/podinfo",
		"--username=test", "--password=test")
	assert.NoError(t, err)
	assert.Equal(t, "► secret 'git-https-test' created in 'flux-system' namespace", lines[0])

	lines, err = run(fluxBinary,
		"create", "secret", "helm", "helm-test",
		"--username=test", "--password=test")
	assert.NoError(t, err)
}
