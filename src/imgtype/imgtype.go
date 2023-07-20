package imgtype

import (
	"errors"
	"log"
	"path/filepath"
	"strings"
)

// ImageFormat represents the supported image formats.
type ImageFormat int

const (
	Unknown             = -1
	JPEG    ImageFormat = iota
	PNG
	WebP
)

var supportedFormats = map[string]ImageFormat{
	".jpg":  JPEG,
	".jpeg": JPEG,
	".png":  PNG,
	".webp": WebP,
}

// GetImageFormat returns the image format based on the file extension.
func GetImageFormat(fileExt string) (ImageFormat, error) {
	ext := strings.ToLower(filepath.Ext(fileExt))
	format, ok := supportedFormats[ext]
	if !ok {
		errMsg := "unsupported image format: " + ext
		log.Printf("Error in GetImageFormat: %s", errMsg)
		return Unknown, errors.New(errMsg)
	}
	return format, nil
}

// GetImageFormatName returns the name of the image format.
func GetImageFormatName(format ImageFormat) string {
	for ext, f := range supportedFormats {
		if f == format {
			return ext
		}
	}
	return ""
}

// GetImageExtensions returns the supported image file extensions.
func GetImageExtensions() []string {
	extensions := make([]string, 0, len(supportedFormats))
	for ext := range supportedFormats {
		extensions = append(extensions, ext)
	}
	return extensions
}

// IsImage checks if a file has a supported image extension.
func IsImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	_, ok := supportedFormats[ext]
	return ok
}
