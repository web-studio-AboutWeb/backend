package filesystem

import (
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
		tempFile  string
		tempDir   string
	}{
		{
			name:      "path does not exists",
			dirPath:   "unknown_dir",
			wantError: true,
		},
		{
			name:      "path is not a directory",
			dirPath:   "temp_dir/temp_file.txt",
			wantError: true,
			tempDir:   "temp_dir",
			tempFile:  "temp_dir/temp_file.txt",
		},
		{
			name:      "should pass",
			dirPath:   "temp_dir",
			wantError: false,
			tempDir:   "temp_dir",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			if tc.tempDir != "" {
				err := os.Mkdir(tc.tempDir, 0744)
				require.NoError(tt, err)
				defer os.RemoveAll(tc.tempDir)
			}

			if tc.tempFile != "" {
				f, err := os.Create(tc.tempFile)
				require.NoError(tt, err)
				defer f.Close()
				defer os.Remove(tc.dirPath)
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

	err := os.Mkdir(tempDir, 0744)
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	f, err := os.Create(filepath.Join(tempDir, tempFile))
	require.NoError(t, err)
	defer f.Close()

	_, err = f.Write([]byte(fileContent))
	require.NoError(t, err)

	fs, err := New(tempDir)
	require.NoError(t, err)

	content, err := fs.Read(tempFile)
	require.NoError(t, err)

	require.Equal(t, fileContent, string(content))
}

func TestFileSystem_Save(t *testing.T) {
	tempDir := "temp_dir"
	fileContent := "hello, world!\n"

	err := os.Mkdir(tempDir, 0744)
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fs, err := New(tempDir)
	require.NoError(t, err)

	fileName, err := fs.Save([]byte(fileContent), "txt")
	require.NoError(t, err)

	content, err := fs.Read(fileName)
	require.NoError(t, err)

	require.Equal(t, fileContent, string(content))
}

func TestFileSystem_Delete(t *testing.T) {
	tempDir := "temp_dir"

	err := os.Mkdir(tempDir, 0744)
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fs, err := New(tempDir)
	require.NoError(t, err)

	fileName, err := fs.Save([]byte("123"), "txt")
	require.NoError(t, err)

	err = fs.Delete(fileName)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(tempDir, fileName))
	require.Error(t, err)
	require.ErrorIs(t, err, os.ErrNotExist)
}
