package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		dirPath   string
		wantError bool
		tempFile  bool
	}{
		{
			name:      "path does not exists",
			dirPath:   "temp_dir",
			wantError: false,
		},
		{
			name:      "path is not a directory",
			dirPath:   "temp_dir/test.txt",
			wantError: true,
			tempFile:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			os.Mkdir(filepath.Dir(tc.dirPath), os.ModePerm)
			defer os.RemoveAll(tc.dirPath)

			if tc.tempFile {
				f, err := os.Create(tc.dirPath)
				require.NoError(tt, err)
				defer f.Close()
			}

			_, err := New(tc.dirPath)
			if tc.wantError {
				require.Error(tt, err)
				return
			}

			require.NoError(tt, err)
		})
	}
}

func TestFileSystem_Read(t *testing.T) {
	tempDir := "temp_dir"
	tempFile := "temp_file.txt"
	fileContent := "hello, world!\n"

	_ = os.Mkdir(tempDir, 0744)
	defer os.RemoveAll(tempDir)

	f, err := os.Create(filepath.Join(tempDir, tempFile))
	require.NoError(t, err)
	defer f.Close()

	_, err = f.Write([]byte(fileContent))
	require.NoError(t, err)

	fs, err := New(tempDir)
	require.NoError(t, err)

	content, err := fs.Read(context.Background(), tempFile)
	require.NoError(t, err)

	require.Equal(t, fileContent, string(content))
}

func TestFileSystem_Save(t *testing.T) {
	tempDir := "temp_dir"
	fileContent := "hello, world!\n"

	_ = os.Mkdir(tempDir, 0744)
	defer os.RemoveAll(tempDir)

	fs, err := New(tempDir)
	require.NoError(t, err)

	fileName := "test.txt"

	err = fs.Save(context.Background(), []byte(fileContent), fileName)
	require.NoError(t, err)

	content, err := fs.Read(context.Background(), fileName)
	require.NoError(t, err)

	require.Equal(t, fileContent, string(content))
}

func TestFileSystem_Delete(t *testing.T) {
	tempDir := "temp_dir"

	_ = os.Mkdir(tempDir, 0744)
	defer os.RemoveAll(tempDir)

	fs, err := New(tempDir)
	require.NoError(t, err)

	fileName := "test.txt"

	err = fs.Save(context.Background(), []byte("123"), fileName)
	require.NoError(t, err)

	err = fs.Delete(context.Background(), fileName)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(tempDir, fileName))
	require.Error(t, err)
	require.ErrorIs(t, err, os.ErrNotExist)
}
