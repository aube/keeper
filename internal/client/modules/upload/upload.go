package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aube/keeper/internal/client/entities"
)

type FileRepository interface {
	Save(ctx context.Context, filename string, data io.Reader) error
	FindAll(ctx context.Context) (*entities.Files, error)
	Delete(ctx context.Context, uuid string) error
	GetFileContent(ctx context.Context, uuid string) (string, error)
	GetPath(filename string) string
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

func Run(repo FileRepository, filename string, http HTTPClient) error {

	ctx := context.Background()

	filepath := repo.GetPath(filename)

	postData := map[string]string{
		"description": "ololo alala",
	}

	responce, err := http.UploadFile(ctx, "/upload", filepath, postData)
	if err != nil {
		return err
	}

	UUID, err := ExtractUUID(responce)
	if err != nil {
		return err
	}

	fmt.Println("Файл загружен под именем:", UUID)

	return nil
}

func ExtractUUID(responseBytes []byte) (string, error) {
	var uploadResp UploadResponse
	if err := json.Unmarshal(responseBytes, &uploadResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal upload response: %v", err)
	}

	return uploadResp.UUID, nil
}
