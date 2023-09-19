package imagetype

import (
	"errors"
	"log"
	"path/filepath"
	"strings"
)

type ImageFormat int

const (
	Unknown ImageFormat = iota - 1
	JPEG
	PNG
	WebP
)

type ImageType interface {
	GetFormatByExtension(fileExt string) (ImageFormat, error)
	GetFormatName(format ImageFormat) string
	GetSupportedExtensions() []string
	IsSupportedExtension(filename string) bool
}

type ImageTypeImpl struct {
	supportedFormats map[string]ImageFormat
}

func New() ImageType {
	return &ImageTypeImpl{
		supportedFormats: map[string]ImageFormat{
			".jpg":  JPEG,
			".jpeg": JPEG,
			".png":  PNG,
			".webp": WebP,
		},
	}
}

func (it *ImageTypeImpl) GetFormatByExtension(fileExt string) (ImageFormat, error) {
	ext := strings.ToLower(filepath.Ext(fileExt))
	format, ok := it.supportedFormats[ext]
	if !ok {
		errMsg := "unsupported image format for extension " + ext
		log.Printf("GetFormatByExtension error: %s", errMsg)
		return Unknown, errors.New(errMsg)
	}
	return format, nil
}

func (it *ImageTypeImpl) GetFormatName(format ImageFormat) string {
	for ext, f := range it.supportedFormats {
		if f == format {
			return ext
		}
	}
	return ""
}

func (it *ImageTypeImpl) GetSupportedExtensions() []string {
	extensions := make([]string, 0, len(it.supportedFormats))
	for ext := range it.supportedFormats {
		extensions = append(extensions, ext)
	}
	return extensions
}

func (it *ImageTypeImpl) IsSupportedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	_, ok := it.supportedFormats[ext]
	return ok
}
