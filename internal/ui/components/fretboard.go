package components

import (
	"fmt"
	"strings"

	"guitui/internal/lesson"
	"guitui/internal/theory"

	"github.com/charmbracelet/lipgloss"
)

// --- STYLES ---
var (
	// Dây đàn & Phím (Màu xám Overlay của Catppuccin)
	fretLineStyle = lipgloss.NewStyle().Foreground(theory.CatOverlay1)

	// Đầu cần đàn (Nut) - Màu chữ chính
	nutStyle = lipgloss.NewStyle().Foreground(theory.CatText).Bold(true)

	// Upcoming Note Styles
	// 1. Ngay lập tức (Dist 1) -> Màu Sky (Xanh da trời nổi bật)
	upcoming1Style = lipgloss.NewStyle().Foreground(theory.CatSky).Bold(true)
	// 2. Tiếp theo (Dist 2) -> Màu Yellow (Vàng)
	upcoming2Style = lipgloss.NewStyle().Foreground(theory.CatYellow)
	// 3. Xa hơn (Dist 3) -> Màu Subtext (Xám sáng)
	upcoming3Style = lipgloss.NewStyle().Foreground(theory.CatSubtext1)

	// Active Note (Đang đánh)
	// Chữ đen (Crust) trên nền Peach (Cam đào) -> Tương phản cao vãi lồn
	activeNoteStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theory.CatCrust).
			Background(theory.CatPeach)

	// Dây buông active: Nền Đỏ
	activeOpenStringStyle = activeNoteStyle.Copy().Background(theory.CatRed)

	upcomingPatterns = []string{" ● ", " : ", " ∴ "}
)

type cellData struct {
	text     string
	style    lipgloss.Style
	priority int
}

type FretboardProps struct {
	ActiveMarkers   []lesson.Marker
	UpcomingMarkers map[string]int // Key: "s_f", Value: distance (1, 2, 3)
	Tuning          []theory.Note
	ScaleConfig     *lesson.GeneratorConfig
	ShowAll         bool
	FretCount       int
}

func RenderFretboard(props FretboardProps) string {
	var b strings.Builder
	grid := make(map[string]cellData)

	// --- LAYER 0: SCALE GHOST NOTES ---
	if props.ShowAll && props.ScaleConfig != nil {
		root := parseNoteSimple(props.ScaleConfig.Root)
		for s := range 6 {
			for f := 0; f <= props.FretCount; f++ {
				note := theory.CalculateNote(props.Tuning[s], f)

				if theory.IsNoteInScale(note, root, props.ScaleConfig.Scale) {
					key := fmt.Sprintf("%d_%d", s, f)

					// Lấy màu từ bảng Catppuccin NoteColors
					color := theory.NoteColors[note]

					// Style chữ màu đó, không nền
					style := lipgloss.NewStyle().Foreground(color)

					grid[key] = cellData{
						text:     fmt.Sprintf("%-3s", theory.NoteNames[note]),
						style:    style,
						priority: 1,
					}
				}
			}
		}
	}

	// --- LAYER 1: UPCOMING NOTES ---
	for key, dist := range props.UpcomingMarkers {
		if dist > 3 {
			continue
		}
		symbol := upcomingPatterns[dist-1]
		var style lipgloss.Style
		switch dist {
		case 1:
			style = upcoming1Style
		case 2:
			style = upcoming2Style
		default:
			style = upcoming3Style
		}
		grid[key] = cellData{text: symbol, style: style, priority: 2}
	}

	// --- LAYER 2: ACTIVE NOTES ---
	for _, m := range props.ActiveMarkers {
		key := fmt.Sprintf("%d_%d", m.StringIndex, m.Fret)

		style := activeNoteStyle
		if m.Fret == 0 {
			style = activeOpenStringStyle
		}

		grid[key] = cellData{
			text:     fmt.Sprintf(" %-2s", theory.NoteNames[m.Note]),
			style:    style,
			priority: 3,
		}
	}

	// --- RENDER ---
	// Header
	b.WriteString("     ")
	for f := 0; f <= props.FretCount; f++ {
		// Số phím màu Lavender (Tím nhạt)
		b.WriteString(lipgloss.NewStyle().Foreground(theory.CatLavender).Render(fmt.Sprintf("%-4d", f)))
	}
	b.WriteString("\n")

	stringLabels := []string{"E", "A", "D", "G", "B", "e"}
	for s := 5; s >= 0; s-- {
		b.WriteString(nutStyle.Render(fmt.Sprintf(" %s ║", stringLabels[s])))

		for f := 0; f <= props.FretCount; f++ {
			key := fmt.Sprintf("%d_%d", s, f)
			if cell, exists := grid[key]; exists {
				b.WriteString(cell.style.Render(cell.text))
				b.WriteString(fretLineStyle.Render("|"))
			} else {
				b.WriteString(fretLineStyle.Render("---|"))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func parseNoteSimple(n string) theory.Note {
	for i, name := range theory.NoteNames {
		if strings.EqualFold(name, strings.TrimSpace(n)) {
			return theory.Note(i)
		}
	}
	return theory.C
}
