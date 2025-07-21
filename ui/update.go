package ui

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"zan/data"
	"zan/models"
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
	case EditView:
		return m.handleEditViewKeys(msg)
	}

	return m, nil
}

func (m Model) handleListViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "enter", "l":
		if len(m.reports) > 0 {
			m.currentView = EditView
			m.isEditing = true
			m.editingIndex = m.cursor
			m.contentArea.SetValue(m.reports[m.cursor].Content)
			return m, m.contentArea.Focus()
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

func (m Model) handleEditViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		m.currentView = ListView
		m.contentArea.SetValue("")
		m.isEditing = false
		// 編集をキャンセルした場合、新規作成中のレポートは削除
		if m.editingIndex != -1 && m.reports[m.editingIndex].Content == "" {
			m.reports = append(m.reports[:m.editingIndex], m.reports[m.editingIndex+1:]...)
			if m.cursor >= len(m.reports) && m.cursor > 0 {
				m.cursor--
			}
		}
		return m, nil
	case "ctrl+s":
		content := m.contentArea.Value()

		if content != "" {
			var reportToSave models.Report
			if m.isEditing && m.editingIndex != -1 {
				// 既存レポートの更新
				reportToSave = m.reports[m.editingIndex]
				reportToSave.Content = content
				reportToSave.Date = time.Now()
			} else {
				// 新規レポートの作成
				reportToSave = models.Report{
					ID:      int(time.Now().UnixNano()), // ユニークなIDを生成
					Content: content,
					Date:    time.Now(),
				}
				m.reports = append(m.reports, reportToSave)
				m.cursor = len(m.reports) - 1
			}

			if err := data.SaveReport(reportToSave); err != nil {
				log.Printf("日報の保存に失敗しました: %v", err)
				// エラーメッセージをユーザーに表示するなどの処理を追加することも可能
			} else {
				// 保存成功後、reportsスライスを最新の状態に更新
				updatedReports, err := data.GetAllReports()
				if err != nil {
					log.Printf("日報の再読み込みに失敗しました: %v", err)
				} else {
					m.reports = updatedReports
					// 保存したレポートがリストのどこにあるか再検索し、カーソルを合わせる
					for i, r := range m.reports {
						if r.ID == reportToSave.ID {
							m.cursor = i
							break
						}
					}
				}
			}

			m.currentView = ListView
			m.contentArea.SetValue("")
			m.isEditing = false
			m.editingIndex = -1 // 編集インデックスをリセット
		}
		return m, nil
	}

	m.contentArea, cmd = m.contentArea.Update(msg)

	return m, cmd
}
