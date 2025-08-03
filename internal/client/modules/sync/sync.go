package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type FileRepository interface {
	Delete(ctx context.Context, uuid string) error
	GetPath(filename string) string
	Exists(filename string) bool
}

type TokenRepository interface {
	Save(ctx context.Context, filename string, data io.Reader) error
	GetFileContent(ctx context.Context, uuid string) (string, error)
}

type HTTPClient interface {
	Get(endpoint string, queryParams map[string]string) ([]byte, error)
	DownloadFile(fileURL, outputPath string) error
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

func Run(username string, fileRepo FileRepository, syncRepo TokenRepository, http HTTPClient) error {

	ctx := context.Background()
	now := time.Now()
	updatedAt := now.AddDate(-1, 0, 0)

	prevSync, err := syncRepo.GetFileContent(ctx, username)
	if err == nil {
		updatedAt, _ = time.Parse(time.RFC3339, prevSync)
	}

	params := make(map[string]string)
	params["deleted"] = "true"
	params["uploaded_at"] = updatedAt.Format("2006-01-02")

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

	err = syncRepo.Save(ctx, username, strings.NewReader(now.Format("2006-01-02")))
	if err != nil {
		return err
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
