package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

// Manager управляет подключением к базе данных
type Manager struct {
	db       *sql.DB
	dbPath   string
	password string
	version  int
}

// NewManager создает новый менеджер базы данных
func NewManager(dbPath, password string) (*Manager, error) {
	m := &Manager{
		dbPath:   dbPath,
		password: password,
	}

	return m, nil
}

// Connect подключается к базе данных
func (m *Manager) Connect() error {
	// Создаем директорию для БД если не существует
	dbDir := filepath.Dir(m.dbPath)
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return fmt.Errorf("не удалось создать директорию для БД: %w", err)
	}

	// Проверяем существует ли файл БД
	isNewDB := !fileExists(m.dbPath)

	// Формируем DSN с параметрами шифрования
	dsn := fmt.Sprintf("file:%s?_pragma_key=%s&_pragma_cipher_page_size=4096", m.dbPath, m.password)

	// Открываем подключение
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("не удалось открыть БД: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	m.db = db

	// Если это новая БД, инициализируем её
	if isNewDB {
		if err := m.initialize(); err != nil {
			m.db.Close()
			return fmt.Errorf("не удалось инициализировать БД: %w", err)
		}
	} else {
		// Проверяем версию БД
		if err := m.checkVersion(); err != nil {
			m.db.Close()
			return fmt.Errorf("не удалось проверить версию БД: %w", err)
		}
	}

	return nil
}

// initialize инициализирует новую базу данных
func (m *Manager) initialize() error {
	// Создаем таблицу версий
	if err := m.createVersionTable(); err != nil {
		return err
	}

	// Применяем все миграции
	if err := m.runMigrations(); err != nil {
		return err
	}

	return nil
}

// createVersionTable создает таблицу для отслеживания версии БД
func (m *Manager) createVersionTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_version (
		version INTEGER PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := m.db.Exec(query)
	return err
}

// checkVersion проверяет версию базы данных
func (m *Manager) checkVersion() error {
	var version int
	err := m.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
	if err != nil {
		return err
	}

	m.version = version

	// Применяем недостающие миграции
	return m.runMigrations()
}

// runMigrations применяет все необходимые миграции
func (m *Manager) runMigrations() error {
	migrations := GetMigrations()

	for _, migration := range migrations {
		if migration.Version > m.version {
			if err := m.applyMigration(migration); err != nil {
				return fmt.Errorf("не удалось применить миграцию %d: %w", migration.Version, err)
			}
		}
	}

	return nil
}

// applyMigration применяет одну миграцию
func (m *Manager) applyMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Выполняем SQL миграции
	if _, err := tx.Exec(migration.SQL); err != nil {
		return err
	}

	// Обновляем версию
	if _, err := tx.Exec("INSERT INTO schema_version (version) VALUES (?)", migration.Version); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	m.version = migration.Version
	return nil
}

// GetDB возвращает экземпляр базы данных
func (m *Manager) GetDB() *sql.DB {
	return m.db
}

// GetVersion возвращает текущую версию БД
func (m *Manager) GetVersion() int {
	return m.version
}

// Close закрывает подключение к базе данных
func (m *Manager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// ChangePassword изменяет пароль базы данных
func (m *Manager) ChangePassword(newPassword string) error {
	if m.db == nil {
		return fmt.Errorf("БД не подключена")
	}

	// Используем PRAGMA rekey для смены пароля
	_, err := m.db.Exec(fmt.Sprintf("PRAGMA rekey = '%s'", newPassword))
	if err != nil {
		return fmt.Errorf("не удалось изменить пароль БД: %w", err)
	}

	m.password = newPassword
	return nil
}

// fileExists проверяет существование файла
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
