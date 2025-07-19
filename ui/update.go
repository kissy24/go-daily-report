package ui

import (
	"time"
	"zan/models"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	}

	switch m.currentView {
	case ListView:
		return m.handleListViewKeys(msg)
	case CreateView:
		return m.handleCreateViewKeys(msg)
	case DetailView:
		return m.handleDetailViewKeys(msg)
	}

	return m, nil
}

func (m Model) handleListViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "n":
		m.currentView = CreateView
		m.isEditing = false
		m.titleInput.SetValue("")
		m.contentArea.SetValue("")
		m.titleInput.Focus()
	case "enter", "l":
		if len(m.reports) > 0 {
			m.currentView = DetailView
		}
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.reports)-1 {
			m.cursor++
		}
	}
	return m, nil
}

func (m Model) handleCreateViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		m.currentView = ListView
		m.titleInput.SetValue("")
		m.contentArea.SetValue("")
		m.titleInput.Blur()
		m.isEditing = false
		return m, nil
	case "ctrl+s":
		title := m.titleInput.Value()
		content := m.contentArea.Value()

		if title != "" {
			if m.isEditing {
				m.reports[m.editingIndex].Title = title
				m.reports[m.editingIndex].Content = content
				m.reports[m.editingIndex].Date = time.Now()
			} else {
				newReport := models.Report{
					ID:      m.nextID,
					Title:   title,
					Content: content,
					Date:    time.Now(),
				}
				m.reports = append(m.reports, newReport)
				m.nextID++
				m.cursor = len(m.reports) - 1
			}

			m.currentView = ListView
			m.titleInput.SetValue("")
			m.contentArea.SetValue("")
			m.titleInput.Blur()
			m.isEditing = false
		}
		return m, nil
	case "tab":
		if m.titleInput.Focused() {
			m.titleInput.Blur()
			cmd = m.contentArea.Focus()
		} else {
			m.contentArea.Blur()
			cmd = m.titleInput.Focus()
		}
		return m, cmd
	}

	if m.titleInput.Focused() {
		m.titleInput, cmd = m.titleInput.Update(msg)
	} else {
		m.contentArea, cmd = m.contentArea.Update(msg)
	}

	return m, cmd
}

func (m Model) handleDetailViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.currentView = ListView
	case "e":
		if len(m.reports) > 0 {
			report := m.reports[m.cursor]
			m.currentView = CreateView
			m.isEditing = true
			m.editingIndex = m.cursor
			m.titleInput.SetValue(report.Title)
			m.contentArea.SetValue(report.Content)
			m.titleInput.Focus()
		}
	case "d":
		if len(m.reports) > 0 {
			m.reports = append(m.reports[:m.cursor], m.reports[m.cursor+1:]...)
			if m.cursor >= len(m.reports) && len(m.reports) > 0 {
				m.cursor = len(m.reports) - 1
			}
			if len(m.reports) == 0 {
				m.currentView = ListView
			}
		}
	}
	return m, nil
}
