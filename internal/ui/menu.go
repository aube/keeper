package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuScreen struct {
	width   int
	height  int
	choices []string
	cursor  int
	screens []string // Соответствие пунктов меню экранам
}

func NewMenuScreen() MenuScreen {
	return MenuScreen{
		choices: []string{
			"Регистрация",
			"Авторизация",
			"Зашифровать файл",
			"Расшифровать файл",
			"Удалить файл",
			"Добавить карту",
			"Удалить карту",
			"Синхронизация",
			"Выход",
		},
		screens: []string{
			"register",
			"login",
			"encrypt",
			"decrypt",
			"delete",
			"card",
			"deletecard",
			"sync",
			"",
		},
	}
}

func (m MenuScreen) Init() tea.Cmd {
	return nil
}

func (m MenuScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == len(m.choices)-1 { // Последний пункт - выход
				return m, tea.Quit
			}
			return m, func() tea.Msg {
				return SwitchScreenMsg{ScreenName: m.screens[m.cursor]}
			}
		}
	}
	return m, nil
}

func (m MenuScreen) View() string {
	title := "МЕНЮ ПРИЛОЖЕНИЯ"
	styledTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Align(lipgloss.Center).
		Bold(true).
		Render(title)

	menuItems := make([]string, len(m.choices))
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		menuItems[i] = fmt.Sprintf("%s %s", cursor, choice)
	}

	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		menuItems...,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			styledTitle,
			"",
			menu,
			"",
			"↑/↓: навигация • Enter: выбор • ctrl+c: выход",
		),
	)
}

func (m *MenuScreen) SetSize(width, height int) {
	m.width = width
	m.height = height
}
