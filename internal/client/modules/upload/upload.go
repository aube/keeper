package upload

import (
	"context"
	"encoding/json"
	"fmt"
)

type FileRepository interface {
	GetPath(filename string) string
}

type HTTPClient interface {
	UploadFile(ctx context.Context, endpoint string, filePath string, formFields map[string]string) ([]byte, error)
}

type UploadResponse struct {
	UUID string `json:"uuid"`
}

func Run(repo FileRepository, filename string, category string, http HTTPClient) error {

	ctx := context.Background()

	filepath := repo.GetPath(filename)

	postData := map[string]string{
		"description": "ololo alala",
		"category":    category,
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
