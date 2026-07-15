package backup

import (
	"archive/zip"
	"backup-maker/internal/config"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	backupPrefix = "backup-maker_"
	backupSuffix = ".zip"
)

func AutoBackup() error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	err = CreateBackup(config.SourceDir, config.EnabledProviders, config.MaxBackups)
	if err != nil {
		return err
	}

	return nil
}

func BackupOnStartup() error {
	if !config.BackupOnStartupActive() {
		return nil
	}

	if err := AutoBackup(); err != nil {
		return err
	}

	return nil
}

func CreateBackup(sourceDir string, enabledProviderIDs []string, limit int) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := backupPrefix + timestamp + backupSuffix
	tempZipPath := filepath.Join(os.TempDir(), filename)
	defer os.Remove(tempZipPath)
	err := func() error {
		zipFile, err := os.Create(tempZipPath)
		if err != nil {
			return err
		}
		defer zipFile.Close()

		archive := zip.NewWriter(zipFile)
		defer archive.Close()

		return filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			relativePath, err := filepath.Rel(sourceDir, path)
			if err != nil {
				return err
			}

			archiveWriter, err := archive.Create(relativePath)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(archiveWriter, file)
			if err != nil {
				return err
			}

			return nil
		})
	}()

	if err != nil {
		return err
	}

	var errorsList []error
	for _, providerid := range enabledProviderIDs {
		provider := GetProvider(providerid)
		if provider == nil {
			continue
		}
		if !provider.IsConfigured() {
			continue
		}

		if err := provider.Send(tempZipPath, filename); err != nil {
			errorsList = append(errorsList, err)
			continue
		}
		if err := provider.Cleanup(limit); err != nil {
			errorsList = append(errorsList, err)
		}
	}

	if len(errorsList) > 0 {
		return errors.Join(errorsList...)
	}

	return nil
}
