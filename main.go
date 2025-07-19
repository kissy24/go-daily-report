package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewType int

const (
	ListView ViewType = iota
	CreateView
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

	// 新規作成用
	titleInput  textinput.Model
	contentArea textarea.Model
	nextID      int
}

// 初期化：ダミーデータを作成
func initialModel() model {
	// タイトル入力の設定
	ti := textinput.New()
	ti.Placeholder = "日報のタイトルを入力..."
	ti.CharLimit = 50
	ti.Width = 50

	// テキストエリアの設定
	ta := textarea.New()
	ta.Placeholder = "日報の内容を入力..."
	ta.SetWidth(60)
	ta.SetHeight(10)

	return model{
		currentView: ListView,
		reports: []Report{
			{ID: 1, Title: "プロジェクトA進捗報告", Content: "要件定義完了", Date: time.Now().AddDate(0, 0, -2)},
			{ID: 2, Title: "Bubble Tea学習", Content: "UIコンポーネントの使い方を学習", Date: time.Now().AddDate(0, 0, -1)},
			{ID: 3, Title: "SQLCセットアップ", Content: "データベース設計を検討", Date: time.Now()},
		},
		cursor:      0,
		titleInput:  ti,
		contentArea: ta,
		nextID:      4,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 全画面共通のキー処理
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	}

	// 画面別のキー処理
	switch m.currentView {
	case ListView:
		return m.handleListViewKeys(msg)
	case CreateView:
		return m.handleCreateViewKeys(msg)
	}

	return m, nil
}

func (m model) handleListViewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "n":
		m.currentView = CreateView
		// 作成ビューに切り替える時、入力欄にフォーカス
		m.titleInput.Focus()
	case "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "j":
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
		// リストビューに戻る
		m.currentView = ListView
		// 入力内容をリセット
		m.titleInput.SetValue("")
		m.contentArea.SetValue("")
		m.titleInput.Blur()
		return m, nil
	case "ctrl+s":
		// 保存処理
		title := m.titleInput.Value()
		content := m.contentArea.Value()

		if title != "" {
			// 新しい日報を追加
			newReport := Report{
				ID:      m.nextID,
				Title:   title,
				Content: content,
				Date:    time.Now(),
			}
			m.reports = append(m.reports, newReport)
			m.nextID++

			// リストビューに戻る
			m.currentView = ListView
			// 入力内容をリセット
			m.titleInput.SetValue("")
			m.contentArea.SetValue("")
			m.titleInput.Blur()
			// 新しく追加したアイテムにカーソルを移動
			m.cursor = len(m.reports) - 1
		}
		return m, nil
	case "tab":
		// タイトル入力とコンテンツエリア間の移動
		if m.titleInput.Focused() {
			m.titleInput.Blur()
			cmd = m.contentArea.Focus()
		} else {
			m.contentArea.Blur()
			cmd = m.titleInput.Focus()
		}
		return m, cmd
	}

	// 各コンポーネントのUpdate処理
	if m.titleInput.Focused() {
		m.titleInput, cmd = m.titleInput.Update(msg)
	} else {
		m.contentArea, cmd = m.contentArea.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	switch m.currentView {
	case ListView:
		return m.renderListView()
	case CreateView:
		return m.renderCreateView()
	default:
		return "不明な画面"
	}
}

func (m model) renderListView() string {
	s := "📝 日報管理システム\n\n"

	// 日報が空の場合
	if len(m.reports) == 0 {
		s += "日報がありません。\n"
		s += "\nn: 新規作成, q: 終了"
		return s
	}

	// 日報リストを表示
	for i, report := range m.reports {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		// 日付を見やすい形式にフォーマット
		dateStr := report.Date.Format("2006/01/02")
		s += fmt.Sprintf("%s [%s] %s\n", cursor, dateStr, report.Title)
	}

	// ヘルプメッセージ
	s += "\n↑/↓: 移動, n: 新規作成, q: 終了"

	return s
}

func (m model) renderCreateView() string {
	s := "📝 新規日報作成\n\n"

	s += "タイトル:\n"
	s += m.titleInput.View() + "\n\n"

	s += "内容:\n"
	s += m.contentArea.View() + "\n\n"

	s += "Tab: フィールド切り替え, Ctrl+S: 保存, Esc: キャンセル"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("エラーが発生しました: %v", err)
	}
}
