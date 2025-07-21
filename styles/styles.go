package styles

import "github.com/charmbracelet/lipgloss"

var (
	// カラーパレット
	PrimaryColor   = lipgloss.Color("#539bf5") // GitHub Dark Blue
	SecondaryColor = lipgloss.Color("#8ddb8c") // GitHub Dark Green
	AccentColor    = lipgloss.Color("#e3b341") // GitHub Dark Yellow/Orange
	SuccessColor   = lipgloss.Color("#57ab5a") // GitHub Dark Success Green
	DangerColor    = lipgloss.Color("#f85149") // GitHub Dark Danger Red
	MutedColor     = lipgloss.Color("#768390") // GitHub Dark Gray
	BgColor        = lipgloss.Color("#22272e") // GitHub Dark Background

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
