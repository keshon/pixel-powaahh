package filesystem

import (
	"app/internal/config"
	"app/internal/imagetype"
	"log"
	"os"
	"path/filepath"
)

// FileSystem is an interface for file system operations related to images.
type FileSystem interface {
	AddSuffixToFileName(filePath, suffix string) string
	ChangeFileExtension(filePath, newExtension string) string
	ClearDirectory(directoryPath string) error
	GetImageFiles(directoryPath string) ([]string, error)
	GetRelativePath(filePath string) string
	ReadFile(filePath string) ([]byte, error)
	SaveFile(relativePath string, data []byte) error
}

// FileSystemImpl implements the FileSystem interface for local file system operations.
type FileSystemImpl struct {
	config *config.Config
}

// NewFileSystemImpl creates a new instance of FileSystemImpl that implements the FileSystem interface.
func NewFileSystemImpl(config *config.Config) FileSystem {
	return &FileSystemImpl{config: config}
}

func (fs *FileSystemImpl) AddSuffixToFileName(filePath, suffix string) string {
	ext := filepath.Ext(filePath)
	filename := filePath[:len(filePath)-len(ext)]
	return filename + suffix + ext
}

func (fs *FileSystemImpl) ChangeFileExtension(filePath, newExtension string) string {
	return filePath[:len(filePath)-len(filepath.Ext(filePath))] + "." + newExtension
}

func (fs *FileSystemImpl) ClearDirectory(directoryPath string) error {
	return os.RemoveAll(directoryPath)
}

func (fs *FileSystemImpl) GetImageFiles(directoryPath string) ([]string, error) {
	var imageFiles []string

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if imagetype.New().IsSupportedExtension(info.Name()) {
			imageFiles = append(imageFiles, path)
		}

		return nil
	})

	return imageFiles, err
}

func (fs *FileSystemImpl) GetRelativePath(filePath string) string {
	baseDir := fs.config.UploadDir
	relativePath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		log.Printf("error getting relative path: %v", err)
		return ""
	}
	return relativePath
}

func (fs *FileSystemImpl) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func (fs *FileSystemImpl) SaveFile(relativePath string, data []byte) error {
	absPath := filepath.Join(fs.config.BinDir, relativePath)

	if err := os.MkdirAll(filepath.Dir(absPath), os.ModePerm); err != nil {
		log.Printf("error creating directory: %v", err)
		return err
	}

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
