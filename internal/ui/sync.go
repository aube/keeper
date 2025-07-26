package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SyncScreen struct {
	width   int
	height  int
	api     SyncScreenAPI
	focused bool // фокус на кнопке
	pressed bool // состояние нажатия
}

func NewSyncScreen(api SyncScreenAPI) SyncScreen {
	return SyncScreen{
		api:     api,
		focused: true, // кнопка в фокусе по умолчанию
	}
}

func (b SyncScreen) Init() tea.Cmd {
	return nil
}

func (b SyncScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if b.focused {
				b.pressed = true
				// Вызываем обработчик и возвращаем команду
				return b, tea.Batch(
					func() tea.Msg {
						b.api.Sync()
						return nil
					},
				)
			}
		case "esc":
			return b, tea.Quit
		}
	}
	return b, nil
}

func (b SyncScreen) View() string {
	// Стили для кнопки
	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("62")).
		Padding(0, 3).
		Margin(1, 0)

	// Текст кнопки
	buttonText := "Нажми меня"
	if b.pressed {
		buttonText = "Нажато!"
	}

	button := buttonStyle.Render(buttonText)

	// Центрируем кнопку
	return lipgloss.Place(
		b.width, b.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"Синхронизация",
			"",
			button,
			"",
			"Enter: Нажать | ESC: Выход",
		),
	)
}

func (b *SyncScreen) SetSize(width, height int) {
	b.width = width
	b.height = height
}
