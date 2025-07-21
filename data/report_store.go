package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"zan/models"
)

const reportsDir = "data/reports"

// InitStore はデータストアを初期化します。
// 日報データを保存するディレクトリが存在しない場合、作成します。
func InitStore() error {
	if _, err := os.Stat(reportsDir); os.IsNotExist(err) {
		return os.MkdirAll(reportsDir, 0755)
	}
	return nil
}

// SaveReport は日報をJSONファイルとして保存します。
// 新規作成の場合は新しいファイルを、更新の場合は既存ファイルを上書きします。
func SaveReport(report models.Report) error {
	// ファイル名を日付形式 (YYYY-MM-DD.json) に変更
	fileName := report.Date.Format("2006-01-02") + ".json"
	filePath := filepath.Join(reportsDir, fileName)

	file, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	return os.WriteFile(filePath, file, 0644)
}

// GetAllReports はすべてのJSONファイルを読み込み、[]models.Report のスライスとして返します。
func GetAllReports() ([]models.Report, error) {
	files, err := os.ReadDir(reportsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Report{}, nil // ディレクトリが存在しない場合は空のスライスを返す
		}
		return nil, fmt.Errorf("failed to read reports directory: %w", err)
	}

	var reports []models.Report
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(reportsDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read report file %s: %w", file.Name(), err)
		}

		var report models.Report
		if err := json.Unmarshal(data, &report); err != nil {
			return nil, fmt.Errorf("failed to unmarshal report file %s: %w", file.Name(), err)
		}
		reports = append(reports, report)
	}
	return reports, nil
}

// GetReportByDate は指定された日付の日報を検索し、models.Report 構造体として返します。
func GetReportByDate(date time.Time) (models.Report, error) {
	fileName := date.Format("2006-01-02") + ".json"
	filePath := filepath.Join(reportsDir, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return models.Report{}, fmt.Errorf("failed to read report file for date %s: %w", date.Format("2006-01-02"), err)
	}

	var report models.Report
	if err := json.Unmarshal(data, &report); err != nil {
		return models.Report{}, fmt.Errorf("failed to unmarshal report for date %s: %w", date.Format("2006-01-02"), err)
	}
	return report, nil
}
