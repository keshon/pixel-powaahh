package pixelita

import (
	"app/internal/imagetype"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type ItemProcessStatus string

const (
	StatusUnknown ItemProcessStatus = "UNKNOWN"
	StatusError   ItemProcessStatus = "ERROR"
	StatusDone    ItemProcessStatus = "DONE"
	StatusNew     ItemProcessStatus = "NEW"
)

type ImageList struct {
	Path            string
	Format          string
	SizeAsInt       int64
	SizeAsString    string
	NewSizeAsInt    int64
	NewSizeAsString string
	NewFormat       string
	Status          ItemProcessStatus
}

type CompressionResult struct {
	index   int
	status  ItemProcessStatus
	newSize string
}

// Create a new List instance for each image file and add it to the list slice
func addFileToList(path string, format string, size int64) {
	image := ImageList{
		Path:         path,
		Format:       format,
		SizeAsInt:    size,
		SizeAsString: formatFileSize(size),
		Status:       StatusNew,
	}

	list = append(list, image)

}

// Function to format file size with appropriate unit (KB, MB, GB, etc.)
func formatFileSize(fileSize int64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
		TB = 1 << 40
		PB = 1 << 50
		EB = 1 << 60
	)

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	value := float64(fileSize)

	var i int
	for i = 0; i < len(units); i++ {
		if value < 1024 {
			break
		}
		value /= 1024
	}

	return fmt.Sprintf("%.2f %s", value, units[i])
}

func worker(id int, jobs <-chan int, results chan<- CompressionResult) {
	for i := range jobs {

		// Process the image
		var processedData []byte

		// Read uploaded file content
		uploadedData, err := fs.ReadFile(filepath.Join(conf.UploadDir, list[i].Path))
		if err != nil {
			log.Printf("Error reading image data: %v", err)
			list[i].Status = StatusUnknown
			return
		}

		// Check the file format
		fileFormat, err := imagetype.New().GetFormatByExtension(list[i].Format)
		if err != nil {
			log.Printf("Error detecting image format: %v", err)
			continue
		}

		switch fileFormat {
		case imagetype.JPEG:
			// Compress JPEG
			processedData, err = jpgEnc.Encode(uploadedData, int(jpgQuality))
			if err != nil {
				log.Printf("Error compressing image: %v", err)
				list[i].Status = StatusError
				continue
			}
		case imagetype.PNG:
			// Compress PNG
			processedData, err = pngEnc.Encode(uploadedData, int(posterization), int(minQuality), int(maxQuality), int(speed))
			if err != nil {
				log.Printf("Error compressing image: %v", err)
				list[i].Status = StatusError
				continue
			}
		default:
			log.Printf("Unsupported image format: %v", fileFormat)
		}

		// Save processed content to file
		destPath := filepath.Join(conf.ProcessedDir, list[i].Path)
		err = fs.SaveFile(destPath, processedData)
		if err != nil {
			log.Printf("Error saving converted image to file: %v", err)
		}
		list[i].Status = StatusDone

		var fileSize int64
		fileInfo, err := os.Stat(destPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fileSize = fileInfo.Size()
		list[i].NewSizeAsInt = fileSize
		list[i].NewSizeAsString = formatFileSize(fileSize)

		// Simulate some work (replace this with your image processing)
		// time.Sleep(time.Millisecond * 500)
		backend.Refresh()

		// Send the result back to the main goroutine
		results <- CompressionResult{
			index:   i,
			status:  list[i].Status,
			newSize: formatFileSize(fileSize),
		}
	}
}
