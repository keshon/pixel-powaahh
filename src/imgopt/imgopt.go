package imgopt

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
)

// ImageCompressor is an interface for compressing images.
type ImageCompressor interface {
	Compress(imageData []byte, quality int) ([]byte, error)
}

// JPEGCompressor implements the ImageCompressor interface for JPEG images.
type JPEGCompressor struct{}

// NewJPEGCompressor creates a new instance of JPEGCompressor.
func NewJPEGCompressor() *JPEGCompressor {
	return &JPEGCompressor{}
}

// Compress compresses and encodes the image as JPEG with the specified quality.
func (jc *JPEGCompressor) Compress(imageData []byte, quality int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err // Error: Unable to decode image
	}

	var compressedImage bytes.Buffer
	err = jpeg.Encode(&compressedImage, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err // Error: Unable to encode image as JPEG
	}

	return compressedImage.Bytes(), nil
}

// PNGCompressor implements the ImageCompressor interface for PNG images.
type PNGCompressor struct{}

// NewPNGCompressor creates a new instance of PNGCompressor.
func NewPNGCompressor() *PNGCompressor {
	return &PNGCompressor{}
}

// Compress compresses and encodes the image as PNG.
func (pc *PNGCompressor) Compress(imageData []byte, _ int) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err // Error: Unable to decode image
	}

	var compressedImage bytes.Buffer
	err = png.Encode(&compressedImage, img)
	if err != nil {
		return nil, err // Error: Unable to encode image as PNG
	}

	return compressedImage.Bytes(), nil
}
