package ui

import (
	"github.com/aube/keeper/internal/ui/screens"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentScreen string
	screens       map[string]tea.Model
	width         int
	height        int
	apis          map[string]screens.ScreenAPI
}

func initialModel(apis map[string]screens.ScreenAPI) model {
	m := model{
		currentScreen: "main",
		screens:       make(map[string]tea.Model),
		apis:          apis,
	}

	// Инициализация экранов
	if mainAPI, ok := apis["main"].(screens.MainScreenAPI); ok {
		m.screens["main"] = screens.NewMainScreen(mainAPI)
	}

	if settingsAPI, ok := apis["settings"].(screens.SettingsScreenAPI); ok {
		m.screens["settings"] = screens.NewSettingsScreen(settingsAPI)
	}

	return m
}

// Init - инициализация главной модели
func (m model) Init() tea.Cmd {
	// Инициализируем текущий экран
	if current, ok := m.screens[m.currentScreen]; ok {
		return current.Init()
	}
	return nil
}

// Update - обработка сообщений в главной модели
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			m.currentScreen = "main"
			cmd = m.screens[m.currentScreen].Init()
			return m, cmd
		case "2":
			m.currentScreen = "settings"
			cmd = m.screens[m.currentScreen].Init()
			return m, cmd
		case "3":
			m.currentScreen = "form"
			cmd = m.screens[m.currentScreen].Init()
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Обновляем размеры всех экранов
		for name, screen := range m.screens {
			updatedScreen, _ := screen.Update(msg)
			m.screens[name] = updatedScreen
		}
		return m, nil

	case SwitchScreenMsg:
		m.currentScreen = msg.ScreenName
		cmd = m.screens[m.currentScreen].Init()
		return m, cmd
	}

	// Делегируем обработку текущему экрану
	if current, ok := m.screens[m.currentScreen]; ok {
		var screenCmd tea.Cmd
		updatedScreen, screenCmd := current.Update(msg)
		m.screens[m.currentScreen] = updatedScreen
		cmds = append(cmds, screenCmd)
	}

	return m, tea.Batch(cmds...)
}

// View - отрисовка текущего экрана
func (m model) View() string {
	if current, ok := m.screens[m.currentScreen]; ok {
		return current.View()
	}
	return "Экран не найден"
}

// Сообщение для переключения экранов
type SwitchScreenMsg struct {
	ScreenName string
}
