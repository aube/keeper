package card

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileRepository interface {
	EncryptFile(inputPath, outputName, password string) error
}

type HTTPClient interface {
	SetHeader(key, value string)
	Get(endpoint string, queryParams map[string]string) ([]byte, error)
	Post(endpoint string, body any) ([]byte, error)
	DownloadFile(fileURL, outputPath string) error
	UploadFile(ctx context.Context, endpoint string, filePath string, formFields map[string]string) ([]byte, error)
}

type UploadResponse struct {
	UUID string `json:"uuid"`
}

type CardJSON struct {
	Number string `json:"number"`
	Date   string `json:"date"`
	CVV    string `json:"cvv"`
}

func Run(Number string, Date string, CVV string, Password string, repo FileRepository, http HTTPClient) (string, error) {

	card := &CardJSON{
		Number: Number,
		Date:   Date,
		CVV:    CVV,
	}
	filename := "card_" + strings.ReplaceAll(Number, " ", "") + ".json"
	path := filepath.Join(os.TempDir(), filename)

	cardJSON, err := json.Marshal(card)
	if err != nil {
		return "", err
	}
	var permissions os.FileMode = 0644 // Read/write for owner, read-only for others

	err = os.WriteFile(path, cardJSON, permissions)
	if err != nil {
		return "", err
	}

	err = repo.EncryptFile(path, filename, Password)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func ExtractUUID(responseBytes []byte) (string, error) {
	var uploadResp UploadResponse
	if err := json.Unmarshal(responseBytes, &uploadResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal upload response: %v", err)
	}

	return uploadResp.UUID, nil
}
