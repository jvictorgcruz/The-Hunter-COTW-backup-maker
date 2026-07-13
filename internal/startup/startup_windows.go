//go:build windows

package startup

import (
	"os"

	"golang.org/x/sys/windows/registry"
)

func SetAutostart() error {

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.SetStringValue("COTWBackupMaker", execPath+" --autostart")
	if err != nil {
		return err
	}

	return nil
}
