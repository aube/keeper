package client

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/aube/keeper/internal/client/config"
	"github.com/aube/keeper/internal/client/entities"
	"github.com/aube/keeper/internal/client/modules/decrypt"
	"github.com/aube/keeper/internal/client/modules/encrypt"
	"github.com/aube/keeper/internal/client/modules/login"
	"github.com/aube/keeper/internal/client/modules/register"
)

type FileRepository interface {
	Save(ctx context.Context, filename string, data io.Reader) error
	FindAll(ctx context.Context) (*entities.Files, error)
	Delete(ctx context.Context, uuid string) error
	GetFileContent(ctx context.Context, uuid string) (io.ReadCloser, error)
	DecryptFile(inputName, outputPath, password string) error
	EncryptFile(inputPath, outputName, password string) error
}

type TokenRepository interface {
	FileRepository
}

type HTTPClient interface {
	SetHeader(key, value string)
	Get(endpoint string, queryParams map[string]string) ([]byte, error)
	Post(endpoint string, body any) ([]byte, error)
	DownloadFile(fileURL, outputPath string) error
}

func Run(command string, cfg config.EnvConfig, filesRepo FileRepository, tokensRepo TokenRepository, http HTTPClient) error {
	ctx := context.Background()

	var err error

	switch command {
	case "register":
		err = register.Run(cfg, http)
	case "login":
		err = login.Run(cfg, tokensRepo, http)
	case "encrypt":
		err = encrypt.Run(cfg, filesRepo)
		// if (err != nil) {
		// 	upload(cfg, filesRepo, http)
		// }
	case "decrypt":
		err = decrypt.Run(cfg, filesRepo)
	case "sync":
		// files4download, files4deletion, err = sync.Run(cfg, tokensRepo, filesRepo, http)
	case "":
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		ctx.Done()
		os.Exit(1)
	}()

	return nil
}
