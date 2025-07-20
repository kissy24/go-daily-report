package ui

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"zan/data"
	"zan/models"
)

type ViewType int

const (
	ListView ViewType = iota
	EditView
)

type Model struct {
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

func InitialModel() Model {
	// データストアの初期化
	if err := data.InitStore(); err != nil {
		log.Fatalf("データストアの初期化に失敗しました: %v", err)
	}

	// 既存の日報データをロード
	reports, err := data.GetAllReports()
	if err != nil {
		log.Fatalf("日報データの読み込みに失敗しました: %v", err)
	}

	// nextID を計算
	nextID, err := data.GetNextID()
	if err != nil {
		log.Fatalf("次のIDの取得に失敗しました: %v", err)
	}

	// テキストエリアの設定（スタイリング）
	ta := textarea.New()
	ta.Placeholder = "日報の内容を入力..."
	ta.SetWidth(60)
	ta.SetHeight(10)

	return Model{
		currentView:  ListView,
		reports:      reports,
		cursor:       0,
		contentArea:  ta,
		nextID:       nextID,
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
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return models.Report{}, -1, false // 日付のパースに失敗
	}

	report, err := data.GetReportByDate(parsedDate)
	if err != nil {
		return models.Report{}, -1, false // レポートが見つからない、またはエラー
	}

	for i, r := range m.reports {
		if r.ID == report.ID {
			return report, i, true
		}
	}
	return models.Report{}, -1, false
}
