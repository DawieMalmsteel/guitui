package components

import (
	"fmt"
	"strings"

	"guitui/internal/audio"
	"guitui/internal/theory"

	"github.com/charmbracelet/lipgloss"
)

var (
	bpmLabelStyle = lipgloss.NewStyle().Foreground(theory.CatSky).Bold(true).PaddingRight(1)
	downBeatStyle = lipgloss.NewStyle().Foreground(theory.CatRed).SetString(" ■ ")
	upBeatStyle   = lipgloss.NewStyle().Foreground(theory.CatYellow).SetString(" ■ ")
	inactiveStyle = lipgloss.NewStyle().Foreground(theory.CatSurface1).SetString(" □ ")

	metroBoxStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(theory.CatMauve).Padding(0, 1)
	metroTitleStyle = lipgloss.NewStyle().Foreground(theory.CatMauve).Bold(true)
	metroKeyStyle   = lipgloss.NewStyle().Foreground(theory.CatGreen).Bold(true)
	metroValueStyle = lipgloss.NewStyle().Foreground(theory.CatYellow)
	metroHintStyle  = lipgloss.NewStyle().Foreground(theory.CatSubtext1).Faint(true)
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

func RenderMetronomeSettings(bpm int, timeSig audio.TimeSignature, soundType string, isActive bool) string {
	var lines []string

	// Title with status
	status := "⏸ PAUSED"
	statusStyle := lipgloss.NewStyle().Foreground(theory.CatYellow)
	if isActive {
		status = "▶ PLAYING"
		statusStyle = lipgloss.NewStyle().Foreground(theory.CatGreen)
	}
	
	title := metroTitleStyle.Render("♪ METRONOME SETTINGS ♪")
	lines = append(lines, title)
	lines = append(lines, statusStyle.Render(status))
	lines = append(lines, "")
	
	separator := metroHintStyle.Render(strings.Repeat("─", 40))
	lines = append(lines, separator)
	lines = append(lines, "")

	// BPM - Big and prominent
	bpmLabel := metroKeyStyle.Render("TEMPO (BPM)")
	bpmValue := lipgloss.NewStyle().Foreground(theory.CatPeach).Bold(true).Render(fmt.Sprintf("▸ %d ◂", bpm))
	lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, bpmLabel, "  ", bpmValue))
	lines = append(lines, metroHintStyle.Render("  Press [+] to increase, [-] to decrease"))
	lines = append(lines, "")

	// Time Signature
	timeSigStr := ""
	switch timeSig {
	case audio.TimeSig4_4:
		timeSigStr = "4/4 (Common Time)"
	case audio.TimeSig3_4:
		timeSigStr = "3/4 (Waltz)"
	case audio.TimeSig6_8:
		timeSigStr = "6/8 (Compound)"
	case audio.TimeSig2_4:
		timeSigStr = "2/4 (March)"
	}
	tsLabel := metroKeyStyle.Render("TIME SIGNATURE")
	tsValue := metroValueStyle.Render(timeSigStr)
	lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, tsLabel, "  ", tsValue))
	lines = append(lines, metroHintStyle.Render("  [1] 4/4   [2] 3/4   [3] 6/8   [4] 2/4"))
	lines = append(lines, "")

	// Sound Type with description
	soundTypeName := ""
	soundDesc := ""
	switch soundType {
	case "wood":
		soundTypeName = "🪵 Wood Block"
		soundDesc = "Natural wood percussion"
	case "mechanical":
		soundTypeName = "⚙️  Mechanical"
		soundDesc = "Sharp mechanical click"
	case "digital":
		soundTypeName = "🔔 Digital Beep"
		soundDesc = "Clean electronic tone"
	default:
		soundTypeName = "🪵 Wood Block"
		soundDesc = "Natural wood percussion"
	}
	soundLabel := metroKeyStyle.Render("SOUND TYPE")
	soundValue := metroValueStyle.Render(soundTypeName)
	lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, soundLabel, "  ", soundValue))
	lines = append(lines, metroHintStyle.Render("  " + soundDesc))
	lines = append(lines, metroHintStyle.Render("  Press [S] to cycle sounds"))
	lines = append(lines, "")
	
	lines = append(lines, separator)
	lines = append(lines, "")

	// Controls section
	controlsTitle := metroKeyStyle.Render("KEYBOARD CONTROLS")
	lines = append(lines, controlsTitle)
	lines = append(lines, "")
	
	controls := []string{
		"  [Space]  Play / Pause metronome",
		"  [M]      Close this menu",
		"  [+/-]    Adjust tempo",
		"  [1-4]    Change time signature",
		"  [S]      Cycle sound types",
	}
	
	for _, ctrl := range controls {
		lines = append(lines, metroHintStyle.Render(ctrl))
	}

	content := strings.Join(lines, "\n")
	
	// Make the box bigger and centered
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theory.CatMauve).
		Padding(1, 3).
		Align(lipgloss.Center)
	
	return boxStyle.Render(content)
}

