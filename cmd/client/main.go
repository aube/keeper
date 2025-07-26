package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/aube/keeper/internal/client"
	"github.com/aube/keeper/internal/client/config"
	"github.com/aube/keeper/internal/client/infrastructure/filestore"
	"github.com/aube/keeper/internal/client/utils/helpers"
	"github.com/aube/keeper/internal/client/utils/httpclient"
	"github.com/aube/keeper/internal/client/utils/logger"
	"github.com/aube/keeper/internal/ui"
)

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

func main() {
	ctx := context.Background()

	fmt.Printf("Build version: %s\n", helpers.StringOrNA(buildVersion))
	fmt.Printf("Build date: %s\n", helpers.StringOrNA(buildTime))
	fmt.Printf("Build commit: %s\n\n", helpers.StringOrNA(buildCommit))

	var command string
	if len(os.Args) > 0 {
		command = os.Args[1]
	}

	// конфиг
	cfg := config.NewConfig()
	if cfg.Username == "" {
		log.Fatalf("Username must be set: -u <username>")
	}

	// логгер
	logger.Init(cfg.LogLevel)

	// инициализация хранилищ
	filesStoragePath := filepath.Join(cfg.StoragePath, "files", cfg.Username)
	filesRepo, err := filestore.NewFileSystemRepository(filesStoragePath)
	if err != nil {
		log.Fatalf("Failed to initialize file repository: %v", err)
	}

	tokensStoragePath := filepath.Join(cfg.StoragePath, "tokens")
	tokensRepo, err := filestore.NewFileSystemRepository(tokensStoragePath)

	if err != nil {
		log.Fatalf("Failed to initialize tokens repository: %v", err)
	}

	// инициализация http-клиента
	http := httpclient.NewHTTPClient(cfg.ServerAddress)

	app := client.NewApp(
		cfg,
		filesRepo,
		tokensRepo,
		http,
	)

	// Запускаем команду из CLI или стартуем UI
	switch command {
	case "register":
		err = app.Register(cfg.Username, cfg.Password, cfg.Email)
	case "login":
		err = app.Login(cfg.Username, cfg.Password)
	case "encrypt":
		err = app.Encrypt(cfg.Password, cfg.Input, cfg.Output)
	case "card":
		err = app.Card(cfg.Password, cfg.Number, cfg.Date, cfg.CVV)
	case "decrypt":
		err = app.Decrypt(cfg.Password, cfg.Input, cfg.Output)
	case "download":
		err = app.Download(cfg.Input)
	case "sync":
		// files4download, files4deletion, err = sync.Run(cfg, tokensRepo, filesRepo, http)
	default:
		err = ui.NewUI(app)
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
}
