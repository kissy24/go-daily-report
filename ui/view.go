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
	case CreateView:
		return m.renderCreateView()
	case DetailView:
		return m.renderDetailView()
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

func (m Model) renderCreateView() string {
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
	helpButtons = append(helpButtons, styles.SecondaryButtonStyle.Render("Tab")+" フィールド切替")
	helpButtons = append(helpButtons, styles.SecondaryButtonStyle.Render("Esc")+" キャンセル")

	help := styles.HelpStyle.Width(m.width - 4).Render(strings.Join(helpButtons, "  "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderDetailView() string {
	if len(m.reports) == 0 {
		return styles.ErrorStyle.Render("表示する日報がありません。")
	}

	report := m.reports[m.cursor]
	var sections []string

	// ヘッダー
	header := styles.HeaderStyle.Width(m.width - 4).Render("📖 Daily Report Details")
	sections = append(sections, header)

	// メタ情報
	metaInfo := fmt.Sprintf(
		"%s  %s  %s",
		styles.DateStyle.Render("📅 "+report.Date.Format("2006年01月02日 15:04")),
		styles.MetaStyle.Render("ID: "+fmt.Sprintf("%d", report.ID)),
		styles.MetaStyle.Render(fmt.Sprintf("(%d/%d)", m.cursor+1, len(m.reports))),
	)

	// タイトル
	sections = append(sections, metaInfo)

	// 内容
	contentSection := styles.ContentBoxStyle.Width(m.width - 4).Height(12).Render(
		styles.LabelStyle.Render("📝 内容:") + "\n\n" +
			report.Content,
	)
	sections = append(sections, contentSection)

	// ヘルプ
	var helpButtons []string
	helpButtons = append(helpButtons, styles.PrimaryButtonStyle.Render("E")+" 編集")
	helpButtons = append(helpButtons, styles.ErrorStyle.Render(lipgloss.NewStyle().Padding(0, 2).Bold(true).Render("D"))+" 削除")
	helpButtons = append(helpButtons, styles.SecondaryButtonStyle.Render("Q/Esc")+" 戻る")

	help := styles.HelpStyle.Width(m.width - 4).Render(strings.Join(helpButtons, "  "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
