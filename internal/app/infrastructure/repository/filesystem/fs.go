package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileSystem struct {
	dir string
}

func New(dirPath string) (*FileSystem, error) {
	fi, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("file stat: %w", err)
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dirPath)
	}

	return &FileSystem{dir: dirPath}, nil
}

func (fs *FileSystem) Save(_ context.Context, data []byte, fileName string) error {
	if len(data) == 0 {
		return fmt.Errorf("data is empty")
	}

	file, err := os.OpenFile(filepath.Join(fs.dir, fileName), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("writing data to file: %w", err)
	}

	return nil
}

func (fs *FileSystem) Read(_ context.Context, fileName string) ([]byte, error) {
	filePath := filepath.Join(fs.dir, fileName)

	_, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("file stat: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading file content: %w", err)
	}

	return data, nil
}

func (fs *FileSystem) Delete(_ context.Context, fileName string) error {
	filePath := filepath.Join(fs.dir, fileName)

	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file stat: %w", err)
	}

	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("deleting the file: %w", err)
	}

	return nil
}
