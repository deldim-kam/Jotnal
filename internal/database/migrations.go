package database

// Migration представляет миграцию базы данных
type Migration struct {
	Version     int
	Description string
	SQL         string
}

// GetMigrations возвращает список всех миграций
func GetMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Description: "Создание базовых таблиц для IDE",
			SQL: `
				-- Таблица проектов
				CREATE TABLE IF NOT EXISTS projects (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL,
					path TEXT NOT NULL UNIQUE,
					description TEXT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Таблица файлов
				CREATE TABLE IF NOT EXISTS files (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					project_id INTEGER NOT NULL,
					path TEXT NOT NULL,
					name TEXT NOT NULL,
					content TEXT,
					size INTEGER DEFAULT 0,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
					UNIQUE(project_id, path)
				);

				-- Таблица настроек проекта
				CREATE TABLE IF NOT EXISTS project_settings (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					project_id INTEGER NOT NULL UNIQUE,
					language TEXT DEFAULT 'go',
					build_command TEXT,
					run_command TEXT,
					test_command TEXT,
					linter_command TEXT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
				);

				-- Индексы
				CREATE INDEX IF NOT EXISTS idx_files_project_id ON files(project_id);
				CREATE INDEX IF NOT EXISTS idx_files_name ON files(name);
				CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
			`,
		},
		{
			Version:     2,
			Description: "Добавление таблицы для истории изменений",
			SQL: `
				-- Таблица истории изменений файлов
				CREATE TABLE IF NOT EXISTS file_history (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					file_id INTEGER NOT NULL,
					content TEXT,
					change_description TEXT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
				);

				-- Индекс
				CREATE INDEX IF NOT EXISTS idx_file_history_file_id ON file_history(file_id);
			`,
		},
		{
			Version:     3,
			Description: "Добавление таблицы сниппетов кода",
			SQL: `
				-- Таблица сниппетов кода
				CREATE TABLE IF NOT EXISTS snippets (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					title TEXT NOT NULL,
					description TEXT,
					language TEXT NOT NULL,
					code TEXT NOT NULL,
					tags TEXT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Индексы
				CREATE INDEX IF NOT EXISTS idx_snippets_language ON snippets(language);
				CREATE INDEX IF NOT EXISTS idx_snippets_title ON snippets(title);
			`,
		},
		{
			Version:     4,
			Description: "Добавление таблицы закладок",
			SQL: `
				-- Таблица закладок
				CREATE TABLE IF NOT EXISTS bookmarks (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					file_id INTEGER NOT NULL,
					line_number INTEGER NOT NULL,
					description TEXT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
				);

				-- Индекс
				CREATE INDEX IF NOT EXISTS idx_bookmarks_file_id ON bookmarks(file_id);
			`,
		},
	}
}
