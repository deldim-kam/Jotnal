package ui

import (
	"database/sql"

	"github.com/deldim-kam/Jotnal/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App представляет главное приложение с UI
type App struct {
	tviewApp      *tview.Application
	pages         *tview.Pages
	db            *sql.DB
	configManager *config.Manager

	// Экраны
	projectsScreen  *ProjectsScreen
	employeesScreen *EmployeesScreen
	snippetsScreen  *SnippetsScreen
	settingsScreen  *SettingsScreen
}

// NewApp создает новый экземпляр приложения
func NewApp(db *sql.DB, configManager *config.Manager) *App {
	app := &App{
		tviewApp:      tview.NewApplication(),
		pages:         tview.NewPages(),
		db:            db,
		configManager: configManager,
	}

	// Инициализируем экраны
	app.projectsScreen = NewProjectsScreen(app)
	app.employeesScreen = NewEmployeesScreen(app)
	app.snippetsScreen = NewSnippetsScreen(app)
	app.settingsScreen = NewSettingsScreen(app)

	// Создаем главное окно
	mainWindow := app.createMainWindow()
	app.pages.AddPage("main", mainWindow, true, true)

	return app
}

// createMainWindow создает главное окно с навигацией
func (a *App) createMainWindow() tview.Primitive {
	// Боковое меню
	menu := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	menu.SetBorder(true).
		SetTitle(" Меню ").
		SetTitleAlign(tview.AlignLeft)

	// Область контента
	content := tview.NewFlex().SetDirection(tview.FlexRow)
	content.SetBorder(true).
		SetTitle(" Добро пожаловать в Jotnal ").
		SetTitleAlign(tview.AlignCenter)

	// Текущий активный экран
	currentScreen := "welcome"

	// Функция для переключения экрана
	switchScreen := func(screenName string, screen tview.Primitive, title string) {
		if currentScreen != screenName {
			content.Clear()
			content.AddItem(screen, 0, 1, true)
			content.SetTitle(" " + title + " ")
			currentScreen = screenName
			a.tviewApp.SetFocus(screen)
		}
	}

	// Добавляем пункты меню
	menu.AddItem("📊 Проекты", "", '1', func() {
		switchScreen("projects", a.projectsScreen.GetView(), "Управление проектами")
		a.projectsScreen.Refresh()
	})

	menu.AddItem("👥 Сотрудники", "", '2', func() {
		switchScreen("employees", a.employeesScreen.GetView(), "Управление сотрудниками")
		a.employeesScreen.Refresh()
	})

	menu.AddItem("📝 Сниппеты", "", '3', func() {
		switchScreen("snippets", a.snippetsScreen.GetView(), "Библиотека сниппетов кода")
		a.snippetsScreen.Refresh()
	})

	menu.AddItem("⚙️  Настройки", "", '4', func() {
		switchScreen("settings", a.settingsScreen.GetView(), "Настройки приложения")
		a.settingsScreen.Refresh()
	})

	menu.AddItem("", "", 0, nil) // Разделитель

	menu.AddItem("❌ Выход", "", 'q', func() {
		a.tviewApp.Stop()
	})

	// Приветственное сообщение
	welcomeText := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("\n\n\n" +
			"╔═══════════════════════════════════════╗\n" +
			"║                                       ║\n" +
			"║         Jotnal IDE v1.0               ║\n" +
			"║                                       ║\n" +
			"║    Система управления проектами       ║\n" +
			"║      и сотрудниками                   ║\n" +
			"║                                       ║\n" +
			"╚═══════════════════════════════════════╝\n\n\n" +
			"Используйте цифры 1-4 для навигации\n" +
			"или выберите пункт из меню слева\n\n" +
			"Нажмите 'q' для выхода")

	content.AddItem(welcomeText, 0, 1, false)

	// Статус бар внизу
	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	cfg := a.configManager.Get()
	statusBar.SetText("[yellow]База данных:[white] " + cfg.Database.Path + " [yellow]| Тема:[white] " + cfg.Interface.Theme + " [yellow]| Язык:[white] " + cfg.Interface.Language)

	// Главный layout
	mainLayout := tview.NewFlex().
		AddItem(menu, 25, 0, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(content, 0, 1, false).
			AddItem(statusBar, 1, 0, false), 0, 1, false)

	// Глобальные горячие клавиши
	mainLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			if currentScreen == "welcome" {
				a.tviewApp.Stop()
				return nil
			}
		case '1':
			menu.SetCurrentItem(0)
			return nil
		case '2':
			menu.SetCurrentItem(1)
			return nil
		case '3':
			menu.SetCurrentItem(2)
			return nil
		case '4':
			menu.SetCurrentItem(3)
			return nil
		}
		return event
	})

	return mainLayout
}

// Run запускает приложение
func (a *App) Run() error {
	return a.tviewApp.SetRoot(a.pages, true).EnableMouse(true).Run()
}

// GetDB возвращает соединение с БД
func (a *App) GetDB() *sql.DB {
	return a.db
}

// GetConfigManager возвращает менеджер конфигурации
func (a *App) GetConfigManager() *config.Manager {
	return a.configManager
}

// ShowModal показывает модальное окно
func (a *App) ShowModal(title, message string, width, height int, doneFunc func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("modal")
			if doneFunc != nil {
				doneFunc()
			}
		})

	modal.SetTitle(" " + title + " ").SetBorder(true)

	a.pages.AddPage("modal", modal, true, true)
}

// ShowConfirm показывает диалог подтверждения
func (a *App) ShowConfirm(title, message string, yesFunc, noFunc func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Да", "Нет"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("confirm")
			if buttonIndex == 0 && yesFunc != nil {
				yesFunc()
			} else if buttonIndex == 1 && noFunc != nil {
				noFunc()
			}
		})

	modal.SetTitle(" " + title + " ").SetBorder(true)

	a.pages.AddPage("confirm", modal, true, true)
}
