//go:build linux

package browser

import (
	"os/exec"
)

func OpenBrowser(url string) error {
	var cmd string
	var args []string
	cmd = "xdg-open"
	args = []string{url}

	return exec.Command(cmd, args...).Start()
}
