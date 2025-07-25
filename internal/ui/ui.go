package ui

import (
	"fmt"

	"github.com/aube/keeper/internal/ui/screens"
	tea "github.com/charmbracelet/bubbletea"
)

type KeeperApp interface {
	Register(Username string, Password string, Email string) error
	Login(Username string, Password string) error
	Encrypt(Password string, Input string, Output string) error
	Decrypt(Password string, Input string, Output string) error
	Upload(Output string) error
	Download(Input string) error
}

// internal/ui/ui.go
func NewUI(app KeeperApp) error {
	apis := map[string]screens.ScreenAPI{
		"main":     app.Login,
		"settings": app.Login,
		"form":     app.Login,
	}
	fmt.Println("lol")
	p := tea.NewProgram(initialModel(apis), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
