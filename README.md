# Jotnal - Журнал смен СРПК

WPF приложение для управления сменами и сотрудниками.

## Описание

Система управления сотрудниками с авторизацией через Windows, иерархической структурой подразделений и системой запросов на регистрацию.

## Сборка исполняемых файлов

### Использование скрипта сборки

Проект включает удобный скрипт `build.sh` для сборки под разные платформы:

```bash
# Собрать для всех платформ (Linux и Windows)
./build.sh all

# Собрать только для Windows (exe файл)
./build.sh windows

# Собрать только для Linux
./build.sh linux

# Очистить директорию сборки
./build.sh clean
```

Собранные файлы будут находиться в директории `./build/`:
- `build/jotnal` - исполняемый файл для Linux
- `build/jotnal.exe` - исполняемый файл для Windows

### Ручная сборка

#### Сборка для Linux:
```bash
go build -o jotnal ./cmd/ide
```

#### Сборка для Windows (кросс-компиляция):
```bash
# Требуется установленный MinGW-w64
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o jotnal.exe ./cmd/ide
```

**Требования для кросс-компиляции Windows:**

Ubuntu/Debian:
```bash
sudo apt-get install mingw-w64
```

Fedora/RHEL:
```bash
sudo dnf install mingw64-gcc
```

### Примечания по сборке

- Exe файл для Windows можно собрать из Linux благодаря кросс-компиляции
- Размер exe файла: ~14 MB
- Поддерживается Windows 7 и выше (x64)
- SQLCipher компилируется статически, дополнительные DLL не требуются

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

#### Positions (Должности)
- Id - уникальный идентификатор
- Name - название должности
- Description - описание
- IsActive - активна ли должность
- CreatedAt, UpdatedAt - даты создания и обновления

#### Departments (Структурные подразделения)
- Id - уникальный идентификатор
- Name - название подразделения
- Description - описание
- ParentDepartmentId - ссылка на родительское подразделение (иерархия)
- IsActive - активно ли подразделение
- CreatedAt, UpdatedAt - даты создания и обновления

#### Employees (Сотрудники/Пользователи)
- Id - уникальный идентификатор
- LastName, FirstName, MiddleName - ФИО
- Phone - телефон
- BirthDate - дата рождения
- PersonnelNumber - табельный номер (уникальный)
- HireDate - дата приема на работу
- IsCurrentlyEmployed - работает ли сейчас
- TerminationDate - дата увольнения (если уволен)
- WindowsUsername - имя пользователя Windows (уникальное)
- IsActive - активен ли пользователь
- Role - роль (User, Administrator, Developer)
- PositionId - ссылка на должность
- DepartmentId - ссылка на подразделение
- LastLoginAt - время последнего входа

#### RegistrationRequests (Запросы на регистрацию)
- Id - уникальный идентификатор
- LastName, FirstName, MiddleName - ФИО
- Phone - телефон
- BirthDate - дата рождения
- PersonnelNumber - табельный номер
- HireDate - дата приема на работу
- WindowsUsername - имя пользователя Windows
- RequestedPosition - запрашиваемая должность
- RequestedDepartment - запрашиваемое подразделение
- Comments - комментарии
- Status - статус (Pending, Approved, Rejected)
- RequestedAt - дата запроса
- ApprovedByEmployeeId - кто обработал запрос
- ProcessedAt - дата обработки
- RejectionReason - причина отклонения

## Особенности

### Авторизация через Windows
- Автоматическая авторизация по имени пользователя Windows
- Проверка наличия пользователя в базе данных
- Возможность смены пользователя

### Специальный пользователь "Разработчик"
- Windows username: DEVELOPER
- Максимальные права в системе
- Минимальные следы (не обновляется LastLoginAt)
- Создается автоматически при инициализации базы данных

### Система ролей
- User - обычный пользователь
- Administrator - администратор (доступ к административным функциям)
- Developer - разработчик (максимальные права)

### Запросы на регистрацию
- Неавторизованные пользователи могут подать запрос на регистрацию
- Администраторы могут одобрить или отклонить запросы
- При одобрении автоматически создается учетная запись сотрудника

## Технологии

- .NET 8.0
- WPF (Windows Presentation Foundation)
- Entity Framework Core 8.0
- SQLite

## Структура проекта

```
Jotnal/
├── Models/              # Модели данных
│   ├── Position.cs
│   ├── Department.cs
│   ├── Employee.cs
│   └── RegistrationRequest.cs
├── Data/                # Контекст базы данных
│   └── JotnalDbContext.cs
├── Services/            # Бизнес-логика
│   ├── AuthenticationService.cs
│   └── RegistrationService.cs
├── Views/               # Представления WPF
├── ViewModels/          # Модели представлений
├── App.xaml             # Точка входа приложения
└── MainWindow.xaml      # Главное окно
```

## Сборка и запуск

1. Установите .NET 8.0 SDK
2. Откройте решение в Visual Studio 2022 или Rider
3. Восстановите NuGet пакеты
4. Соберите проект
5. Запустите приложение

База данных SQLite будет создана автоматически при первом запуске.

## Начальные данные

При первом запуске автоматически создаются:
- Должность "Администратор системы"
- Должность "Разработчик"
- Подразделение "Головная организация"
- Пользователь "Системный Разработчик" (DEVELOPER)
