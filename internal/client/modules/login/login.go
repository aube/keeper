package login

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aube/keeper/internal/client/config"
)

type FileRepository interface {
	Save(ctx context.Context, filename string, data io.Reader) error
	GetFileContent(ctx context.Context, uuid string) (string, error)
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

func Run(cfg config.EnvConfig, repo FileRepository, http HTTPClient) error {

	postData := map[string]interface{}{
		"username": cfg.Username,
		"password": cfg.Password,
	}

	response, err := http.Post("/login", postData)
	if err != nil {
		return err
	}

	token, err := ExtractToken(response)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = repo.Save(ctx, cfg.Username, strings.NewReader(token))
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
