package sync

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
	Exists(filename string) bool
}

type TokenRepository interface {
	Save(ctx context.Context, filename string, data io.Reader) error
	GetFileContent(ctx context.Context, uuid string) (string, error)
}

type HTTPClient interface {
	SetHeader(key, value string)
	Get(endpoint string, queryParams map[string]string) ([]byte, error)
	Post(endpoint string, body any) ([]byte, error)
	DownloadFile(fileURL, outputPath string) error
	UploadFile(ctx context.Context, endpoint string, filePath string, formFields map[string]string) ([]byte, error)
}

type Response struct {
	Pagination Pagination `json:"pagination"`
	Rows       []Row      `json:"rows"`
}

type Pagination struct {
	Size  int `json:"size"`
	Page  int `json:"page"`
	Total int `json:"total"`
}

type Row struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Size        int    `json:"size"`
	ContentType string `json:"content_type"`
	Description string `json:"description"`
}

func Run(fileRepo FileRepository, syncRepo TokenRepository, http HTTPClient) error {

	ctx := context.Background()

	params := make(map[string]string)
	params["deleted"] = "true"

	deletedFilesResponse, err := http.Get("/uploads", params)
	if err != nil {
		return err
	}

	// local delete files marks deleted on server
	deletedFiles, err := ExtractRows(deletedFilesResponse)
	if err != nil {
		return err
	}

	for _, row := range deletedFiles {
		if !fileRepo.Exists(row.Name) {
			continue
		}
		err = fileRepo.Delete(ctx, row.Name)
		if err != nil {
			return err
		}
	}

	// download new files
	params["deleted"] = "false"
	newFilesResponse, err := http.Get("/uploads", params)
	if err != nil {
		return err
	}
	newFiles, err := ExtractRows(newFilesResponse)
	if err != nil {
		return err
	}

	for _, row := range newFiles {
		if fileRepo.Exists(row.Name) {
			continue
		}
		url := "/file?name=" + row.Name
		filepath := fileRepo.GetPath(row.Name)

		err := http.DownloadFile(url, filepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExtractRows(responseBytes []byte) ([]Row, error) {
	var response Response
	err := json.Unmarshal([]byte(responseBytes), &response)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}

	return response.Rows, nil
}
