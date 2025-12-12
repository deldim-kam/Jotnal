package models

import "time"

// Project представляет проект в IDE
type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// File представляет файл в проекте
type File struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProjectSettings представляет настройки проекта
type ProjectSettings struct {
	ID            int64     `json:"id"`
	ProjectID     int64     `json:"project_id"`
	Language      string    `json:"language"`
	BuildCommand  string    `json:"build_command"`
	RunCommand    string    `json:"run_command"`
	TestCommand   string    `json:"test_command"`
	LinterCommand string    `json:"linter_command"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// FileHistory представляет историю изменений файла
type FileHistory struct {
	ID                int64     `json:"id"`
	FileID            int64     `json:"file_id"`
	Content           string    `json:"content"`
	ChangeDescription string    `json:"change_description"`
	CreatedAt         time.Time `json:"created_at"`
}

// Snippet представляет сниппет кода
type Snippet struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Code        string    `json:"code"`
	Tags        string    `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Bookmark представляет закладку в файле
type Bookmark struct {
	ID          int64     `json:"id"`
	FileID      int64     `json:"file_id"`
	LineNumber  int       `json:"line_number"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Employee представляет сотрудника с иерархической структурой
type Employee struct {
	ID         int64     `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
	Email      string    `json:"email"`
	Position   string    `json:"position"`
	Department string    `json:"department"`
	ManagerID  *int64    `json:"manager_id"` // NULL для главного руководителя
	Phone      string    `json:"phone"`
	HireDate   time.Time `json:"hire_date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
