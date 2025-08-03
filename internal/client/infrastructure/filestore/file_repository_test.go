package filestore

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSystemRepository(t *testing.T) {
	// Setup temp directory
	tempDir, err := os.MkdirTemp("", "filestore-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo, err := NewFileSystemRepository(tempDir)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("Save and GetFile", func(t *testing.T) {
		testContent := "test file content"
		filename := "testfile.txt"

		err := repo.Save(ctx, filename, strings.NewReader(testContent))
		assert.NoError(t, err)

		// Test GetFile
		reader, err := repo.GetFile(ctx, filename)
		assert.NoError(t, err)
		defer reader.Close()

		content, err := io.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, testContent, string(content))

		// Test GetFileContent
		contentStr, err := repo.GetFileContent(ctx, filename)
		assert.NoError(t, err)
		assert.Equal(t, testContent, contentStr)

		// Test Exists
		assert.True(t, repo.Exists(filename))
		assert.False(t, repo.Exists("nonexistent.txt"))
	})

	t.Run("Delete", func(t *testing.T) {
		filename := "todelete.txt"
		err := repo.Save(ctx, filename, strings.NewReader("content"))
		assert.NoError(t, err)

		err = repo.Delete(ctx, filename)
		assert.NoError(t, err)

		_, err = repo.GetFile(ctx, filename)
		assert.Error(t, err)
	})

	t.Run("FindAll", func(t *testing.T) {
		// Clean up before test
		os.RemoveAll(tempDir)
		os.MkdirAll(tempDir, 0755)

		files := []string{"file1.txt", "file2.txt", "file3.txt"}
		for _, f := range files {
			err := repo.Save(ctx, f, strings.NewReader("content"))
			assert.NoError(t, err)
		}

		result, err := repo.FindAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, *result, 3)

		for _, f := range *result {
			assert.Contains(t, files, f.Name)
			assert.Equal(t, filepath.Join(tempDir, f.Name), f.Path)
		}
	})

	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		// Create a test file
		testContent := "this is a secret message"
		inputFile := filepath.Join(tempDir, "plaintext.txt")
		err := os.WriteFile(inputFile, []byte(testContent), 0644)
		assert.NoError(t, err)

		// Encrypt
		password := "securepassword123"
		encryptedName := "encrypted.dat"
		err = repo.EncryptFile(inputFile, encryptedName, password)
		assert.NoError(t, err)

		// Decrypt to new file
		outputFile := filepath.Join(tempDir, "decrypted.txt")
		err = repo.DecryptFile(encryptedName, outputFile, password)
		assert.NoError(t, err)

		// Verify decrypted content
		decryptedContent, err := os.ReadFile(outputFile)
		assert.NoError(t, err)
		assert.Equal(t, testContent, string(decryptedContent))

		// Test wrong password
		wrongOutput := filepath.Join(tempDir, "wrong.txt")
		err = repo.DecryptFile(encryptedName, wrongOutput, "wrongpassword")
		assert.Error(t, err)
	})
}
