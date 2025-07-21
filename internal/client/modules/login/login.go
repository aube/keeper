package login

import (
	"fmt"

	"github.com/aube/keeper/internal/client/config"
)

func Run(config config.EnvConfig) {
	fmt.Println("config", config)
}
