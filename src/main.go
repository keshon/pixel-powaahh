package main

import (
	"fmt"
	"os"
	"pixel-powaahh/src/config"
	"pixel-powaahh/src/filesystem"
	"pixel-powaahh/src/imageprocessing"
	"sync"

	"github.com/spf13/cobra"
)

func main() {
	var jpgOnly bool
	var pngOnly bool
	var toWebp bool
	var quality int

	rootCmd := &cobra.Command{
		Use:   "pp",
		Short: "Pixel Powaahh - JPG and PNG optimizer and converter",
		Run: func(cmd *cobra.Command, args []string) {
			runPP(jpgOnly, pngOnly, toWebp, quality)
		},
	}

	rootCmd.Flags().BoolVarP(&jpgOnly, "jpg", "j", false, "Optimize JPEG files only")
	rootCmd.Flags().BoolVarP(&pngOnly, "png", "p", false, "Optimize PNG files only")
	rootCmd.Flags().BoolVarP(&toWebp, "webp", "w", false, "Convert images to WebP format")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 80, "Compression ratio for JPEG or WebP")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runPP(jpgOnly, pngOnly, toWebp bool, quality int) {
	// Define the input directory path to scan
	inputDirectory := "./" + config.UPLOAD_DIR

	// Create the input directory if it doesn't exist
	if err := os.MkdirAll(inputDirectory, os.ModePerm); err != nil {
		fmt.Println("Error creating input directory:", err)
		return
	}

	// Define the output directory path to save processed images
	outputDirectory := "./" + config.PROCCESED_DIR

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// Empty the processed folder
	if err := filesystem.EmptyProcessedFolder(outputDirectory); err != nil {
		fmt.Println("Error emptying processed folder:", err)
		return
	}

	// Fetch all image files from the directory
	imageFiles, err := filesystem.FetchImageFiles(inputDirectory)
	if err != nil {
		fmt.Println("Error fetching image files:", err)
		return
	}

	// Define the maximum number of concurrent goroutines
	maxConcurrency := 12

	// Create a wait group to wait for goroutines to finish
	var wg sync.WaitGroup

	// Create a semaphore to limit concurrent goroutines
	semaphore := make(chan struct{}, maxConcurrency)

	// Iterate over the image files and process them concurrently
	for _, file := range imageFiles {
		semaphore <- struct{}{} // Acquire a semaphore slot
		wg.Add(1)
		go func(file string) {
			defer func() {
				<-semaphore // Release the semaphore slot
				wg.Done()
			}()
			if toWebp {
				imageprocessing.ConvertToWebP(file, quality, outputDirectory)
			} else {
				imageprocessing.OptimizeImage(file, jpgOnly, pngOnly, quality, outputDirectory)
			}
		}(file)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("Done")
}
