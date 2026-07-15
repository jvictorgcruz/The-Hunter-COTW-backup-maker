//go:build windows

package browser

import "os/exec"

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	cmd = "cmd"
	args = []string{"/c", "start", url}

	return exec.Command(cmd, args...).Start()
}
