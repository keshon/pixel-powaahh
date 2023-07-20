package imgopt

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/ultimate-guitar/go-imagequant"
)

// ImageOptimizer is an interface for processing images.
type ImageOptimizer interface {
	CompressImage(imgData []byte, quality int) ([]byte, error)
}

// JPEGOptimizer implements the ImageOptimizer interface for JPEG images.
type JPEGOptimize struct{}

// NewJPEGOptimize creates a new instance of JPEGOptimize that implements the ImageOptimizer interface.
func NewJPEGOptimize() *JPEGOptimize {
	return &JPEGOptimize{}
}

// CompressImage compresses and encodes the JPEG image with the specified quality.
func (jp *JPEGOptimize) CompressImage(imgData []byte, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// PNGOptimizer implements the ImageOptimizer interface for PNG images.
type PNGOptimize struct{}

// NewPNGOptimize creates a new instance of PNGOptimize that implements the ImageOptimizer interface.
func NewPNGOptimize() *PNGOptimize {
	return &PNGOptimize{}
}

// CompressImage compresses and encodes the PNG image with the specified quality.
func (pp *PNGOptimize) CompressImage(imgData []byte, _ int) ([]byte, error) {
	buf := new(bytes.Buffer)
	img, err := png.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	attr, err := imagequant.NewAttributes()
	if err != nil {
		log.Printf("failed to create image attributes: %v", err)
		return nil, err
	}
	defer attr.Release()

	speed := 3

	err = attr.SetSpeed(speed)
	if err != nil {
		log.Printf("failed to set speed: %v", err)
		return nil, err
	}

	rgba32data := imageToRGBA32(img)

	imq, err := imagequant.NewImage(attr, string(rgba32data), width, height, 0)
	if err != nil {
		log.Printf("failed to create image quantization: %v", err)
		return nil, err
	}
	defer imq.Release()

	res, err := imq.Quantize(attr)
	if err != nil {
		log.Printf("failed to perform image quantization: %v", err)
		return nil, err
	}
	defer res.Release()

	rgb8data, err := res.WriteRemappedImage()
	if err != nil {
		log.Printf("failed to write remapped image: %v", err)
		return nil, err
	}

	prepImage := RGB8ToImage(res.GetImageWidth(), res.GetImageHeight(), rgb8data, res.GetPalette())

	err = png.Encode(buf, prepImage)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
