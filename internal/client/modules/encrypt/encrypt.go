package encrypt

import (
	"encoding/json"
	"errors"
	"fmt"
)

type FileRepository interface {
	EncryptFile(inputPath, outputName, password string) error
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
