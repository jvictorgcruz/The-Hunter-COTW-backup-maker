//go:build windows

package browser

import "os/exec"

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	cmd = "rundll32"
	args = []string{"url.dll,FileProtocolHandler", url}

	return exec.Command(cmd, args...).Start()
}
