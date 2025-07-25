package client

import (
	"context"
	"io"

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

type KeeperApp interface {
	Register(Username string, Password string, Email string) error
	Login(Username string, Password string) error
	Encrypt(Password string, Input string, Output string) error
	Decrypt(Password string, Input string, Output string) error
	Upload(Output string) error
	Download(Input string) error
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

func (a *App) Register(Username string, Password string, Email string) error {
	return register.Run(Username, Password, Email, a.http)
}
func (a *App) Login(Username string, Password string) error {
	return login.Run(Username, Password, a.tokensRepo, a.http)
}
func (a *App) Encrypt(Password string, Input string, Output string) error {
	return encrypt.Run(Password, Input, Output, a.filesRepo)
}
func (a *App) Decrypt(Password string, Input string, Output string) error {
	return decrypt.Run(a.cfg.Password, a.cfg.Input, a.cfg.Output, a.filesRepo)
}
func (a *App) Upload(Output string) error {
	return upload.Run(a.filesRepo, Output, a.http)
}
func (a *App) Download(Input string) error {
	return download.Run(Input, a.filesRepo, a.http)
}
func (a *App) Sync() error {
	return nil
	// return login.Run(a.cfg, a.tokensRepo, a.http)
}
