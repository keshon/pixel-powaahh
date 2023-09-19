package pixelita

import (
	"app/internal/config"
	"app/internal/filesystem"
	"app/internal/imageencode"
	"app/internal/imagetype"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func (px *Pixelita) StartCLI(jpgOnly, pngOnly, toWebp bool, quality int) {
	var logBuffer bytes.Buffer

	log.SetOutput(&logBuffer)

	// Create objects
	config := config.NewConfig()
	fs := filesystem.NewFileSystemImpl(config)
	jpeg := imageencode.NewJPEGEncoder()
	png := imageencode.NewPNGEncoder()

	// Define the input directory path to scan
	inputDirectory := "./" + config.UploadDir

	// Create the input directory if it doesn't exist
	if err := os.MkdirAll(inputDirectory, os.ModePerm); err != nil {
		log.Fatalf("Error creating input directory: %v", err)
	}

	// Define the output directory path to save processed images
	outputDirectory := "./" + config.ProcessedDir

	// Create the processed directory if it doesn't exist
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// Empty the processed directory if it exists
	if err := fs.ClearDirectory(outputDirectory); err != nil {
		log.Fatalf("Error emptying processed folder: %v", err)
	}

	// Fetch all image files from the upload directory
	files, err := fs.GetImageFiles(inputDirectory)
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
			uploadedData, err := fs.ReadFile(inputDirectory + "/" + file)
			if err != nil {
				log.Printf("Error reading image data: %v", err)
				return
			}

			// Check the file format
			fileFormat, err := imagetype.New().GetFormatByExtension(file)
			if err != nil {
				log.Printf("Error detecting image format: %v", err)
				return
			}

			var processedData []byte
			var destFile string

			destFile = filepath.Join(outputDirectory, file)

			if toWebp {
				destFile = fs.ChangeFileExtension(destFile, ".webp")

				// Convert to WebP
				processedData, err = px.convertBetweenFormats(uploadedData, fileFormat, imagetype.WebP, quality)
				if err != nil {
					log.Printf("Error converting image: %v", err)
					return
				}

			} else {

				if fileFormat == imagetype.JPEG && pngOnly {
					return
				}

				if fileFormat == imagetype.PNG && jpgOnly {
					return
				}

				switch fileFormat {
				case imagetype.JPEG:
					// Compress JPEG
					processedData, err = jpeg.Encode(uploadedData, quality)
					if err != nil {
						log.Printf("Error compressing image: %v", err)
						return
					}
				case imagetype.PNG:
					// Compress PNG
					processedData, err = png.Encode(uploadedData, 0, 10, 100, 3)
					if err != nil {
						log.Printf("Error compressing image: %v", err)
						return
					}
				default:
					log.Printf("Unsupported image format: %v", fileFormat)
				}
			}

			// Save processed content to file
			err = fs.SaveFile(destFile, processedData)
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
