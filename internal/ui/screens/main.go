package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// internal/ui/screens/main.go
type MainScreen struct {
	width  int
	height int
	api    MainScreenAPI
	data   MainScreenData
}

func NewMainScreen(api MainScreenAPI) MainScreen {
	return MainScreen{
		api:  api,
		data: MainScreenData{SelectedOption: "default"},
	}
}

func (m MainScreen) Init() tea.Cmd {
	return nil
}

func (m MainScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Вызываем API при действии на главном экране
			if err := m.api.OnMainAction(m.data); err != nil {
				// Обработка ошибки
				return m, tea.Quit
			}
			return m, nil
		case "up", "down":
			// Обновляем данные экрана
			if m.data.SelectedOption == "option1" {
				m.data.SelectedOption = "option2"
			} else {
				m.data.SelectedOption = "option1"
			}
		case "esc":
			// m.api.OnCancel()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MainScreen) View() string {
	// Отрисовка интерфейса с использованием m.data
	return fmt.Sprintf("Главный экран\nВыбрано: %s", m.data.SelectedOption)
}

func (m *MainScreen) SetSize(width, height int) {
	m.width = width
	m.height = height
}
