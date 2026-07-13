package paths

import (
	"os"
	"path/filepath"
)

// DetectSavePath tries to find the default game save directory.
func DetectSavePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	possiblePaths := []string{
		// Epic
		filepath.Join(home, "Documents", "Avalanche Studios", "Epic Games Store"),
		// Epic (OneDrive)
		filepath.Join(home, "OneDrive", "Documents", "Avalanche Studios", "Epic Games Store"),
		// Steam
		filepath.Join(home, "Documents", "Avalanche Studios"),
		// Steam (OneDrive)
		filepath.Join(home, "OneDrive", "Documents", "Avalanche Studios"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}
