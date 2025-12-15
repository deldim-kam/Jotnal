package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/deldim-kam/Jotnal/internal/config"
	"github.com/deldim-kam/Jotnal/internal/database"
	"github.com/deldim-kam/Jotnal/internal/ui"
	"golang.org/x/term"
)

func main() {
	fmt.Println("=== Jotnal IDE ===")
	fmt.Println("Запуск приложения...")

	// Инициализация конфигурации
	cfgManager, err := config.NewManager("")
	if err != nil {
		log.Fatalf("Ошибка при инициализации конфигурации: %v", err)
	}

	cfg := cfgManager.Get()
	fmt.Printf("Конфигурация загружена из: %s/.jotnal/config.json\n", mustGetHomeDir())

	// Проверяем наличие пароля БД
	if cfg.Database.Password == "" {
		fmt.Println("\nПервый запуск: необходимо установить пароль для базы данных")
		password := promptPassword()
		if err := cfgManager.UpdateDatabasePassword(password); err != nil {
			log.Fatalf("Ошибка при сохранении пароля: %v", err)
		}
		cfg = cfgManager.Get()
		fmt.Println("Пароль успешно установлен!")
	}

	// Автоматическое подключение к БД
	fmt.Printf("\nПодключение к базе данных: %s\n", cfg.Database.Path)
	dbManager, err := database.NewManager(cfg.Database.Path, cfg.Database.Password)
	if err != nil {
		log.Fatalf("Ошибка при создании менеджера БД: %v", err)
	}

	if err := dbManager.Connect(); err != nil {
		log.Fatalf("Ошибка при подключении к БД: %v", err)
	}
	defer dbManager.Close()

	fmt.Printf("✓ Успешно подключено к базе данных\n")
	fmt.Printf("✓ Версия схемы БД: %d\n", dbManager.GetVersion())

	// Выбор интерфейса
	fmt.Println("\n=== Выбор интерфейса ===")
	fmt.Println("1. Графический интерфейс (TUI) - рекомендуется")
	fmt.Println("2. Текстовое меню (старый интерфейс)")
	fmt.Print("\nВыберите интерфейс (1-2, Enter для графического): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "2" {
		// Старый текстовый интерфейс
		showMenu(cfgManager, dbManager)
	} else {
		// Новый графический интерфейс
		fmt.Println("\nЗапуск графического интерфейса...")
		app := ui.NewApp(dbManager.GetDB(), cfgManager)
		if err := app.Run(); err != nil {
			log.Fatalf("Ошибка при запуске UI: %v", err)
		}
	}
}

func showMenu(cfgManager *config.Manager, dbManager *database.Manager) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Главное меню ===")
		fmt.Println("1. Показать информацию о БД")
		fmt.Println("2. Изменить путь к БД")
		fmt.Println("3. Изменить пароль БД")
		fmt.Println("4. Показать настройки интерфейса")
		fmt.Println("5. Изменить настройки интерфейса")
		fmt.Println("6. Выход")
		fmt.Print("\nВыберите действие (1-6): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			showDatabaseInfo(dbManager, cfgManager)
		case "2":
			changeDatabasePath(cfgManager, dbManager)
		case "3":
			changeDatabasePassword(cfgManager, dbManager)
		case "4":
			showInterfaceSettings(cfgManager)
		case "5":
			changeInterfaceSettings(cfgManager)
		case "6":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неверный выбор, попробуйте снова")
		}
	}
}

func showDatabaseInfo(dbManager *database.Manager, cfgManager *config.Manager) {
	cfg := cfgManager.Get()
	fmt.Println("\n=== Информация о базе данных ===")
	fmt.Printf("Путь: %s\n", cfg.Database.Path)
	fmt.Printf("Версия схемы: %d\n", dbManager.GetVersion())

	// Проверяем размер файла БД
	if info, err := os.Stat(cfg.Database.Path); err == nil {
		fmt.Printf("Размер файла: %.2f КБ\n", float64(info.Size())/1024)
	}

	// Показываем список таблиц
	db := dbManager.GetDB()
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		fmt.Printf("Ошибка при получении списка таблиц: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("\nТаблицы в БД:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err == nil {
			fmt.Printf("  - %s\n", tableName)
		}
	}
}

