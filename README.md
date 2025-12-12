# Jotnal

**Jotnal** - это IDE приложение на Go с зашифрованной базой данных SQLite для управления проектами и файлами.

## Особенности

- ✅ **Зашифрованная база данных SQLCipher** - все данные защищены паролем
- ✅ **Автоматическое создание БД** при первом запуске
- ✅ **Система миграций** для версионирования схемы БД
- ✅ **Автоматическое обновление схемы** - создание недостающих таблиц
- ✅ **Гибкая конфигурация** - изменение пути к БД и настроек интерфейса
- ✅ **Автоматическое подключение** при наличии пути к БД
- ✅ **Хранение конфигурации** в JSON файле

## Структура проекта

```
Jotnal/
├── cmd/
│   └── ide/
│       └── main.go              # Главный файл приложения
├── internal/
│   ├── config/
│   │   └── config.go            # Менеджер конфигурации
│   └── database/
│       ├── database.go          # Менеджер БД с шифрованием
│       └── migrations.go        # Система миграций
├── pkg/
│   └── models/
│       └── models.go            # Модели данных
├── go.mod
└── README.md
```

## Требования

### Системные зависимости

Для работы с зашифрованной SQLite необходимо установить SQLCipher:

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install -y gcc libsqlcipher-dev
```

**macOS:**
```bash
brew install sqlcipher
```

**Fedora/RHEL:**
```bash
sudo dnf install sqlcipher-devel gcc
```

### Go зависимости

```bash
go get github.com/mutecomm/go-sqlcipher/v4
go get golang.org/x/term
```

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/deldim-kam/Jotnal.git
cd Jotnal
```

2. Установите зависимости:
```bash
go mod download
```

3. Соберите приложение:
```bash
go build -o jotnal ./cmd/ide
```

4. Запустите:
```bash
./jotnal
```

При первом запуске вам будет предложено установить пароль для базы данных.

## Конфигурация

Конфигурация хранится в файле `~/.jotnal/config.json`:

```json
{
  "database": {
    "path": "/home/user/.jotnal/jotnal.db",
    "password": "ваш_пароль"
  },
  "interface": {
    "theme": "dark",
    "font_size": 14,
    "window_size": {
      "width": 1280,
      "height": 720
    },
    "language": "ru"
  }
}
```

## База данных

### Таблицы

База данных включает следующие таблицы:

1. **projects** - хранение информации о проектах
2. **files** - файлы проектов
3. **project_settings** - настройки каждого проекта
4. **file_history** - история изменений файлов
5. **snippets** - сниппеты кода
6. **bookmarks** - закладки в файлах
7. **employees** - сотрудники с иерархической структурой (начальник-подчиненный)
8. **schema_version** - версионирование схемы БД

### Миграции

Система миграций автоматически:
- Создает новую БД при первом запуске
- Проверяет версию схемы
- Применяет недостающие миграции
- Обновляет существующую БД до актуальной версии

Текущая версия схемы: **5**

#### Иерархическая структура сотрудников

Таблица **employees** поддерживает многоуровневую иерархию через поле `manager_id`:
- Каждый сотрудник может иметь руководителя (manager_id ссылается на id другого сотрудника)
- Главный руководитель имеет manager_id = NULL
- Поддерживается неограниченная вложенность: начальник → подчиненный → подчиненный подчиненного и т.д.
- При удалении руководителя, manager_id его подчиненных автоматически устанавливается в NULL

## Безопасность

### Шифрование

- База данных использует **SQLCipher** для шифрования
- Алгоритм: 256-bit AES
- Размер страницы: 4096 байт
- БД доступна только с правильным паролем
- Пароль хранится в конфигурационном файле с правами доступа 0600

### Рекомендации

1. Используйте сложный пароль (минимум 16 символов)
2. Регулярно делайте резервные копии БД
3. Не передавайте конфигурационный файл третьим лицам
4. Храните пароль в надежном месте

## Возможности через меню

1. **Показать информацию о БД** - просмотр пути, версии и списка таблиц
2. **Изменить путь к БД** - переключение на другую базу данных
3. **Изменить пароль БД** - смена пароля шифрования
4. **Показать настройки интерфейса** - просмотр текущих настроек
5. **Изменить настройки интерфейса** - настройка темы, шрифта, размера окна и языка

## Разработка

### Добавление новых миграций

Редактируйте файл `internal/database/migrations.go`:

```go
{
    Version:     5,
    Description: "Описание новой миграции",
    SQL: `
        CREATE TABLE new_table (...);
    `,
},
```

### Добавление новых моделей

Добавьте структуры в `pkg/models/models.go`:

```go
type NewModel struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

## Примеры использования

### Программное создание проекта

```go
db := dbManager.GetDB()
result, err := db.Exec(
    "INSERT INTO projects (name, path, description) VALUES (?, ?, ?)",
    "Мой проект",
    "/path/to/project",
    "Описание проекта",
)
```

### Получение списка проектов

```go
rows, err := db.Query("SELECT id, name, path FROM projects")
defer rows.Close()

for rows.Next() {
    var id int64
    var name, path string
    rows.Scan(&id, &name, &path)
    fmt.Printf("Проект: %s (%s)\n", name, path)
}
```

### Работа с иерархией сотрудников

#### Создание сотрудника

```go
// Создание главного руководителя (без начальника)
result, err := db.Exec(
    "INSERT INTO employees (first_name, last_name, email, position, manager_id) VALUES (?, ?, ?, ?, ?)",
    "Иван", "Иванов", "ivanov@example.com", "Генеральный директор", nil,
)

// Создание подчиненного
result, err := db.Exec(
    "INSERT INTO employees (first_name, last_name, email, position, manager_id) VALUES (?, ?, ?, ?, ?)",
    "Петр", "Петров", "petrov@example.com", "Менеджер отдела", 1, // 1 - ID руководителя
)
```

#### Получение всех подчиненных руководителя

```go
rows, err := db.Query(`
    SELECT id, first_name, last_name, position
    FROM employees
    WHERE manager_id = ?
`, managerID)
```

#### Получение полной иерархии (рекурсивный запрос)

```go
rows, err := db.Query(`
    WITH RECURSIVE employee_hierarchy AS (
        -- Начальная точка: выбранный сотрудник
        SELECT id, first_name, last_name, position, manager_id, 0 as level
        FROM employees
        WHERE id = ?

        UNION ALL

        -- Рекурсивная часть: все подчиненные
        SELECT e.id, e.first_name, e.last_name, e.position, e.manager_id, eh.level + 1
        FROM employees e
        INNER JOIN employee_hierarchy eh ON e.manager_id = eh.id
    )
    SELECT * FROM employee_hierarchy
    ORDER BY level, last_name
`, rootEmployeeID)
```

## Лицензия

MIT License

## Автор

deldim-kam

## Содействие

Pull requests приветствуются! Для серьезных изменений сначала откройте issue для обсуждения.
