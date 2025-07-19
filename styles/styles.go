package styles

import "github.com/charmbracelet/lipgloss"

var (
	// カラーパレット
	PrimaryColor   = lipgloss.Color("#7C3AED") // 紫
	SecondaryColor = lipgloss.Color("#06B6D4") // シアン
	AccentColor    = lipgloss.Color("#F59E0B") // オレンジ
	SuccessColor   = lipgloss.Color("#10B981") // 緑
	DangerColor    = lipgloss.Color("#EF4444") // 赤
	MutedColor     = lipgloss.Color("#6B7280") // グレー
	BgColor        = lipgloss.Color("#1F2937") // ダークグレー

	// スタイル定義
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Background(PrimaryColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 2).
			Margin(0, 0, 1, 0)

	SelectedItemStyle = lipgloss.NewStyle().
				Background(SecondaryColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true).
				Padding(0, 1).
				Margin(0, 1)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")).
			Padding(0, 1).
			Margin(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			Margin(1, 0, 0, 0).
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(MutedColor)

	ContentBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1).
			Margin(0, 0, 1, 0)

	LabelStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			Margin(0, 0, 0, 0)

	MetaStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(DangerColor).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	DateStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	// ボタンスタイル
	PrimaryButtonStyle = lipgloss.NewStyle().
				Background(PrimaryColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 2).
				Margin(0, 1, 0, 0).
				Bold(true)

	SecondaryButtonStyle = lipgloss.NewStyle().
				Background(MutedColor).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 2).
				Margin(0, 1, 0, 0)
)
