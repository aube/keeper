package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MainScreen - главный экран
type MainScreen struct {
	width  int
	height int
}

func NewMainScreen() MainScreen {
	return MainScreen{}
}

func (m MainScreen) Init() tea.Cmd {
	return nil
}

func (m MainScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Можно отправить сообщение для переключения экрана
			// return m, func() tea.Msg { return SwitchScreenMsg{"settings"} }
		}
	}
	return m, nil
}

func (m MainScreen) View() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height/2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Foreground(lipgloss.Color("229")).
		Render("Главный экран\n\nНажмите 2 для перехода в настройки")

	return style
}

func (m *MainScreen) SetSize(width, height int) {
	m.width = width
	m.height = height
}
