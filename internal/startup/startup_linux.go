//go:build linux

package startup

import (
	"os"
	"path/filepath"
	"strings"
)

func SetAutostart() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	startupFile := filepath.Join(home, ".config", "autostart", "backup-maker.desktop")
	dir := filepath.Dir(startupFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(startupFile)
	if err != nil {
		return err
	}
	defer file.Close()

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	var builder strings.Builder
	builder.WriteString("[Desktop Entry]\n")
	builder.WriteString("Type=Application\n")
	builder.WriteString("Name=COTW Backup Maker\n")
	builder.WriteString("Exec=")
	builder.WriteString(execPath)
	builder.WriteString(" --autostart\n")
	builder.WriteString("Hidden=false\n")
	builder.WriteString("NoDisplay=false\n")
	builder.WriteString("X-GNOME-Autostart-enabled=true")
	content := builder.String()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
