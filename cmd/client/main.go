package main

import (
	"fmt"
	"os"

	"github.com/aube/keeper/internal/client/config"
	"github.com/aube/keeper/internal/client/modules/login"
	"github.com/aube/keeper/internal/common"
)

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

func main() {
	fmt.Printf("Build version: %s\n", common.StringOrNA(buildVersion))
	fmt.Printf("Build date: %s\n", common.StringOrNA(buildTime))
	fmt.Printf("Build commit: %s\n\n", common.StringOrNA(buildCommit))

	// read config
	Config := config.NewConfig()

	if len(os.Args) == 1 {
		panic("command not found")
	}

	command := os.Args[1]

	switch command {
	case "register":
	case "login":
		login.Run(Config)
	case "encrypt":
	case "decrypt":
	case "sync":
	case "":

	}
	// start modules

}
