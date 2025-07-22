package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aube/keeper/internal/client"
	"github.com/aube/keeper/internal/client/config"
	"github.com/aube/keeper/internal/client/infrastructure/filestore"
	"github.com/aube/keeper/internal/client/utils/httpclient"
	"github.com/aube/keeper/internal/client/utils/logger"
	"github.com/aube/keeper/internal/common"
)

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

func main() {
	// ctx := context.Background()

	fmt.Printf("Build version: %s\n", common.StringOrNA(buildVersion))
	fmt.Printf("Build date: %s\n", common.StringOrNA(buildTime))
	fmt.Printf("Build commit: %s\n\n", common.StringOrNA(buildCommit))

	if len(os.Args) == 1 {
		log.Fatalf("command not found")
	}
	command := os.Args[1]

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

	// запуск приложения
	client.Run(
		command,
		cfg,
		filesRepo,
		tokensRepo,
		http,
	)

}
