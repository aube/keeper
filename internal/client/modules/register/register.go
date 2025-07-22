package register

import (
	"context"
	"io"

	"github.com/aube/keeper/internal/client/config"
	"github.com/aube/keeper/internal/client/entities"
)

type FileRepository interface {
	Save(ctx context.Context, file *entities.File, data io.Reader) error
	FindAll(ctx context.Context) (*entities.Files, error)
	Delete(ctx context.Context, uuid string) error
	GetFileContent(ctx context.Context, uuid string) (io.ReadCloser, error)
}

type HTTPClient interface {
	SetHeader(key, value string)
	Get(endpoint string, queryParams map[string]string) ([]byte, error)
	Post(endpoint string, body any) ([]byte, error)
	DownloadFile(fileURL, outputPath string) error
}

func Run(cfg config.EnvConfig, http HTTPClient) error {

	postData := map[string]interface{}{
		"username": cfg.Username,
		"password": cfg.Password,
		"email":    cfg.Email,
	}
	_, err := http.Post("/register", postData)
	if err != nil {
		return err
	}

	return nil
}
