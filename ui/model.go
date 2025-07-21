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

	// テキストエリアの設定（スタイリング）
	ta := textarea.New()
	ta.Placeholder = "日報の内容を入力..."
	ta.SetWidth(60)
	ta.SetHeight(10)

	model := Model{
		currentView:  ListView,
		reports:      reports,
		cursor:       0,
		contentArea:  ta,
		editingIndex: -1,
		isEditing:    false,
		width:        80,
		height:       24,
	}

	// 起動時に本日の日報がない場合、新規作成し保存
	today := time.Now()
	_, err = data.GetReportByDate(today)
	if err != nil {
		// 本日の日報が存在しない場合
		newReport := models.Report{
			ID:      int(time.Now().UnixNano()), // ユニークなIDを生成
			Content: "",                         // 空の内容で作成
			Date:    today,
		}

		if err := data.SaveReport(newReport); err != nil {
			log.Printf("起動時の日報の自動保存に失敗しました: %v", err)
		} else {
			// 保存成功後、reportsスライスを最新の状態に更新
			updatedReports, err := data.GetAllReports()
			if err != nil {
				log.Printf("日報の再読み込みに失敗しました: %v", err)
			} else {
				model.reports = updatedReports
				// 保存したレポートがリストのどこにあるか再検索し、カーソルを合わせる
				for i, r := range model.reports {
					if r.ID == newReport.ID {
						model.cursor = i
						break
					}
				}
			}
		}
		model.currentView = ListView // ListViewに遷移
		model.contentArea.SetValue("")
		model.isEditing = false
		model.editingIndex = -1 // 編集インデックスをリセット
	}

	return model
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
		if r.Date.Format("2006-01-02") == report.Date.Format("2006-01-02") {
			return report, i, true
		}
	}
	return models.Report{}, -1, false
}
