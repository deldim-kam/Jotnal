package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SettingsScreen экран настроек
type SettingsScreen struct {
	app  *App
	view *tview.Flex
	form *tview.Form
	info *tview.TextView
}

// NewSettingsScreen создает новый экран настроек
func NewSettingsScreen(app *App) *SettingsScreen {
	s := &SettingsScreen{
		app:  app,
		form: tview.NewForm(),
		info: tview.NewTextView().SetDynamicColors(true),
	}

	s.form.SetBorder(true).
		SetTitle(" Настройки приложения ").
		SetTitleAlign(tview.AlignLeft)

	s.info.SetBorder(true).
		SetTitle(" Информация о БД ").
		SetTitleAlign(tview.AlignLeft)

	s.view = tview.NewFlex().
		AddItem(s.form, 0, 2, true).
		AddItem(s.info, 0, 1, false)

	return s
}

func (s *SettingsScreen) Refresh() {
	s.form.Clear(true)
	cfg := s.app.GetConfigManager().Get()

	var theme, language string
	var fontSize, width, height int

	theme = cfg.Interface.Theme
	fontSize = cfg.Interface.FontSize
	width = cfg.Interface.WindowSize.Width
	height = cfg.Interface.WindowSize.Height
	language = cfg.Interface.Language

	s.form.AddInputField("Тема (dark/light):", theme, 20, nil, func(text string) {
		theme = text
	})
	s.form.AddInputField("Размер шрифта:", fmt.Sprintf("%d", fontSize), 10, nil, func(text string) {
		fmt.Sscanf(text, "%d", &fontSize)
	})
	s.form.AddInputField("Ширина окна:", fmt.Sprintf("%d", width), 10, nil, func(text string) {
		fmt.Sscanf(text, "%d", &width)
	})
	s.form.AddInputField("Высота окна:", fmt.Sprintf("%d", height), 10, nil, func(text string) {
		fmt.Sscanf(text, "%d", &height)
	})
	s.form.AddInputField("Язык (ru/en):", language, 10, nil, func(text string) {
		language = text
	})

	s.form.AddButton("Сохранить", func() {
		err := s.app.GetConfigManager().UpdateInterfaceSettings(
			theme, fontSize, width, height, language,
		)
		if err != nil {
			s.app.ShowModal("Ошибка", "Не удалось сохранить настройки: "+err.Error(), 50, 10, nil)
			return
		}

		s.app.ShowModal("Успех", "Настройки сохранены!\nПерезапустите приложение для применения изменений.", 50, 10, nil)
	})

	s.form.AddButton("Сбросить", func() {
		s.Refresh()
	})

	// Обновляем информацию о БД
	s.updateDBInfo()
}

func (s *SettingsScreen) updateDBInfo() {
	cfg := s.app.GetConfigManager().Get()
	db := s.app.GetDB()

	// Получаем количество записей в таблицах
	var projectsCount, employeesCount, snippetsCount int

	db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&projectsCount)
	db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&employeesCount)
	db.QueryRow("SELECT COUNT(*) FROM snippets").Scan(&snippetsCount)

	info := fmt.Sprintf(
		"\n[yellow]Путь к БД:[white]\n%s\n\n"+
			"[yellow]Статистика:[white]\n\n"+
			"  Проектов: %d\n"+
			"  Сотрудников: %d\n"+
			"  Сниппетов: %d\n\n"+
			"[yellow]Горячие клавиши:[white]\n\n"+
			"  [green]Ctrl+D[white] - Изменить пароль БД\n"+
			"  [green]Ctrl+P[white] - Изменить путь к БД\n",
		cfg.Database.Path,
		projectsCount, employeesCount, snippetsCount,
	)

	s.info.SetText(info)
}

func (s *SettingsScreen) GetView() tview.Primitive {
	s.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlD {
			s.changePassword()
			return nil
		}
		if event.Key() == tcell.KeyCtrlP {
			s.changePath()
			return nil
		}
		return event
	})

	return s.view
}

func (s *SettingsScreen) changePassword() {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Смена пароля БД ").SetTitleAlign(tview.AlignLeft)

	var newPassword, confirmPassword string

	form.AddPasswordField("Новый пароль:", "", 30, '*', func(text string) {
		newPassword = text
	})
	form.AddPasswordField("Подтвердите пароль:", "", 30, '*', func(text string) {
		confirmPassword = text
	})

	form.AddButton("Сменить", func() {
		if newPassword == "" {
			s.app.ShowModal("Ошибка", "Пароль не может быть пустым", 40, 8, nil)
			return
		}

		if newPassword != confirmPassword {
			s.app.ShowModal("Ошибка", "Пароли не совпадают", 40, 8, nil)
			return
		}

		s.app.ShowModal("Внимание", "Для смены пароля БД перезапустите приложение\nи используйте старый интерфейс (опция 3 в меню)", 60, 12, func() {
			s.app.pages.RemovePage("password-form")
		})
	})

	form.AddButton("Отмена", func() {
		s.app.pages.RemovePage("password-form")
	})

	s.app.pages.AddPage("password-form", center(form, 60, 12), true, true)
}

func (s *SettingsScreen) changePath() {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Изменение пути к БД ").SetTitleAlign(tview.AlignLeft)

	cfg := s.app.GetConfigManager().Get()
	var newPath string
	newPath = cfg.Database.Path

	form.AddInputField("Путь к БД:", newPath, 60, nil, func(text string) {
		newPath = text
	})

	form.AddButton("Сохранить", func() {
		if newPath == "" {
			s.app.ShowModal("Ошибка", "Путь не может быть пустым", 40, 8, nil)
			return
		}

		err := s.app.GetConfigManager().UpdateDatabasePath(newPath)
		if err != nil {
			s.app.ShowModal("Ошибка", "Не удалось сохранить путь: "+err.Error(), 50, 10, nil)
			return
		}

		s.app.ShowModal("Успех", "Путь к БД обновлен!\nПерезапустите приложение.", 50, 10, func() {
			s.app.pages.RemovePage("path-form")
		})
	})

	form.AddButton("Отмена", func() {
		s.app.pages.RemovePage("path-form")
	})

	s.app.pages.AddPage("path-form", center(form, 80, 10), true, true)
}
