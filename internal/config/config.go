package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config представляет конфигурацию приложения
type Config struct {
	Database  DatabaseConfig  `json:"database"`
	Interface InterfaceConfig `json:"interface"`
}

// DatabaseConfig содержит настройки базы данных
type DatabaseConfig struct {
	Path     string `json:"path"`
	Password string `json:"password"`
}

// InterfaceConfig содержит настройки интерфейса
type InterfaceConfig struct {
	Theme      string `json:"theme"`
	FontSize   int    `json:"font_size"`
	WindowSize struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"window_size"`
	Language string `json:"language"`
}

// Manager управляет конфигурацией приложения
type Manager struct {
	configPath string
	config     *Config
}

// NewManager создает новый менеджер конфигурации
func NewManager(configPath string) (*Manager, error) {
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(homeDir, ".jotnal", "config.json")
	}

	m := &Manager{
		configPath: configPath,
	}

	// Создаем директорию если не существует
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, err
	}

	// Загружаем или создаем конфигурацию
	if err := m.Load(); err != nil {
		// Если файл не существует, создаем конфигурацию по умолчанию
		if os.IsNotExist(err) {
			m.config = m.defaultConfig()
			if err := m.Save(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return m, nil
}

// defaultConfig возвращает конфигурацию по умолчанию
func (m *Manager) defaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	defaultDBPath := filepath.Join(homeDir, ".jotnal", "jotnal.db")

	cfg := &Config{
		Database: DatabaseConfig{
			Path:     defaultDBPath,
			Password: "", // Пароль будет установлен при первом запуске
		},
		Interface: InterfaceConfig{
			Theme:    "dark",
			FontSize: 14,
			Language: "ru",
		},
	}

	cfg.Interface.WindowSize.Width = 1280
	cfg.Interface.WindowSize.Height = 720

	return cfg
}

// Load загружает конфигурацию из файла
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	m.config = &Config{}
	return json.Unmarshal(data, m.config)
}

// Save сохраняет конфигурацию в файл
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0600)
}

// Get возвращает текущую конфигурацию
func (m *Manager) Get() *Config {
	return m.config
}

// UpdateDatabasePath обновляет путь к базе данных
func (m *Manager) UpdateDatabasePath(path string) error {
	m.config.Database.Path = path
	return m.Save()
}

// UpdateDatabasePassword обновляет пароль базы данных
func (m *Manager) UpdateDatabasePassword(password string) error {
	m.config.Database.Password = password
	return m.Save()
}

// UpdateInterfaceSettings обновляет настройки интерфейса
func (m *Manager) UpdateInterfaceSettings(theme string, fontSize int, width, height int, language string) error {
	m.config.Interface.Theme = theme
	m.config.Interface.FontSize = fontSize
	m.config.Interface.WindowSize.Width = width
	m.config.Interface.WindowSize.Height = height
	m.config.Interface.Language = language
	return m.Save()
}
