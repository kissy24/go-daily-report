package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// カラーパレット
	primaryColor   = lipgloss.Color("#7C3AED") // 紫
	secondaryColor = lipgloss.Color("#06B6D4") // シアン
	accentColor    = lipgloss.Color("#F59E0B") // オレンジ
	successColor   = lipgloss.Color("#10B981") // 緑
	dangerColor    = lipgloss.Color("#EF4444") // 赤
	mutedColor     = lipgloss.Color("#6B7280") // グレー
	bgColor        = lipgloss.Color("#1F2937") // ダークグレー

	// スタイル定義
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Background(primaryColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 2).
			Margin(0, 0, 1, 0)

	selectedItemStyle = lipgloss.NewStyle().
				Background(secondaryColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true).
				Padding(0, 1).
				Margin(0, 1)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			Padding(0, 1).
			Margin(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Margin(1, 0, 0, 0).
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor)

	contentBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1).
			Margin(0, 0, 1, 0)

	labelStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Margin(0, 0, 0, 0)

	metaStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	dateStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	// ボタンスタイル
	primaryButtonStyle = lipgloss.NewStyle().
				Background(primaryColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 2).
				Margin(0, 1, 0, 0).
				Bold(true)

	secondaryButtonStyle = lipgloss.NewStyle().
				Background(mutedColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 2).
				Margin(0, 1, 0, 0)
)

type ViewType int

const (
	ListView ViewType = iota
	CreateView
	DetailView
)

type Report struct {
	ID      int
	Title   string
	Content string
	Date    time.Time
}

type model struct {
	currentView ViewType
	reports     []Report
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

func initialModel() model {
	// タイトル入力の設定（スタイリング）
	ti := textinput.New()
	ti.Placeholder = "日報のタイトルを入力..."
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(mutedColor)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(primaryColor)
	ti.CharLimit = 50
	ti.Width = 50

	// テキストエリアの設定（スタイリング）
	ta := textarea.New()
	ta.Placeholder = "日報の内容を入力..."
	ta.SetWidth(60)
	ta.SetHeight(10)

	return model{
		currentView: ListView,
		reports: []Report{
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

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

func (m model) handleListViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

func (m model) handleCreateViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
				newReport := Report{
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

func (m model) handleDetailViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

func (m model) View() string {
	switch m.currentView {
	case ListView:
		return m.renderListView()
	case CreateView:
		return m.renderCreateView()
	case DetailView:
		return m.renderDetailView()
	default:
		return errorStyle.Render("不明な画面です")
	}
}

func (m model) renderListView() string {
	var sections []string

	// ヘッダー
	header := headerStyle.Width(m.width - 4).Render("📝 Daily Report Manager")
	sections = append(sections, header)

	if len(m.reports) == 0 {
		emptyMsg := contentBoxStyle.Width(m.width - 4).Render(
			labelStyle.Render("📋 日報がありません") + "\n\n" +
				"新しい日報を作成してください。",
		)
		sections = append(sections, emptyMsg)
	} else {
		// 日報リスト
		var items []string
		for i, report := range m.reports {
			dateStr := dateStyle.Render(report.Date.Format("01/02 15:04"))

			if m.cursor == i {
				item := selectedItemStyle.Render(fmt.Sprintf("▶ %s  %s", dateStr, report.Title))
				items = append(items, item)
			} else {
				item := normalItemStyle.Render(fmt.Sprintf("  %s  %s", dateStr, report.Title))
				items = append(items, item)
			}
		}

		listBox := contentBoxStyle.Width(m.width - 4).Render(strings.Join(items, "\n"))
		sections = append(sections, listBox)
	}

	// ヘルプ
	var helpButtons []string
	helpButtons = append(helpButtons, primaryButtonStyle.Render("Enter")+" 詳細表示")
	helpButtons = append(helpButtons, primaryButtonStyle.Render("N")+" 新規作成")
	helpButtons = append(helpButtons, secondaryButtonStyle.Render("Q")+" 終了")
	helpButtons = append(helpButtons, secondaryButtonStyle.Render("↑↓")+" 移動")

	help := helpStyle.Width(m.width - 4).Render(strings.Join(helpButtons, "  "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m model) renderCreateView() string {
	var sections []string

	// ヘッダー
	var headerText string
	if m.isEditing {
		headerText = "✏️ Edit Daily Report"
	} else {
		headerText = "➕ Create New Daily Report"
	}
	header := headerStyle.Width(m.width - 4).Render(headerText)
	sections = append(sections, header)

	// タイトル入力
	titleLabel := labelStyle.Render("📌 タイトル:")
	titleBox := contentBoxStyle.Width(m.width - 4).Render(
		titleLabel + "\n" + m.titleInput.View(),
	)
	sections = append(sections, titleBox)

	// 内容入力
	contentLabel := labelStyle.Render("📝 内容:")
	contentBox := contentBoxStyle.Width(m.width - 4).Render(
		contentLabel + "\n" + m.contentArea.View(),
	)
	sections = append(sections, contentBox)

	// ヘルプ
	var helpButtons []string
	helpButtons = append(helpButtons, successStyle.Render("Ctrl+S")+" 保存")
	helpButtons = append(helpButtons, secondaryButtonStyle.Render("Tab")+" フィールド切替")
	helpButtons = append(helpButtons, secondaryButtonStyle.Render("Esc")+" キャンセル")

	help := helpStyle.Width(m.width - 4).Render(strings.Join(helpButtons, "  "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m model) renderDetailView() string {
	if len(m.reports) == 0 {
		return errorStyle.Render("表示する日報がありません。")
	}

	report := m.reports[m.cursor]
	var sections []string

	// ヘッダー
	header := headerStyle.Width(m.width - 4).Render("📖 Daily Report Details")
	sections = append(sections, header)

	// メタ情報
	metaInfo := fmt.Sprintf(
		"%s  %s  %s",
		dateStyle.Render("📅 "+report.Date.Format("2006年01月02日 15:04")),
		metaStyle.Render("ID: "+fmt.Sprintf("%d", report.ID)),
		metaStyle.Render(fmt.Sprintf("(%d/%d)", m.cursor+1, len(m.reports))),
	)

	// タイトル
	titleSection := contentBoxStyle.Width(m.width - 4).Render(
		labelStyle.Render("📌 タイトル:") + "\n" +
			titleStyle.Render(report.Title) + "\n\n" +
			metaInfo,
	)
	sections = append(sections, titleSection)

	// 内容
	contentSection := contentBoxStyle.Width(m.width - 4).Height(12).Render(
		labelStyle.Render("📝 内容:") + "\n\n" +
			report.Content,
	)
	sections = append(sections, contentSection)

	// ヘルプ
	var helpButtons []string
	helpButtons = append(helpButtons, primaryButtonStyle.Render("E")+" 編集")
	helpButtons = append(helpButtons, errorStyle.Render(lipgloss.NewStyle().Padding(0, 2).Bold(true).Render("D"))+" 削除")
	helpButtons = append(helpButtons, secondaryButtonStyle.Render("Q/Esc")+" 戻る")

	help := helpStyle.Width(m.width - 4).Render(strings.Join(helpButtons, "  "))
	sections = append(sections, help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("エラーが発生しました: %v", err)
	}
}
