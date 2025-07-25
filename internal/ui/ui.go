package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeeperApp interface {
	Register(Username string, Password string, Email string) error
	Login(Username string, Password string) error
	Encrypt(Password string, Input string, Output string) error
	Decrypt(Password string, Input string, Output string) error
	Upload(Output string) error
	Download(Input string) error
}

func NewUI(app KeeperApp) {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ошибка: %v", err)
		os.Exit(1)
	}
}

// Главная модель UI
type model struct {
	currentScreen string               // Текущий экран
	screens       map[string]tea.Model // Все экраны
	width         int
	height        int
}

func initialModel() model {
	m := model{
		currentScreen: "main",
		screens:       make(map[string]tea.Model),
	}

	m.screens["main"] = NewMainScreen()
	m.screens["settings"] = NewSettingsScreen()
	// m.screens["help"] = NewHelpScreen()

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			m.currentScreen = "main"
			return m, nil
		case "2":
			m.currentScreen = "settings"
			return m, nil
		case "3":
			m.currentScreen = "help"
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Передаем размеры всем экранам
		for _, screen := range m.screens {
			if sz, ok := screen.(interface{ SetSize(int, int) }); ok {
				sz.SetSize(msg.Width, msg.Height)
			}
		}
		return m, nil
	}

	current, cmd := m.screens[m.currentScreen].Update(msg)
	m.screens[m.currentScreen] = current
	return m, cmd
}

func (m model) View() string {
	// Получаем view текущего экрана
	screenView := m.screens[m.currentScreen].View()

	// Добавляем навигацию
	nav := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		MarginTop(1).
		Render(
			"1: Главная | 2: Настройки | 3: Помощь | q: Выход",
		)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		screenView,
		nav,
	)
}
