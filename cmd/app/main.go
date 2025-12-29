package main

import (
	"fmt"
	"os"

	"guitui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen()) // WithAltScreen để chiếm full màn hình
	if _, err := p.Run(); err != nil {
		fmt.Printf("Chết cụ rồi: %v", err)
		os.Exit(1)
	}
}
