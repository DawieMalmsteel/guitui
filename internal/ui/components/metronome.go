package components

import (
	"fmt"
	"strings"

	"guitui/internal/theory" // Để lấy màu Catppuccin

	"github.com/charmbracelet/lipgloss"
)

var (
	bpmLabelStyle = lipgloss.NewStyle().Foreground(theory.CatSky).Bold(true).PaddingRight(1)

	// Downbeat (Phách 1) - Màu Đỏ
	downBeatStyle = lipgloss.NewStyle().Foreground(theory.CatRed).SetString(" ■ ")

	// Upbeat (Phách thường) - Màu Vàng
	upBeatStyle = lipgloss.NewStyle().Foreground(theory.CatYellow).SetString(" ■ ")

	// Inactive - Màu nền tối
	inactiveStyle = lipgloss.NewStyle().Foreground(theory.CatSurface1).SetString(" □ ")
)

func RenderMetronome(currentBeat int, totalBeats int, bpm int) string {
	var b strings.Builder
	b.WriteString(bpmLabelStyle.Render(fmt.Sprintf("BPM: %d ┃", bpm)))

	for i := 0; i < totalBeats; i++ {
		if i == currentBeat {
			if i == 0 {
				b.WriteString(downBeatStyle.String())
			} else {
				b.WriteString(upBeatStyle.String())
			}
		} else {
			b.WriteString(inactiveStyle.String())
		}
	}
	return b.String()
}
