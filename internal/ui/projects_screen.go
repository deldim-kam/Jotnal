package ui

import (
	"fmt"
	"time"

	"github.com/deldim-kam/Jotnal/pkg/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ProjectsScreen экран управления проектами
type ProjectsScreen struct {
	app   *App
	view  *tview.Flex
	table *tview.Table
	info  *tview.TextView
}

// NewProjectsScreen создает новый экран проектов
func NewProjectsScreen(app *App) *ProjectsScreen {
	s := &ProjectsScreen{
		app:   app,
		table: tview.NewTable().SetBorders(false).SetSelectable(true, false),
		info:  tview.NewTextView().SetDynamicColors(true),
	}

	s.table.SetBorder(true).
		SetTitle(" Список проектов ").
		SetTitleAlign(tview.AlignLeft)

	s.info.SetText("\n  [yellow]Горячие клавиши:[white]\n\n" +
		"  [green]a[white] - Добавить проект\n" +
		"  [green]e[white] - Редактировать\n" +
		"  [green]d[white] - Удалить\n" +
		"  [green]Enter[white] - Просмотр деталей\n" +
		"  [green]r[white] - Обновить список\n")
	s.info.SetBorder(true).
		SetTitle(" Информация ").
		SetTitleAlign(tview.AlignLeft)

	// Настраиваем таблицу
	s.setupTable()

	// Создаем layout
	s.view = tview.NewFlex().
		AddItem(s.table, 0, 3, true).
		AddItem(s.info, 40, 0, false)

	// Обработка клавиш
	s.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			s.addProject()
			return nil
		case 'e':
			s.editProject()
			return nil
		case 'd':
			s.deleteProject()
			return nil
		case 'r':
			s.Refresh()
			return nil
		}

		if event.Key() == tcell.KeyEnter {
			s.showDetails()
			return nil
		}

		return event
	})

	return s
}

// setupTable настраивает заголовки таблицы
func (s *ProjectsScreen) setupTable() {
	headers := []string{"ID", "Название", "Путь", "Описание", "Создан"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false).
			SetAttributes(tcell.AttrBold)
		s.table.SetCell(0, i, cell)
	}
}

// Refresh обновляет список проектов
func (s *ProjectsScreen) Refresh() {
	// Очищаем таблицу (кроме заголовка)
	s.table.Clear()
	s.setupTable()

	// Загружаем проекты из БД
	projects, err := s.loadProjects()
	if err != nil {
		s.app.ShowModal("Ошибка", "Не удалось загрузить проекты: "+err.Error(), 50, 10, nil)
		return
	}

	// Заполняем таблицу
	for i, project := range projects {
		row := i + 1

		s.table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", project.ID)).SetAlign(tview.AlignCenter))
		s.table.SetCell(row, 1, tview.NewTableCell(project.Name))
		s.table.SetCell(row, 2, tview.NewTableCell(project.Path))
		s.table.SetCell(row, 3, tview.NewTableCell(truncate(project.Description, 30)))
		s.table.SetCell(row, 4, tview.NewTableCell(project.CreatedAt.Format("2006-01-02")))
	}

	if len(projects) > 0 {
		s.table.Select(1, 0)
	}
}

