package ui

import (
	"fmt"
	"time"

	"github.com/deldim-kam/Jotnal/pkg/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SnippetsScreen экран управления сниппетами
type SnippetsScreen struct {
	app     *App
	view    *tview.Flex
	list    *tview.List
	preview *tview.TextView
}

// NewSnippetsScreen создает новый экран сниппетов
func NewSnippetsScreen(app *App) *SnippetsScreen {
	s := &SnippetsScreen{
		app:     app,
		list:    tview.NewList().ShowSecondaryText(true),
		preview: tview.NewTextView().SetDynamicColors(true).SetScrollable(true),
	}

	s.list.SetBorder(true).
		SetTitle(" Список сниппетов ").
		SetTitleAlign(tview.AlignLeft)

	s.preview.SetText("\n\n  Выберите сниппет для просмотра")
	s.preview.SetBorder(true).
		SetTitle(" Предпросмотр ").
		SetTitleAlign(tview.AlignLeft)

	s.view = tview.NewFlex().
		AddItem(s.list, 0, 1, true).
		AddItem(s.preview, 0, 2, false)

	s.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			s.addSnippet()
			return nil
		case 'e':
			s.editSnippet()
			return nil
		case 'd':
			s.deleteSnippet()
			return nil
		case 'r':
			s.Refresh()
			return nil
		}
		return event
	})

	s.list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		s.showPreview(index)
	})

	return s
}

func (s *SnippetsScreen) Refresh() {
	s.list.Clear()

	snippets, err := s.loadSnippets()
	if err != nil {
		s.app.ShowModal("Ошибка", "Не удалось загрузить сниппеты: "+err.Error(), 50, 10, nil)
		return
	}

	for _, snippet := range snippets {
		title := fmt.Sprintf("[%s] %s", snippet.Language, snippet.Title)
		s.list.AddItem(title, truncate(snippet.Description, 50), 0, nil)
	}

	if len(snippets) > 0 {
		s.showPreview(0)
	}
}

func (s *SnippetsScreen) loadSnippets() ([]models.Snippet, error) {
	rows, err := s.app.GetDB().Query(
		"SELECT id, title, description, language, code, tags, created_at, updated_at FROM snippets ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []models.Snippet
	for rows.Next() {
		var snippet models.Snippet
		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Description,
			&snippet.Language, &snippet.Code, &snippet.Tags,
			&snippet.CreatedAt, &snippet.UpdatedAt)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	return snippets, nil
}

func (s *SnippetsScreen) showPreview(index int) {
	snippets, _ := s.loadSnippets()
	if index < 0 || index >= len(snippets) {
		return
	}

	snippet := snippets[index]
	preview := fmt.Sprintf(
		"\n[yellow]Название:[white] %s\n\n"+
			"[yellow]Язык:[white] %s\n\n"+
			"[yellow]Описание:[white] %s\n\n"+
			"[yellow]Теги:[white] %s\n\n"+
			"[yellow]Создан:[white] %s\n\n"+
			"[yellow]Код:[white]\n\n%s",
		snippet.Title, snippet.Language, snippet.Description,
		snippet.Tags, snippet.CreatedAt.Format("2006-01-02 15:04:05"),
		snippet.Code,
	)

	s.preview.SetText(preview)
}

func (s *SnippetsScreen) addSnippet() {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Новый сниппет ").SetTitleAlign(tview.AlignLeft)

	var title, description, language, code, tags string

	form.AddInputField("Название:*", "", 50, nil, func(text string) {
		title = text
	})
	form.AddInputField("Язык:*", "go", 20, nil, func(text string) {
		language = text
	})
	form.AddInputField("Теги:", "", 50, nil, func(text string) {
		tags = text
	})
	form.AddTextArea("Описание:", "", 60, 2, 0, func(text string) {
		description = text
	})
	form.AddTextArea("Код:*", "", 60, 8, 0, func(text string) {
		code = text
	})

	form.AddButton("Сохранить", func() {
		if title == "" || language == "" || code == "" {
			s.app.ShowModal("Ошибка", "Название, язык и код обязательны", 50, 10, nil)
			return
		}

		_, err := s.app.GetDB().Exec(
			`INSERT INTO snippets (title, description, language, code, tags, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			title, description, language, code, tags, time.Now(), time.Now(),
		)
		if err != nil {
			s.app.ShowModal("Ошибка", "Не удалось создать сниппет: "+err.Error(), 50, 10, nil)
			return
		}

		s.app.pages.RemovePage("form")
		s.Refresh()
		s.app.ShowModal("Успех", "Сниппет успешно создан!", 40, 8, nil)
	})

	form.AddButton("Отмена", func() {
		s.app.pages.RemovePage("form")
	})

	s.app.pages.AddPage("form", center(form, 80, 25), true, true)
}

func (s *SnippetsScreen) editSnippet() {
	index := s.list.GetCurrentItem()
	snippets, _ := s.loadSnippets()
	if index < 0 || index >= len(snippets) {
		return
	}

	snippet := snippets[index]

	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Редактирование сниппета ").SetTitleAlign(tview.AlignLeft)

	var title, description, language, code, tags string
	title, description, language = snippet.Title, snippet.Description, snippet.Language
	code, tags = snippet.Code, snippet.Tags

	form.AddInputField("Название:*", title, 50, nil, func(text string) {
		title = text
	})
	form.AddInputField("Язык:*", language, 20, nil, func(text string) {
		language = text
	})
	form.AddInputField("Теги:", tags, 50, nil, func(text string) {
		tags = text
	})
	form.AddTextArea("Описание:", description, 60, 2, 0, func(text string) {
		description = text
	})
	form.AddTextArea("Код:*", code, 60, 8, 0, func(text string) {
		code = text
	})

	form.AddButton("Сохранить", func() {
		_, err := s.app.GetDB().Exec(
			`UPDATE snippets SET title = ?, description = ?, language = ?, code = ?, tags = ?, updated_at = ?
			 WHERE id = ?`,
			title, description, language, code, tags, time.Now(), snippet.ID,
		)
		if err != nil {
			s.app.ShowModal("Ошибка", "Не удалось обновить сниппет: "+err.Error(), 50, 10, nil)
			return
		}

		s.app.pages.RemovePage("form")
		s.Refresh()
		s.app.ShowModal("Успех", "Сниппет успешно обновлен!", 40, 8, nil)
	})

	form.AddButton("Отмена", func() {
		s.app.pages.RemovePage("form")
	})

	s.app.pages.AddPage("form", center(form, 80, 25), true, true)
}

func (s *SnippetsScreen) deleteSnippet() {
	index := s.list.GetCurrentItem()
	snippets, _ := s.loadSnippets()
	if index < 0 || index >= len(snippets) {
		return
	}

	snippet := snippets[index]

	s.app.ShowConfirm(
		"Подтверждение удаления",
		fmt.Sprintf("Вы уверены, что хотите удалить сниппет '%s'?", snippet.Title),
		func() {
			_, err := s.app.GetDB().Exec("DELETE FROM snippets WHERE id = ?", snippet.ID)
			if err != nil {
				s.app.ShowModal("Ошибка", "Не удалось удалить сниппет: "+err.Error(), 50, 10, nil)
				return
			}
			s.Refresh()
			s.app.ShowModal("Успех", "Сниппет удален!", 40, 8, nil)
		},
		nil,
	)
}

func (s *SnippetsScreen) GetView() tview.Primitive {
	return s.view
}
