package ui

import (
	"fmt"
	"time"

	"github.com/deldim-kam/Jotnal/pkg/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// EmployeesScreen экран управления сотрудниками
type EmployeesScreen struct {
	app   *App
	view  *tview.Flex
	table *tview.Table
	info  *tview.TextView
}

// NewEmployeesScreen создает новый экран сотрудников
func NewEmployeesScreen(app *App) *EmployeesScreen {
	s := &EmployeesScreen{
		app:   app,
		table: tview.NewTable().SetBorders(false).SetSelectable(true, false),
		info:  tview.NewTextView().SetDynamicColors(true),
	}

	s.table.SetBorder(true).
		SetTitle(" Список сотрудников ").
		SetTitleAlign(tview.AlignLeft)

	s.info.SetText("\n  [yellow]Горячие клавиши:[white]\n\n" +
		"  [green]a[white] - Добавить сотрудника\n" +
		"  [green]e[white] - Редактировать\n" +
		"  [green]d[white] - Удалить\n" +
		"  [green]Enter[white] - Просмотр деталей\n" +
		"  [green]r[white] - Обновить список\n")
	s.info.SetBorder(true).
		SetTitle(" Информация ").
		SetTitleAlign(tview.AlignLeft)

	s.setupTable()

	s.view = tview.NewFlex().
		AddItem(s.table, 0, 3, true).
		AddItem(s.info, 40, 0, false)

	s.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			s.addEmployee()
			return nil
		case 'e':
			s.editEmployee()
			return nil
		case 'd':
			s.deleteEmployee()
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

func (s *EmployeesScreen) setupTable() {
	headers := []string{"ID", "Фамилия", "Имя", "Email", "Должность", "Отдел", "Дата найма"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false).
			SetAttributes(tcell.AttrBold)
		s.table.SetCell(0, i, cell)
	}
}

func (s *EmployeesScreen) Refresh() {
	s.table.Clear()
	s.setupTable()

	employees, err := s.loadEmployees()
	if err != nil {
		s.app.ShowModal("Ошибка", "Не удалось загрузить сотрудников: "+err.Error(), 50, 10, nil)
		return
	}

	for i, emp := range employees {
		row := i + 1
		s.table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", emp.ID)).SetAlign(tview.AlignCenter))
		s.table.SetCell(row, 1, tview.NewTableCell(emp.LastName))
		s.table.SetCell(row, 2, tview.NewTableCell(emp.FirstName))
		s.table.SetCell(row, 3, tview.NewTableCell(emp.Email))
		s.table.SetCell(row, 4, tview.NewTableCell(emp.Position))
		s.table.SetCell(row, 5, tview.NewTableCell(emp.Department))
		s.table.SetCell(row, 6, tview.NewTableCell(emp.HireDate.Format("2006-01-02")))
	}

	if len(employees) > 0 {
		s.table.Select(1, 0)
	}
}

func (s *EmployeesScreen) loadEmployees() ([]models.Employee, error) {
	query := `SELECT id, first_name, last_name, middle_name, email, position,
			  department, manager_id, phone, hire_date, created_at, updated_at
			  FROM employees ORDER BY last_name, first_name`

	rows, err := s.app.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		err := rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.MiddleName, &e.Email,
			&e.Position, &e.Department, &e.ManagerID, &e.Phone, &e.HireDate,
			&e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}

	return employees, nil
}

