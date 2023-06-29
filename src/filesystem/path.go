package filesystem

import (
	"path/filepath"
	"pixel-powaahh/src/config"
	"strings"
)

// addFileSuffix adds a suffix to a file path.
func AddFileSuffix(path, suffix string) string {
	ext := filepath.Ext(path)
	filename := strings.TrimSuffix(path, ext)
	return filename + suffix + ext
}

// changeFileExtension changes the file extension of a path to a new extension.
func ChangeFileExtension(path, newExtension string) string {
	ext := filepath.Ext(path)
	filename := strings.TrimSuffix(path, ext)
	return filename + "." + strings.TrimPrefix(newExtension, ".")
}

// getRelativeFilePath returns the relative path of a file by removing the base directory.
func GetRelativeFilePath(filePath string) string {
	baseDir := config.UPLOAD_DIR // Base directory of the uploaded files
	relativePath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		// Handle error if the file path is not relative to the base directory
		return ""
	}
	return relativePath
}
