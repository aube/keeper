package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// internal/ui/ui.go
func NewUI(app KeeperApp) error {
	apis := map[string]ScreenAPI{
		"login":      app,
		"register":   app,
		"encrypt":    app,
		"decrypt":    app,
		"delete":     app,
		"card":       app,
		"deletecard": app,
		"sync":       app,
	}
	fmt.Println("lol")

	p := tea.NewProgram(initialModel(apis), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

type model struct {
	currentScreen string
	screens       map[string]tea.Model
	width         int
	height        int
	apis          map[string]ScreenAPI
}

func initialModel(apis map[string]ScreenAPI) model {
	m := model{
		currentScreen: "menu",
		screens:       make(map[string]tea.Model),
		apis:          apis,
	}

	// Добавляем экран меню
	m.screens["menu"] = NewMenuScreen()

	// // Инициализация остальных экранов
	if mainAPI, ok := apis["register"].(RegisterScreenAPI); ok {
		m.screens["register"] = NewRegisterScreen(mainAPI)
	}

	if loginAPI, ok := apis["login"].(LoginScreenAPI); ok {
		m.screens["login"] = NewLoginScreen(loginAPI)
	}

	if encryptAPI, ok := apis["encrypt"].(EncryptScreenAPI); ok {
		m.screens["encrypt"] = NewEncryptScreen(encryptAPI)
	}

	if decryptAPI, ok := apis["decrypt"].(DecryptScreenAPI); ok {
		m.screens["decrypt"] = NewDecryptScreen(decryptAPI)
	}

	if deleteAPI, ok := apis["delete"].(DeleteScreenAPI); ok {
		m.screens["delete"] = NewDeleteScreen(deleteAPI)
	}

	if cardAPI, ok := apis["card"].(CardScreenAPI); ok {
		m.screens["card"] = NewCardScreen(cardAPI)
	}

	if deletecardAPI, ok := apis["deletecard"].(DeletecardScreenAPI); ok {
		m.screens["deletecard"] = NewDeletecardScreen(deletecardAPI)
	}

	if syncAPI, ok := apis["sync"].(SyncScreenAPI); ok {
		m.screens["sync"] = NewSyncScreen(syncAPI)
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
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.currentScreen = "menu"
			return m, m.screens[m.currentScreen].Init()
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
	return "Экран не найден. m.currentScreen == " + m.currentScreen
}

// Сообщение для переключения экранов
type SwitchScreenMsg struct {
	ScreenName string
}
