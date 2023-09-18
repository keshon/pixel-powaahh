package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pp/src/conf"
	"pp/src/fileio"
	"pp/src/imgconv"
	"pp/src/imgopt"
	"pp/src/imgtype"
	"sync"

	"github.com/spf13/cobra"
)

func main() {
	cobra.MousetrapHelpText = ""

	var jpgOnly bool
	var pngOnly bool
	var toWebp bool
	var quality int

	rootCmd := &cobra.Command{
		Use:   "pp",
		Short: "Pixel Powaahh - JPG and PNG optimizer and converter",
		Run: func(cmd *cobra.Command, args []string) {
			pp(jpgOnly, pngOnly, toWebp, quality)
		},
	}

	rootCmd.Flags().BoolVarP(&jpgOnly, "jpg", "j", false, "Optimize JPEG files only")
	rootCmd.Flags().BoolVarP(&pngOnly, "png", "p", false, "Optimize PNG files only")
	rootCmd.Flags().BoolVarP(&toWebp, "webp", "w", false, "Convert images to WebP format")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 80, "Compression ratio for JPEG or WebP: 1-100")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func pp(jpgOnly, pngOnly, toWebp bool, quality int) {
	// Create objects
	fio := fileio.NewFileOp()
	conv := imgconv.NewImgConvert()
	jpeg := imgopt.NewJPEGCompressor()
	png := imgopt.NewPNGCompressor()

	// Define the input directory path to scan
	inputDirectory := "./" + conf.UPLOAD_DIR

	// Create the input directory if it doesn't exist
	if err := os.MkdirAll(inputDirectory, os.ModePerm); err != nil {
		log.Fatalf("Error creating input directory: %v", err)
	}

	// Define the output directory path to save processed images
	outputDirectory := "./" + conf.PROCCESED_DIR

	// Create the processed directory if it doesn't exist
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// Empty the processed directory if it exists
	if err := fio.EmptyDir(outputDirectory); err != nil {
		log.Fatalf("Error emptying processed folder: %v", err)
	}

	// Fetch all image files from the upload directory
	files, err := fio.FetchFiles(inputDirectory)
	if err != nil {
		log.Fatalf("Error fetching image files: %v", err)
	}

	// Define the maximum number of concurrent goroutines
	maxConcurrency := 12

	// Create a wait group to wait for goroutines to finish
	var wg sync.WaitGroup

	// Create a semaphore to limit concurrent goroutines
	semaphore := make(chan struct{}, maxConcurrency)

	// Iterate over the image files and process them concurrently
	for _, file := range files {
		semaphore <- struct{}{} // Acquire a semaphore slot
		wg.Add(1)
		go func(file string) {
			defer func() {
				<-semaphore // Release the semaphore slot
				wg.Done()
			}()

			// Read uploaded file content
			uploadedData, err := fio.ReadFile(inputDirectory + "/" + file)
			if err != nil {
				log.Printf("Error reading image data: %v", err)
				return
			}

			// Check the file format
			fileFormat, err := imgtype.GetImageFormat(file)
			if err != nil {
				log.Printf("Error detecting image format: %v", err)
				return
			}

			var processedData []byte
			var destFile string

			destFile = filepath.Join(outputDirectory, file)

			if toWebp {
				destFile = fio.ChangeExt(destFile, ".webp")

				// Convert to WebP
				processedData, err = conv.ConvertImg(uploadedData, fileFormat, imgtype.WebP, quality)
				if err != nil {
					log.Printf("Error converting image: %v", err)
					return
				}

			} else {

				if fileFormat == imgtype.JPEG && pngOnly {
					return
				}

				if fileFormat == imgtype.PNG && jpgOnly {
					return
				}

				switch fileFormat {
				case imgtype.JPEG:
					// Compress JPEG
					processedData, err = jpeg.Compress(uploadedData, quality)
					if err != nil {
						log.Printf("Error compressing image: %v", err)
						return
					}
				case imgtype.PNG:
					// Compress PNG
					processedData, err = png.Compress(uploadedData, quality)
					if err != nil {
						log.Printf("Error compressing image: %v", err)
						return
					}
				default:
					log.Printf("Unsupported image format: %v", fileFormat)
				}
			}

			// Save processed content to file
			err = fio.SaveFile(destFile, processedData)
			if err != nil {
				log.Printf("Error saving converted image to file: %v", err)
			}

			fmt.Println("Optimization completed:", destFile)
		}(file)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("Done")
}
