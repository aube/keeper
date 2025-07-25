package decrypt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/aube/keeper/internal/client/entities"
)

type FileRepository interface {
	Save(ctx context.Context, filename string, data io.Reader) error
	FindAll(ctx context.Context) (*entities.Files, error)
	Delete(ctx context.Context, uuid string) error
	GetFileContent(ctx context.Context, uuid string) (string, error)
	DecryptFile(inputName, outputPath, password string) error
}

type HTTPClient interface {
	SetHeader(key, value string)
	Get(endpoint string, queryParams map[string]string) ([]byte, error)
	Post(endpoint string, body any) ([]byte, error)
	DownloadFile(fileURL, outputPath string) error
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Run(password string, inputName string, outputPath string, repo FileRepository) error {

	if password == "" {
		return errors.New("empty password")
	}
	if inputName == "" {
		return errors.New("empty input file name")
	}
	if outputPath == "" {
		return errors.New("empty output file path")
	}

	err := repo.DecryptFile(inputName, outputPath, password)
	if err != nil {
		return err
	}

	return nil
}

func ExtractToken(responseBytes []byte) (string, error) {
	var tokenResp LoginResponse
	if err := json.Unmarshal(responseBytes, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal token response: %v", err)
	}

	return tokenResp.Token, nil
}
