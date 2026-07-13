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
		// Steam
		filepath.Join(home, "Documents", "Avalanche Studios"),
		// Epic
		filepath.Join(home, "Documents", "Avalanche Studios", "Epic Games Store"),
		// Steam (OneDrive)
		filepath.Join(home, "OneDrive", "Documents", "Avalanche Studios"),
		// Epic (OneDrive)
		filepath.Join(home, "OneDrive", "Documents", "Avalanche Studios", "Epic Games Store"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}
