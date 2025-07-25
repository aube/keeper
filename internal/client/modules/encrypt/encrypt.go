package encrypt

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
	EncryptFile(inputPath, outputName, password string) error
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

func Run(password string, inputPath string, outputName string, repo FileRepository) error {
	if password == "" {
		return errors.New("empty password")
	}
	if inputPath == "" {
		return errors.New("empty input file path")
	}
	if outputName == "" {
		return errors.New("empty output file name")
	}

	err := repo.EncryptFile(inputPath, outputName, password)
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
