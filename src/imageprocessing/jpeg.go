package imageprocessing

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"pixel-powaahh/src/filesystem"
)

// optimizeJPEG optimizes a JPEG image by compressing and encoding it with the specified quality.
func optimizeJPEG(sourcefile, destfile string, quality int, outputDirectory string) {
	// Open the image file
	f, err := os.Open(sourcefile)
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

	// Create the output file path for the optimized image
	destFile := filepath.Join(outputDirectory, filesystem.GetRelativeFilePath(destfile))

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

	// Compress and encode the image into the output file
	if quality == 0 {
		quality = 80
	}
	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: quality})
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	fmt.Println("Optimization completed:", destFile)
}
