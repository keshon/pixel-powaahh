package fileio

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"pp/src/conf"
	"pp/src/imgtype"
	"strings"
)

// FileOps is an interface for file system operations related to images.
type FileOps interface {
	AddSuffix(path, suffix string) string
	ChangeExt(path, newExt string) string
	EmptyDir(outputDir string) error
	FetchFiles(directory string) ([]string, error)

	GetRelPath(filePath string) string
	ReadFile(path string) ([]byte, error)
	SaveFile(relPath string, data []byte) error
}

// FileOp implements the FileOps interface for local file system operations.
type FileOp struct{}

// NewFileOp creates a new instance of FileOp that implements the FileOps interface.
func NewFileOp() FileOps {
	return &FileOp{}
}

// AddSuffix adds a suffix to a file path.
func (fo *FileOp) AddSuffix(path, suffix string) string {
	ext := filepath.Ext(path)
	filename := strings.TrimSuffix(path, ext)
	return filename + suffix + ext
}

// ChangeExt changes the file extension of a path to a new extension.
func (fo *FileOp) ChangeExt(path, newExt string) string {
	return path[:len(path)-len(filepath.Ext(path))] + "." + strings.TrimPrefix(newExt, ".")
}

// EmptyDir deletes all files and directories within the processed folder.
func (fo *FileOp) EmptyDir(outputDir string) error {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err := os.RemoveAll(filepath.Join(outputDir, entry.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

// FetchFiles recursively fetches all image files from a directory and its subdirectories.
func (fo *FileOp) FetchFiles(directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if imgtype.IsImage(info.Name()) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Printf("error fetching image files: %v", err)
		return nil, err
	}

	// Make the file paths relative to the base directory
	for i, file := range files {
		relativePath, err := filepath.Rel(conf.UPLOAD_DIR, file)
		if err != nil {
			// Handle error if the file path is not relative to the base directory
			log.Printf("error getting relative path: %v", err)
			return nil, err
		}
		files[i] = relativePath
	}

	return files, nil
}

// GetRelPath returns the relative path of a file by removing the base directory.
func (fo *FileOp) GetRelPath(path string) string {
	baseDir := conf.UPLOAD_DIR // Base directory of the uploaded files
	relativePath, err := filepath.Rel(baseDir, path)
	if err != nil {
		log.Printf("error getting relative path: %v", err)
		return ""
	}
	return relativePath
}

func (fo *FileOp) ReadFile(path string) ([]byte, error) {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		log.Printf("error opening file: %v", err)
		return nil, err
	}
	defer f.Close()

	// Read the file data
	data, err := io.ReadAll(f)
	if err != nil {
		log.Printf("error reading file data: %v", err)
		return nil, err
	}

	return data, nil
}

// SaveFile saves the data to a file.
func (fo *FileOp) SaveFile(relPath string, data []byte) error {
	// Get the current working directory (the directory where the binary is located)
	executablePath, err := os.Executable()
	if err != nil {
		log.Printf("error getting executable path: %v", err)
		return err
	}

	// Get the directory path of the binary
	binDir := filepath.Dir(executablePath)

	// Combine the binary directory path with the relative file path
	absPath := filepath.Join(binDir, relPath)

	// Create the directory and any necessary intermediate directories
	err = os.MkdirAll(filepath.Dir(absPath), os.ModePerm)
	if err != nil {
		log.Printf("error creating directory: %v", err)
		return err
	}

	// Create the file and write the data
	file, err := os.Create(absPath)
	if err != nil {
		log.Printf("error creating file: %v", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Printf("error writing data to file: %v", err)
	}
	return err
}
