package utilities

import (
	"path/filepath"
	"pixel-powaahh/src/config"
	"strings"
)

// IsImageFile checks if a file has an image extension.
func IsImageFile(filename string) bool {
	extensions := []string{".jpg", ".jpeg", ".png"}
	ext := strings.ToLower(filepath.Ext(filename))

	for _, extension := range extensions {
		if extension == ext {
			return true
		}
	}

	return false
}

// GetImageFormat returns the image format based on the file extension.
func GetImageFormat(fileExt string) config.ImageFormat {
	switch fileExt {
	case ".jpg", ".jpeg":
		return config.JPEG
	case ".png":
		return config.PNG
	case ".webp": // New case for WebP format
		return config.WebP
	default:
		return -1
	}
}
