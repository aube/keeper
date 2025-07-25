package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SettingsScreen - экран настроек
type SettingsScreen struct {
	width  int
	height int
	option string
}

func NewSettingsScreen() SettingsScreen {
	return SettingsScreen{option: "option1"}
}

func (s SettingsScreen) Init() tea.Cmd {
	return nil
}

func (s SettingsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down":
			if s.option == "option1" {
				s.option = "option2"
			} else {
				s.option = "option1"
			}
		case "enter":
			fmt.Println("ololo")
		}
	}
	return s, nil
}

func (s SettingsScreen) View() string {
	option1 := "option1"
	option2 := "option2"

	if s.option == "option1" {
		option1 = "> " + option1
		option2 = "  " + option2
	} else {
		option1 = "  " + option1
		option2 = "> " + option2
	}

	content := fmt.Sprintf("Экран настроек\n\n%s\n%s\n\nНажмите 1 для возврата", option1, option2)

	return lipgloss.NewStyle().
		Width(s.width).
		Height(s.height/2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Foreground(lipgloss.Color("229")).
		Render(content)
}

func (s *SettingsScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}
