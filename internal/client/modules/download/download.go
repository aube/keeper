package download

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aube/keeper/internal/client/config"
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

func Run(cfg config.EnvConfig, repo FileRepository, http HTTPClient) error {
	// ctx := context.Background()

	url := "/file?name=" + cfg.Input
	filepath := repo.GetPath(cfg.Input)

	err := http.DownloadFile(url, filepath)
	if err != nil {
		return err
	}

	fmt.Println("Файл скачан")

	return nil
}

func ExtractUUID(responseBytes []byte) (string, error) {
	var uploadResp UploadResponse
	if err := json.Unmarshal(responseBytes, &uploadResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal upload response: %v", err)
	}

	return uploadResp.UUID, nil
}