func (s *EmployeesScreen) addEmployee() {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Новый сотрудник ").SetTitleAlign(tview.AlignLeft)

	var lastName, firstName, middleName, email, position, department, phone string

	form.AddInputField("Фамилия:*", "", 30, nil, func(text string) {
		lastName = text
	})
	form.AddInputField("Имя:*", "", 30, nil, func(text string) {
		firstName = text
	})
	form.AddInputField("Отчество:", "", 30, nil, func(text string) {
		middleName = text
	})
	form.AddInputField("Email:", "", 40, nil, func(text string) {
		email = text
	})
	form.AddInputField("Должность:*", "", 40, nil, func(text string) {
		position = text
	})
	form.AddInputField("Отдел:", "", 40, nil, func(text string) {
		department = text
	})
	form.AddInputField("Телефон:", "", 20, nil, func(text string) {
		phone = text
	})

	form.AddButton("Сохранить", func() {
		if lastName == "" || firstName == "" || position == "" {
			s.app.ShowModal("Ошибка", "Фамилия, имя и должность обязательны", 50, 10, nil)
			return
		}

		_, err := s.app.GetDB().Exec(
			`INSERT INTO employees (first_name, last_name, middle_name, email, position,
			 department, phone, hire_date, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			firstName, lastName, middleName, email, position, department, phone,
			time.Now(), time.Now(), time.Now(),
		)
		if err != nil {
			s.app.ShowModal("Ошибка", "Не удалось создать сотрудника: "+err.Error(), 50, 10, nil)
			return
		}

		s.app.pages.RemovePage("form")
		s.Refresh()
		s.app.ShowModal("Успех", "Сотрудник успешно добавлен!", 40, 8, nil)
	})

	form.AddButton("Отмена", func() {
		s.app.pages.RemovePage("form")
	})

	s.app.pages.AddPage("form", center(form, 70, 20), true, true)
}

func (s *EmployeesScreen) editEmployee() {
	row, _ := s.table.GetSelection()
	if row == 0 {
		return
	}

	idCell := s.table.GetCell(row, 0)
	var empID int64
	fmt.Sscanf(idCell.Text, "%d", &empID)

	var e models.Employee
	err := s.app.GetDB().QueryRow(
		`SELECT id, first_name, last_name, middle_name, email, position,
		 department, phone FROM employees WHERE id = ?`,
		empID,
	).Scan(&e.ID, &e.FirstName, &e.LastName, &e.MiddleName, &e.Email,
		&e.Position, &e.Department, &e.Phone)

	if err != nil {
		s.app.ShowModal("Ошибка", "Не удалось загрузить сотрудника: "+err.Error(), 50, 10, nil)
		return
	}

	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Редактирование сотрудника ").SetTitleAlign(tview.AlignLeft)

	var lastName, firstName, middleName, email, position, department, phone string
	lastName, firstName, middleName = e.LastName, e.FirstName, e.MiddleName
	email, position, department, phone = e.Email, e.Position, e.Department, e.Phone

	form.AddInputField("Фамилия:*", lastName, 30, nil, func(text string) {
		lastName = text
	})
	form.AddInputField("Имя:*", firstName, 30, nil, func(text string) {
		firstName = text
	})
	form.AddInputField("Отчество:", middleName, 30, nil, func(text string) {
		middleName = text
	})
	form.AddInputField("Email:", email, 40, nil, func(text string) {
		email = text
	})
	form.AddInputField("Должность:*", position, 40, nil, func(text string) {
		position = text
	})
	form.AddInputField("Отдел:", department, 40, nil, func(text string) {
		department = text
	})
	form.AddInputField("Телефон:", phone, 20, nil, func(text string) {
		phone = text
	})

	form.AddButton("Сохранить", func() {
		_, err := s.app.GetDB().Exec(
			`UPDATE employees SET first_name = ?, last_name = ?, middle_name = ?,
			 email = ?, position = ?, department = ?, phone = ?, updated_at = ?
			 WHERE id = ?`,
			firstName, lastName, middleName, email, position, department, phone,
			time.Now(), empID,
		)
		if err != nil {
			s.app.ShowModal("Ошибка", "Не удалось обновить сотрудника: "+err.Error(), 50, 10, nil)
			return
		}

		s.app.pages.RemovePage("form")
		s.Refresh()
		s.app.ShowModal("Успех", "Сотрудник успешно обновлен!", 40, 8, nil)
	})

	form.AddButton("Отмена", func() {
		s.app.pages.RemovePage("form")
	})

	s.app.pages.AddPage("form", center(form, 70, 22), true, true)
}

func (s *EmployeesScreen) deleteEmployee() {
	row, _ := s.table.GetSelection()
	if row == 0 {
		return
	}

	nameCell := s.table.GetCell(row, 1)
	idCell := s.table.GetCell(row, 0)
	var empID int64
	fmt.Sscanf(idCell.Text, "%d", &empID)

	s.app.ShowConfirm(
		"Подтверждение удаления",
		fmt.Sprintf("Вы уверены, что хотите удалить сотрудника '%s'?", nameCell.Text),
		func() {
			_, err := s.app.GetDB().Exec("DELETE FROM employees WHERE id = ?", empID)
			if err != nil {
				s.app.ShowModal("Ошибка", "Не удалось удалить сотрудника: "+err.Error(), 50, 10, nil)
				return
			}
			s.Refresh()
			s.app.ShowModal("Успех", "Сотрудник удален!", 40, 8, nil)
		},
		nil,
	)
}

func (s *EmployeesScreen) showDetails() {
	row, _ := s.table.GetSelection()
	if row == 0 {
		return
	}

	idCell := s.table.GetCell(row, 0)
	var empID int64
	fmt.Sscanf(idCell.Text, "%d", &empID)

	var e models.Employee
	err := s.app.GetDB().QueryRow(
		`SELECT id, first_name, last_name, middle_name, email, position,
		 department, manager_id, phone, hire_date, created_at, updated_at
		 FROM employees WHERE id = ?`,
		empID,
	).Scan(&e.ID, &e.FirstName, &e.LastName, &e.MiddleName, &e.Email,
		&e.Position, &e.Department, &e.ManagerID, &e.Phone, &e.HireDate,
		&e.CreatedAt, &e.UpdatedAt)

	if err != nil {
		s.app.ShowModal("Ошибка", "Не удалось загрузить сотрудника: "+err.Error(), 50, 10, nil)
		return
	}

	managerName := "Нет"
	if e.ManagerID != nil {
		var mFirstName, mLastName string
		s.app.GetDB().QueryRow(
			"SELECT first_name, last_name FROM employees WHERE id = ?",
			*e.ManagerID,
		).Scan(&mFirstName, &mLastName)
		managerName = fmt.Sprintf("%s %s", mLastName, mFirstName)
	}

	details := fmt.Sprintf(
		"\n[yellow]ID:[white] %d\n\n"+
			"[yellow]ФИО:[white] %s %s %s\n\n"+
			"[yellow]Email:[white] %s\n\n"+
			"[yellow]Должность:[white] %s\n\n"+
			"[yellow]Отдел:[white] %s\n\n"+
			"[yellow]Телефон:[white] %s\n\n"+
			"[yellow]Руководитель:[white] %s\n\n"+
			"[yellow]Дата найма:[white] %s\n\n"+
			"[yellow]Создан:[white] %s\n",
		e.ID, e.LastName, e.FirstName, e.MiddleName, e.Email,
		e.Position, e.Department, e.Phone, managerName,
		e.HireDate.Format("2006-01-02"),
		e.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(details).
		SetBorder(true).
		SetTitle(" Детали сотрудника ")

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		s.app.pages.RemovePage("details")
		return nil
	})

	s.app.pages.AddPage("details", center(textView, 70, 22), true, true)
}

func (s *EmployeesScreen) GetView() tview.Primitive {
	return s.view
}
