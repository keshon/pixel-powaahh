package imageprocessing

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"pixel-powaahh/src/filesystem"

	"github.com/ultimate-guitar/go-imagequant"
)

// optimizePNG optimizes a PNG image by quantizing its colors and replacing it with the optimized version.
func optimizePNG(sourcefile, destfile string, speed int, outputDirectory string) {
	fh, err := os.Open(sourcefile)
	if err != nil {
		fmt.Println("failed to open source file: %w", err)
	}
	defer fh.Close()

	img, err := png.Decode(fh)
	if err != nil {
		fmt.Println("failed to decode PNG: %w", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	attr, err := imagequant.NewAttributes()
	if err != nil {
		fmt.Println("failed to create image attributes: %w", err)
	}
	defer attr.Release()

	err = attr.SetSpeed(speed)
	if err != nil {
		fmt.Println("failed to set speed: %w", err)
	}

	rgba32data := goImageToRGBA32(img)

	iqm, err := imagequant.NewImage(attr, string(rgba32data), width, height, 0)
	if err != nil {
		fmt.Println("failed to create image quantization: %w", err)
	}
	defer iqm.Release()

	res, err := iqm.Quantize(attr)
	if err != nil {
		fmt.Println("failed to perform image quantization: %w", err)
	}
	defer res.Release()

	rgb8data, err := res.WriteRemappedImage()
	if err != nil {
		fmt.Println("failed to write remapped image: %w", err)
	}

	im2 := rgb8PaletteToGoImage(res.GetImageWidth(), res.GetImageHeight(), rgb8data, res.GetPalette())

	// Create the output file path for the optimized image
	destFile := filepath.Join(outputDirectory, filesystem.GetRelativeFilePath(destfile))

	// Create the directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(destFile), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
	}

	fh2, err := os.Create(destFile)
	if err != nil {
		fmt.Println("failed to create destination file: %w", err)
	}
	defer fh2.Close()

	err = png.Encode(fh2, im2)
	if err != nil {
		fmt.Println("failed to encode PNG: %w", err)
	}

	fmt.Println("Optimization completed:", destFile)
}
