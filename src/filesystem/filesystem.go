package filesystem

import (
	"os"
	"path/filepath"
	"pixel-powaahh/src/utilities"
)

// fetchImageFiles recursively fetches all image files from a directory and its subdirectories.
func FetchImageFiles(directory string) ([]string, error) {
	var imageFiles []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !utilities.IsImageFile(info.Name()) {
			return nil
		}

		imageFiles = append(imageFiles, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return imageFiles, nil
}

// emptyProcessedFolder deletes all files and directories within the processed folder.
func EmptyProcessedFolder(outputDirectory string) error {
	// Read the processed directory
	entries, err := os.ReadDir(outputDirectory)
	if err != nil {
		return err
	}

	// Delete all files and directories within the processed directory
	for _, entry := range entries {
		err := os.RemoveAll(filepath.Join(outputDirectory, entry.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
