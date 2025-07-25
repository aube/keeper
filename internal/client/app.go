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
	"github.com/aube/keeper/internal/client/modules/download"
	"github.com/aube/keeper/internal/client/modules/encrypt"
	"github.com/aube/keeper/internal/client/modules/login"
	"github.com/aube/keeper/internal/client/modules/register"
	"github.com/aube/keeper/internal/client/modules/upload"
)

type FileRepository interface {
	Save(ctx context.Context, filename string, data io.Reader) error
	FindAll(ctx context.Context) (*entities.Files, error)
	Delete(ctx context.Context, uuid string) error
	GetFile(ctx context.Context, uuid string) (io.ReadCloser, error)
	GetFileContent(ctx context.Context, uuid string) (string, error)
	DecryptFile(inputName, outputPath, password string) error
	EncryptFile(inputPath, outputName, password string) error
	GetPath(filename string) string
}

type TokenRepository interface {
	Save(ctx context.Context, filename string, data io.Reader) error
	GetFileContent(ctx context.Context, filename string) (string, error)
}

type HTTPClient interface {
	SetHeader(key, value string)
	Get(endpoint string, queryParams map[string]string) ([]byte, error)
	Post(endpoint string, body any) ([]byte, error)
	DownloadFile(fileURL, outputPath string) error
	UploadFile(ctx context.Context, endpoint string, filePath string, formFields map[string]string) ([]byte, error)
}

type App struct {
	filesRepo  FileRepository
	tokensRepo TokenRepository
	http       HTTPClient
	cfg        config.EnvConfig
}

func NewApp(cfg config.EnvConfig, filesRepo FileRepository, tokensRepo TokenRepository, http HTTPClient) *App {
	return &App{
		filesRepo:  filesRepo,
		tokensRepo: tokensRepo,
		http:       http,
		cfg:        cfg,
	}
}

func (a *App) Register() error {
	return register.Run(a.cfg, a.http)
}
func (a *App) Login() error {
	return login.Run(a.cfg, a.tokensRepo, a.http)
}
func (a *App) Encrypt() error {
	return encrypt.Run(a.cfg, a.filesRepo)
}
func (a *App) Decrypt() error {
	return decrypt.Run(a.cfg, a.filesRepo)
}
func (a *App) Upload() error {
	return upload.Run(a.filesRepo, a.cfg.Output, a.http)
}
func (a *App) Download() error {
	return download.Run(a.cfg, a.filesRepo, a.http)
}
func (a *App) Sync() error {
	return nil
	// return login.Run(a.cfg, a.tokensRepo, a.http)
}

func Run(command string, cfg config.EnvConfig, filesRepo FileRepository, tokensRepo TokenRepository, http HTTPClient) error {
	ctx := context.Background()

	var err error

	token, _ := tokensRepo.GetFileContent(ctx, cfg.Username)
	if token != "" {
		http.SetHeader("Authorization", "Bearer "+string(token))
	}

	switch command {
	case "register":
		err = register.Run(cfg, http)
	case "login":
		err = login.Run(cfg, tokensRepo, http)
	case "encrypt":
		err = encrypt.Run(cfg, filesRepo)
		if err == nil {
			err = upload.Run(filesRepo, cfg.Output, http)
		}
	case "decrypt":
		err = decrypt.Run(cfg, filesRepo)
	case "download":
		err = download.Run(cfg, filesRepo, http)
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
