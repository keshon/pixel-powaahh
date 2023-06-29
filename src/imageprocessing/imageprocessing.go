package imageprocessing

import (
	"fmt"
	"path/filepath"
	"pixel-powaahh/src/config"
	"pixel-powaahh/src/filesystem"
	"pixel-powaahh/src/utilities"
	"strings"
)

// optimizeImage optimizes an image based on its file format.
func OptimizeImage(file string, jpgOnly bool, pngOnly bool, quality int, outputDirectory string) {
	// Check the file format
	fileExt := strings.ToLower(filepath.Ext(file))
	imageFormat := utilities.GetImageFormat(fileExt)

	if imageFormat == config.JPEG && pngOnly {
		return
	}

	if imageFormat == config.PNG && jpgOnly {
		return
	}

	fmt.Println("Optimizing:", file)

	switch imageFormat {
	case config.JPEG:
		optimizeJPEG(file, filesystem.AddFileSuffix(file, ""), quality, outputDirectory)
	case config.PNG:
		optimizePNG(file, filesystem.AddFileSuffix(file, ""), 3, outputDirectory)
	default:
		fmt.Println("Unsupported image format:", fileExt)
	}
}
