package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"zan/models"
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
	contentArea  textarea.Model
	nextID       int
	editingIndex int
	isEditing    bool

	// UI状態
	width  int
	height int
}

func InitialModel() Model { // 'initialModel' を 'InitialModel' に変更してエクスポート

	// テキストエリアの設定（スタイリング）
	ta := textarea.New()
	ta.Placeholder = "日報の内容を入力..."
	ta.SetWidth(60)
	ta.SetHeight(10)

	return Model{
		currentView: ListView,
		reports: []models.Report{
			{ID: 1, Content: "要件定義が完了しました。\n次は基本設計に入ります。\n\n✅ 完了項目:\n- 要件ヒアリング\n- 仕様書作成\n\n⚠️ 課題:\n- リソース確保\n- スケジュール調整", Date: time.Now().AddDate(0, 0, -2)},
			{ID: 2, Content: "Go言語のTUIライブラリBubble Teaを学習しました。\n\n📖 学んだこと:\n- Model/Update/Viewアーキテクチャ\n- コンポーネントの使い方\n- キー入力の処理方法\n\n🎯 次の目標:\n- Lipglossでスタイリング\n- 実際のアプリケーション作成", Date: time.Now().AddDate(0, 0, -1)},
			{ID: 3, Content: "SQLCを使ったタイプセーフなSQL操作の環境構築を行いました。\n\n🔧 作業内容:\n- sqlc.yamlの設定\n- マイグレーションファイル作成\n- クエリファイルの準備\n\n💡 所感:\n- 型安全性が向上\n- SQLの記述がより明確に", Date: time.Now()},
			{ID: 4, Content: "SQLCを使ったタイプセーフなSQL操作の環境構築を行いました。\n\n🔧 作業内容:\n- sqlc.yamlの設定\n- マイグレーションファイル作成\n- クエリファイルの準備\n\n💡 所感:\n- 型安全性が向上\n- SQLの記述がより明確に", Date: time.Now().AddDate(0, 0, 2)},
			{ID: 5, Content: "SQLCを使ったタイプセーフなSQL操作の環境構築を行いました。\n\n🔧 作業内容:\n- sqlc.yamlの設定\n- マイグレーションファイル作成\n- クエリファイルの準備\n\n💡 所感:\n- 型安全性が向上\n- SQLの記述がより明確に", Date: time.Now().AddDate(0, 0, 3)},
			{ID: 6, Content: "SQLCを使ったタイプセーフなSQL操作の環境構築を行いました。\n\n🔧 作業内容:\n- sqlc.yamlの設定\n- マイグレーションファイル作成\n- クエリファイルの準備\n\n💡 所感:\n- 型安全性が向上\n- SQLの記述がより明確に", Date: time.Now().AddDate(0, 0, 4)},
			{ID: 7, Content: "SQLCを使ったタイプセーフなSQL操作の環境構築を行いました。\n\n🔧 作業内容:\n- sqlc.yamlの設定\n- マイグレーションファイル作成\n- クエリファイルの準備\n\n💡 所感:\n- 型安全性が向上\n- SQLの記述がより明確に", Date: time.Now().AddDate(0, 0, 5)},
		},
		cursor:       0,
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

// FindReportByDate は指定された日付（YYYY-MM-DD形式）の日報を検索します。
func (m Model) FindReportByDate(date string) (models.Report, int, bool) {
	for i, report := range m.reports {
		if report.Date.Format("2006-01-02") == date {
			return report, i, true
		}
	}
	return models.Report{}, -1, false
}
