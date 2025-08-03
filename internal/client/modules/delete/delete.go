package upload

import (
	"context"
	"encoding/json"
	"fmt"
)

type FileRepository interface {
	Delete(ctx context.Context, uuid string) error
}

type HTTPClient interface {
	Delete(filename string) error
}

type UploadResponse struct {
	UUID string `json:"uuid"`
}

func Run(repo FileRepository, filename string, http HTTPClient) error {

	ctx := context.Background()

	err := http.Delete("/delete?name" + filename)
	if err != nil {
		return err
	}

	err = repo.Delete(ctx, filename)
	if err != nil {
		return err
	}

	fmt.Println("Файл удалён")

	return nil
}

func ExtractUUID(responseBytes []byte) (string, error) {
	var uploadResp UploadResponse
	if err := json.Unmarshal(responseBytes, &uploadResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal upload response: %v", err)
	}

	return uploadResp.UUID, nil
}
