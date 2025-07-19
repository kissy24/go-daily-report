package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"zan/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("エラーが発生しました: %v", err)
	}
}
