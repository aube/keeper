package download

import (
	"encoding/json"
	"fmt"
)

type FileRepository interface {
	GetPath(filename string) string
}

type HTTPClient interface {
	DownloadFile(fileURL, outputPath string) error
}

type UploadResponse struct {
	UUID string `json:"uuid"`
}

func Run(inputName string, repo FileRepository, http HTTPClient) error {
	// ctx := context.Background()

	url := "/file?name=" + inputName
	filepath := repo.GetPath(inputName)

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
