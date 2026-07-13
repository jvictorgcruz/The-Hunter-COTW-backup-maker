package backup

import (
	"archive/zip"
	"backup-maker/internal/config"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	err = CreateBackup(config.SourceDir, config.DestinationDir, config.MaxBackups)
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

func CreateBackup(sourceDir string, destinationPath string, limit int) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	destinationZipPath := filepath.Join(destinationPath, backupPrefix+timestamp+backupSuffix)
	zipFile, err := os.Create(destinationZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
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

	if err != nil {
		return err
	}

	if err := cleanupOldBackups(destinationPath, limit); err != nil {
		return err
	}

	return nil
}

func cleanupOldBackups(destinationPath string, limit int) error {
	files, err := os.ReadDir(destinationPath)
	if err != nil {
		return err
	}
	var backupFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), backupPrefix) && strings.HasSuffix(f.Name(), backupSuffix) {
			backupFiles = append(backupFiles, filepath.Join(destinationPath, f.Name()))
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
