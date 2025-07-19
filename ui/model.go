package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"zan/models"
	"zan/styles"
)

type ViewType int

const (
	ListView ViewType = iota
	CreateView
	DetailView
)

type Model struct { // 'model' を 'Model' に変更してエクスポート
	currentView ViewType
	reports     []models.Report
	cursor      int

	// 新規作成・編集用
	titleInput   textinput.Model
	contentArea  textarea.Model
	nextID       int
	editingIndex int
	isEditing    bool

	// UI状態
	width  int
	height int
}

func InitialModel() Model { // 'initialModel' を 'InitialModel' に変更してエクスポート
	// タイトル入力の設定（スタイリング）
	ti := textinput.New()
	ti.Placeholder = "日報のタイトルを入力..."
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(styles.MutedColor)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(styles.PrimaryColor)
	ti.CharLimit = 50
	ti.Width = 50

	// テキストエリアの設定（スタイリング）
	ta := textarea.New()
	ta.Placeholder = "日報の内容を入力..."
	ta.SetWidth(60)
	ta.SetHeight(10)

	return Model{
		currentView: ListView,
		reports: []models.Report{
			{ID: 1, Title: "🚀 プロジェクトA進捗報告", Content: "要件定義が完了しました。\n次は基本設計に入ります。\n\n✅ 完了項目:\n- 要件ヒアリング\n- 仕様書作成\n\n⚠️ 課題:\n- リソース確保\n- スケジュール調整", Date: time.Now().AddDate(0, 0, -2)},
			{ID: 2, Title: "📚 Bubble Tea学習", Content: "Go言語のTUIライブラリBubble Teaを学習しました。\n\n📖 学んだこと:\n- Model/Update/Viewアーキテクチャ\n- コンポーネントの使い方\n- キー入力の処理方法\n\n🎯 次の目標:\n- Lipglossでスタイリング\n- 実際のアプリケーション作成", Date: time.Now().AddDate(0, 0, -1)},
			{ID: 3, Title: "⚙️ SQLCセットアップ", Content: "SQLCを使ったタイプセーフなSQL操作の環境構築を行いました。\n\n🔧 作業内容:\n- sqlc.yamlの設定\n- マイグレーションファイル作成\n- クエリファイルの準備\n\n💡 所感:\n- 型安全性が向上\n- SQLの記述がより明確に", Date: time.Now()},
		},
		cursor:       0,
		titleInput:   ti,
		contentArea:  ta,
		nextID:       4,
		editingIndex: -1,
		isEditing:    false,
		width:        80,
		height:       24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

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
				item := styles.SelectedItemStyle.Render(fmt.Sprintf("▶ %s  %s", dateStr, report.Title))
				items = append(items, item)
			} else {
				item := styles.NormalItemStyle.Render(fmt.Sprintf("  %s  %s", dateStr, report.Title))
				items = append(items, item)
			}
		}

		listBox := styles.ContentBoxStyle.Width(m.width - 4).Render(strings.Join(items, "\n"))
		sections = append(sections, listBox)
	}

	// ヘルプ
	var helpButtons []string
	helpButtons = append(helpButtons, styles.PrimaryButtonStyle.Render("Enter")+" 詳細表示")
	helpButtons = append(helpButtons, styles.PrimaryButtonStyle.Render("N")+" 新規作成")
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

	// タイトル入力
	titleLabel := styles.LabelStyle.Render("📌 タイトル:")
	titleBox := styles.ContentBoxStyle.Width(m.width - 4).Render(
		titleLabel + "\n" + m.titleInput.View(),
	)
	sections = append(sections, titleBox)

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
	titleSection := styles.ContentBoxStyle.Width(m.width - 4).Render(
		styles.LabelStyle.Render("📌 タイトル:") + "\n" +
			styles.TitleStyle.Render(report.Title) + "\n\n" +
			metaInfo,
	)
	sections = append(sections, titleSection)

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
