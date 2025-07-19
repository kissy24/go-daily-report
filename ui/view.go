package ui

import (
	"fmt"
	"strings"
	"zan/styles"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	switch m.currentView {
	case ListView:
		return m.renderListView()
	case EditView:
		return m.renderEditView()
	default:
		return styles.ErrorStyle.Render("不明な画面です")
	}
}

func (m Model) renderListView() string {
	var sections []string

	// ヘッダー
	header := styles.HeaderStyle.Width(m.width - 4).Render("📝 Daily Report Manager")
	sections = append(sections, header)

	if len(m.reports) == 0 {
		emptyMsg := styles.ContentBoxStyle.Width(m.width - 4).Render(
			styles.LabelStyle.Render("📋 日報がありません") + "\n\n" +
				"新しい日報を作成してください。",
		)
		sections = append(sections, emptyMsg)
	} else {
		// 日報リスト
		var items []string
		for i, report := range m.reports {
			dateStr := styles.DateStyle.Render(report.Date.Format("01/02 15:04"))

			if m.cursor == i {
				item := styles.SelectedItemStyle.Render(fmt.Sprintf("▶ %s", dateStr))
				items = append(items, item)
			} else {
				item := styles.NormalItemStyle.Render(fmt.Sprintf("  %s", dateStr))
				items = append(items, item)
			}
		}

		listBox := styles.ContentBoxStyle.Width(m.width - 4).Render(strings.Join(items, "\n"))
		sections = append(sections, listBox)
	}

	// ヘルプ
	var helpButtons []string
	helpButtons = append(helpButtons, styles.PrimaryButtonStyle.Render("Enter")+" 詳細表示")
	helpButtons = append(helpButtons, styles.PrimaryButtonStyle.Render("N")+" 本日の日報")
	helpButtons = append(helpButtons, styles.SecondaryButtonStyle.Render("Q")+" 終了")
	helpButtons = append(helpButtons, styles.SecondaryButtonStyle.Render("↑↓")+" 移動")

	help := styles.HelpStyle.Width(m.width - 4).Render(strings.Join(helpButtons, "  "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderEditView() string {
	var sections []string

	// ヘッダー
	var headerText string
	if m.isEditing {
		headerText = "✏️ Edit Daily Report"
	} else {
		headerText = "➕ Create New Daily Report"
	}
	header := styles.HeaderStyle.Width(m.width - 4).Render(headerText)
	sections = append(sections, header)

	// 内容入力
	contentLabel := styles.LabelStyle.Render("📝 内容:")
	contentBox := styles.ContentBoxStyle.Width(m.width - 4).Render(
		contentLabel + "\n" + m.contentArea.View(),
	)
	sections = append(sections, contentBox)

	// ヘルプ
	var helpButtons []string
	helpButtons = append(helpButtons, styles.SuccessStyle.Render("Ctrl+S")+" 保存")
	helpButtons = append(helpButtons, styles.SecondaryButtonStyle.Render("Esc")+" キャンセル")

	help := styles.HelpStyle.Width(m.width - 4).Render(strings.Join(helpButtons, "  "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