// loadProjects загружает проекты из БД
func (s *ProjectsScreen) loadProjects() ([]models.Project, error) {
	rows, err := s.app.GetDB().Query("SELECT id, name, path, description, created_at, updated_at FROM projects ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Path, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// addProject добавляет новый проект
func (s *ProjectsScreen) addProject() {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Новый проект ").SetTitleAlign(tview.AlignLeft)

	var name, path, description string

	form.AddInputField("Название:", "", 40, nil, func(text string) {
		name = text
	})
	form.AddInputField("Путь:", "", 60, nil, func(text string) {
		path = text
	})
	form.AddTextArea("Описание:", "", 60, 3, 0, func(text string) {
		description = text
	})

	form.AddButton("Сохранить", func() {
		if name == "" || path == "" {
			s.app.ShowModal("Ошибка", "Название и путь обязательны для заполнения", 50, 10, nil)
			return
		}

		_, err := s.app.GetDB().Exec(
			"INSERT INTO projects (name, path, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			name, path, description, time.Now(), time.Now(),
		)
		if err != nil {
			s.app.ShowModal("Ошибка", "Не удалось создать проект: "+err.Error(), 50, 10, nil)
			return
		}

		s.app.pages.RemovePage("form")
		s.Refresh()
		s.app.ShowModal("Успех", "Проект успешно создан!", 40, 8, nil)
	})

	form.AddButton("Отмена", func() {
		s.app.pages.RemovePage("form")
	})

	s.app.pages.AddPage("form", center(form, 80, 15), true, true)
}

// editProject редактирует выбранный проект
func (s *ProjectsScreen) editProject() {
	row, _ := s.table.GetSelection()
	if row == 0 {
		return
	}

	idCell := s.table.GetCell(row, 0)
	var projectID int64
	fmt.Sscanf(idCell.Text, "%d", &projectID)

	// Загружаем данные проекта
	var p models.Project
	err := s.app.GetDB().QueryRow(
		"SELECT id, name, path, description FROM projects WHERE id = ?",
		projectID,
	).Scan(&p.ID, &p.Name, &p.Path, &p.Description)

	if err != nil {
		s.app.ShowModal("Ошибка", "Не удалось загрузить проект: "+err.Error(), 50, 10, nil)
		return
	}

	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Редактирование проекта ").SetTitleAlign(tview.AlignLeft)

	var name, path, description string
	name, path, description = p.Name, p.Path, p.Description

	form.AddInputField("Название:", name, 40, nil, func(text string) {
		name = text
	})
	form.AddInputField("Путь:", path, 60, nil, func(text string) {
		path = text
	})
	form.AddTextArea("Описание:", description, 60, 3, 0, func(text string) {
		description = text
	})

	form.AddButton("Сохранить", func() {
		_, err := s.app.GetDB().Exec(
			"UPDATE projects SET name = ?, path = ?, description = ?, updated_at = ? WHERE id = ?",
			name, path, description, time.Now(), projectID,
		)
		if err != nil {
			s.app.ShowModal("Ошибка", "Не удалось обновить проект: "+err.Error(), 50, 10, nil)
			return
		}

		s.app.pages.RemovePage("form")
		s.Refresh()
		s.app.ShowModal("Успех", "Проект успешно обновлен!", 40, 8, nil)
	})

	form.AddButton("Отмена", func() {
		s.app.pages.RemovePage("form")
	})

	s.app.pages.AddPage("form", center(form, 80, 15), true, true)
}

// deleteProject удаляет выбранный проект
func (s *ProjectsScreen) deleteProject() {
	row, _ := s.table.GetSelection()
	if row == 0 {
		return
	}

	nameCell := s.table.GetCell(row, 1)
	idCell := s.table.GetCell(row, 0)
	var projectID int64
	fmt.Sscanf(idCell.Text, "%d", &projectID)

	s.app.ShowConfirm(
		"Подтверждение удаления",
		fmt.Sprintf("Вы уверены, что хотите удалить проект '%s'?", nameCell.Text),
		func() {
			_, err := s.app.GetDB().Exec("DELETE FROM projects WHERE id = ?", projectID)
			if err != nil {
				s.app.ShowModal("Ошибка", "Не удалось удалить проект: "+err.Error(), 50, 10, nil)
				return
			}
			s.Refresh()
			s.app.ShowModal("Успех", "Проект удален!", 40, 8, nil)
		},
		nil,
	)
}

// showDetails показывает детальную информацию о проекте
func (s *ProjectsScreen) showDetails() {
	row, _ := s.table.GetSelection()
	if row == 0 {
		return
	}

	idCell := s.table.GetCell(row, 0)
	var projectID int64
	fmt.Sscanf(idCell.Text, "%d", &projectID)

	var p models.Project
	err := s.app.GetDB().QueryRow(
		"SELECT id, name, path, description, created_at, updated_at FROM projects WHERE id = ?",
		projectID,
	).Scan(&p.ID, &p.Name, &p.Path, &p.Description, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		s.app.ShowModal("Ошибка", "Не удалось загрузить проект: "+err.Error(), 50, 10, nil)
		return
	}

	details := fmt.Sprintf(
		"\n[yellow]ID:[white] %d\n\n"+
			"[yellow]Название:[white] %s\n\n"+
			"[yellow]Путь:[white] %s\n\n"+
			"[yellow]Описание:[white] %s\n\n"+
			"[yellow]Создан:[white] %s\n\n"+
			"[yellow]Обновлен:[white] %s\n",
		p.ID, p.Name, p.Path, p.Description,
		p.CreatedAt.Format("2006-01-02 15:04:05"),
		p.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(details).
		SetBorder(true).
		SetTitle(" Детали проекта ")

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		s.app.pages.RemovePage("details")
		return nil
	})

	s.app.pages.AddPage("details", center(textView, 70, 20), true, true)
}

// GetView возвращает view экрана
func (s *ProjectsScreen) GetView() tview.Primitive {
	return s.view
}

// truncate обрезает строку до указанной длины
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// center центрирует примитив на экране
func center(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
