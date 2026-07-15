package backup

import (
	"backup-maker/internal/config"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type LocalProvider struct{}

func (lp *LocalProvider) ID() string {
	return "local"
}

func (lp *LocalProvider) IsConfigured() bool {
	cfg, err := config.LoadConfig()
	if err != nil || cfg.DestinationDir == "" {
		return false
	}
	info, err := os.Stat(cfg.DestinationDir)
	return err == nil && info.IsDir()
}

func (lp *LocalProvider) Send(localZipPath string, filename string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	destPath := filepath.Join(cfg.DestinationDir, filename)

	srcFile, err := os.Open(localZipPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func (lp *LocalProvider) Cleanup(limit int) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	files, err := os.ReadDir(cfg.DestinationDir)
	if err != nil {
		return err
	}

	var backupFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), backupPrefix) && strings.HasSuffix(f.Name(), backupSuffix) {
			backupFiles = append(backupFiles, filepath.Join(cfg.DestinationDir, f.Name()))
		}
	}

	if len(backupFiles) <= limit {
		return nil
	}

	sort.Strings(backupFiles)

	filesToRemove := backupFiles[:len(backupFiles)-limit]
	for _, path := range filesToRemove {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	RegisterProvider(&LocalProvider{})
}
