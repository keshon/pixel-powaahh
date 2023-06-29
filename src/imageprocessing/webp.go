package imageprocessing

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"pixel-powaahh/src/filesystem"

	"github.com/chai2010/webp"
)

// ConvertToWebP converts an image to the WebP format.
func ConvertToWebP(file string, quality int, outputDirectory string) {
	// Open the image file
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	// Decode the image file
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Create the output file path for the converted image
	destFile := filepath.Join(outputDirectory, filesystem.GetRelativeFilePath(file))
	destFile = filesystem.ChangeFileExtension(destFile, ".webp")
	fmt.Println(destFile)

	// Create the directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(destFile), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	outputFile, err := os.Create(destFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Encode the image into the WebP format
	if err = webp.Encode(outputFile, img, &webp.Options{Quality: float32(quality)}); err != nil {
		fmt.Println("Error encoding image to WebP:", err)
		return
	}

	fmt.Println("Conversion to WebP completed:", destFile)
}