func changeDatabasePath(cfgManager *config.Manager, dbManager *database.Manager) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nВведите новый путь к БД (или Enter для отмены): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println("Отменено")
		return
	}

	if err := cfgManager.UpdateDatabasePath(input); err != nil {
		fmt.Printf("Ошибка при обновлении пути: %v\n", err)
		return
	}

	fmt.Println("✓ Путь к БД обновлен. Перезапустите приложение для применения изменений.")
}

func changeDatabasePassword(cfgManager *config.Manager, dbManager *database.Manager) {
	fmt.Println("\n=== Смена пароля базы данных ===")
	newPassword := promptPassword()

	if err := dbManager.ChangePassword(newPassword); err != nil {
		fmt.Printf("Ошибка при смене пароля БД: %v\n", err)
		return
	}

	if err := cfgManager.UpdateDatabasePassword(newPassword); err != nil {
		fmt.Printf("Ошибка при сохранении пароля в конфигурации: %v\n", err)
		return
	}

	fmt.Println("✓ Пароль успешно изменен!")
}

func showInterfaceSettings(cfgManager *config.Manager) {
	cfg := cfgManager.Get()
	fmt.Println("\n=== Настройки интерфейса ===")
	fmt.Printf("Тема: %s\n", cfg.Interface.Theme)
	fmt.Printf("Размер шрифта: %d\n", cfg.Interface.FontSize)
	fmt.Printf("Размер окна: %dx%d\n", cfg.Interface.WindowSize.Width, cfg.Interface.WindowSize.Height)
	fmt.Printf("Язык: %s\n", cfg.Interface.Language)
}

func changeInterfaceSettings(cfgManager *config.Manager) {
	reader := bufio.NewReader(os.Stdin)
	cfg := cfgManager.Get()

	fmt.Println("\n=== Изменение настроек интерфейса ===")

	fmt.Printf("Тема (текущая: %s, например: dark/light): ", cfg.Interface.Theme)
	theme, _ := reader.ReadString('\n')
	theme = strings.TrimSpace(theme)
	if theme == "" {
		theme = cfg.Interface.Theme
	}

	fmt.Printf("Размер шрифта (текущий: %d): ", cfg.Interface.FontSize)
	var fontSize int
	fmt.Fscanf(reader, "%d\n", &fontSize)
	if fontSize == 0 {
		fontSize = cfg.Interface.FontSize
	}

	fmt.Printf("Ширина окна (текущая: %d): ", cfg.Interface.WindowSize.Width)
	var width int
	fmt.Fscanf(reader, "%d\n", &width)
	if width == 0 {
		width = cfg.Interface.WindowSize.Width
	}

	fmt.Printf("Высота окна (текущая: %d): ", cfg.Interface.WindowSize.Height)
	var height int
	fmt.Fscanf(reader, "%d\n", &height)
	if height == 0 {
		height = cfg.Interface.WindowSize.Height
	}

	fmt.Printf("Язык (текущий: %s, например: ru/en): ", cfg.Interface.Language)
	language, _ := reader.ReadString('\n')
	language = strings.TrimSpace(language)
	if language == "" {
		language = cfg.Interface.Language
	}

	if err := cfgManager.UpdateInterfaceSettings(theme, fontSize, width, height, language); err != nil {
		fmt.Printf("Ошибка при обновлении настроек: %v\n", err)
		return
	}

	fmt.Println("✓ Настройки успешно обновлены!")
}

func promptPassword() string {
	fmt.Print("Введите пароль для базы данных: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("Ошибка при чтении пароля: %v", err)
	}
	fmt.Println()

	fmt.Print("Повторите пароль: ")
	password2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("Ошибка при чтении пароля: %v", err)
	}
	fmt.Println()

	if string(password) != string(password2) {
		fmt.Println("Пароли не совпадают, попробуйте снова")
		return promptPassword()
	}

	return string(password)
}

func mustGetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~"
	}
	return home
}
