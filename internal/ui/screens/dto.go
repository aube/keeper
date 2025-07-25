package screens

// Базовый интерфейс для всех экранов
type ScreenAPI interface {
	// OnCancel()
}

// API для главного экрана
type MainScreenAPI interface {
	ScreenAPI
	OnMainAction(data MainScreenData) error
}

// API для экрана настроек
type SettingsScreenAPI interface {
	ScreenAPI
	OnSettingsSave(settings SettingsData) error
	GetCurrentSettings() SettingsData
}

// API для формы ввода
type FormScreenAPI interface {
	ScreenAPI
	OnFormSubmit(form FormData) error
	GetFormDefaults() FormData
}

type MainScreenData struct {
	SelectedOption string
}

type SettingsData struct {
	Theme         string
	Notifications bool
}

type FormData struct {
	Username string
	Email    string
	Age      int
}
