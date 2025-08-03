package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		filePath string
		fileSize int64
	}{
		{
			name:     "basic file",
			fileName: "test.txt",
			filePath: "/path/to/test.txt",
			fileSize: 1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := NewFile(tt.fileName, tt.filePath, tt.fileSize)

			assert.Equal(t, tt.fileName, file.Name)
			assert.Equal(t, tt.filePath, file.Path)
			assert.Equal(t, tt.fileSize, file.Size)
		})
	}
}

func TestFiles(t *testing.T) {
	files := Files{
		*NewFile("file1.txt", "/path/to/file1.txt", 100),
		*NewFile("file2.txt", "/path/to/file2.txt", 200),
	}

	assert.Len(t, files, 2)
	assert.Equal(t, "file1.txt", files[0].Name)
	assert.Equal(t, "file2.txt", files[1].Name)
}
