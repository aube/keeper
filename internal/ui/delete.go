// internal/ui/screens/auth.go
package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DeleteScreen struct {
	width    int
	height   int
	api      DeleteScreenAPI
	inputs   []textinput.Model
	focus    int
	errorMsg string
}

func NewDeleteScreen(api DeleteScreenAPI) DeleteScreen {
	a := DeleteScreen{
		api:    api,
		inputs: make([]textinput.Model, 1),
	}

	// Поле логина
	a.inputs[0] = textinput.New()
	a.inputs[0].Placeholder = "Имя зашифрованного файла"
	a.inputs[0].CharLimit = 32
	a.inputs[0].Focus()
	a.inputs[0].Prompt = "┃ "
	a.inputs[0].TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	// a.inputs[0].SetValue(api.GetLastLogin())

	return a
}

func (a DeleteScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (a DeleteScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			// Обработка навигации между полями
			s := msg.String()

			if s == "enter" && a.focus == len(a.inputs)-1 {
				// Отправка формы при нажатии Enter в поле пароля
				filename := a.inputs[0].Value()

				if err := a.api.Delete(filename); err != nil {
					a.errorMsg = err.Error()
					return a, nil
				}

				return a, func() tea.Msg {
					return SwitchScreenMsg{ScreenName: "menu"}
				}
			}

			// Циклическая навигация между полями
			if s == "up" || s == "shift+tab" {
				a.focus--
			} else {
				a.focus++
			}

			if a.focus >= len(a.inputs) {
				a.focus = 0
			} else if a.focus < 0 {
				a.focus = len(a.inputs) - 1
			}

			// Устанавливаем фокус на текущее поле
			cmds = make([]tea.Cmd, len(a.inputs))
			for i := range a.inputs {
				if i == a.focus {
					cmds[i] = a.inputs[i].Focus()
				} else {
					a.inputs[i].Blur()
				}
			}
			return a, tea.Batch(cmds...)

		case "esc":

			return a, tea.Quit
		}
	}

	// Обновляем текущее поле ввода
	var cmd tea.Cmd
	a.inputs[a.focus], cmd = a.inputs[a.focus].Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a DeleteScreen) View() string {
	title := "Авторизация"
	styledTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Align(lipgloss.Center).
		Bold(true).
		Render(title)

	// Стили для полей ввода
	inputStyle := lipgloss.NewStyle().
		Width(30).
		Padding(0, 1)

	// Собираем поля ввода
	inputs := []string{
		inputStyle.Render(a.inputs[0].View()),
	}

	// Добавляем подписи
	form := lipgloss.JoinVertical(
		lipgloss.Left,
		a.inputs[0].Placeholder+":",
		inputs[0],
	)

	// Добавляем сообщение об ошибке
	if a.errorMsg != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render("Ошибка: " + a.errorMsg)
		form = lipgloss.JoinVertical(lipgloss.Left, form, "", errorStyle)
	}

	// Кнопка отправки
	submitBtn := " "
	if a.focus == len(a.inputs)-1 {
		submitBtn = ">"
	}
	submit := lipgloss.NewStyle().
		MarginTop(1).
		Render(fmt.Sprintf("%s Войти (Enter)", submitBtn))

	// Возврат в меню
	back := lipgloss.NewStyle().
		MarginTop(1).
		Render("ESC: Отмена")

	return lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			styledTitle,
			"",
			form,
			"",
			submit,
			back,
		),
	)
}

func (a *DeleteScreen) SetSize(width, height int) {
	a.width = width
	a.height = height
}
