package config

import (
	"os"
	"path/filepath"
)

// Config contains the application configuration.
type Config struct {
	UploadDir    string
	ProcessedDir string
	BinDir       string
	RunGUIMode   bool
}

// DefaultConfig creates a new instance of Config with default values.
func NewConfig() *Config {
	uploadDir := "uploads"
	processedDir := "processed"

	// Get the absolute path to the executable
	binDir, err := executableDir()
	if err != nil {
		// Handle the error (e.g., fallback to a default path)
		binDir = "."
	}

	return &Config{
		UploadDir:    uploadDir,
		ProcessedDir: processedDir,
		BinDir:       binDir,
		RunGUIMode:   false,
	}
}

// executableDir returns the absolute path to the directory containing the executable.
func executableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(exePath), nil
}
