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
	fretLineStyle = lipgloss.NewStyle().Foreground(theory.CatOverlay1)
	nutStyle      = lipgloss.NewStyle().Foreground(theory.CatText).Bold(true)

	// Active Note Styles (Giữ nguyên)
	activeNoteStyle       = lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatPeach)
	activeOpenStringStyle = activeNoteStyle.Copy().Background(theory.CatRed)

	// Finger Backgrounds (Cho Active Note)
	fingerStyles = map[int]lipgloss.Style{
		0: activeOpenStringStyle,
		1: activeNoteStyle.Copy().Background(theory.CatTeal),
		2: activeNoteStyle.Copy().Background(theory.CatYellow),
		3: activeNoteStyle.Copy().Background(theory.CatPeach),
		4: activeNoteStyle.Copy().Background(theory.CatRed),
	}

	// --- NEW: FINGER FOREGROUNDS (Cho Upcoming Patterns) ---
	// Chỉ tô màu chữ, không tô nền
	fingerFgStyles = map[int]lipgloss.Style{
		0: lipgloss.NewStyle().Foreground(theory.CatText),   // Open: Màu trắng
		1: lipgloss.NewStyle().Foreground(theory.CatTeal),   // Index: Xanh
		2: lipgloss.NewStyle().Foreground(theory.CatYellow), // Middle: Vàng
		3: lipgloss.NewStyle().Foreground(theory.CatPeach),  // Ring: Cam
		4: lipgloss.NewStyle().Foreground(theory.CatRed),    // Pinky: Đỏ
	}

	// Upcoming Patterns: Gần -> Xa
	upcomingPatterns = []string{" ● ", " : ", " ∴ "}
)

type cellData struct {
	text     string
	style    lipgloss.Style
	priority int
}

// Struct mới để chứa info cho Upcoming (Distance + Finger)
type UpcomingItem struct {
	Distance int // 1, 2, 3
	Finger   int // 0-4
}

type FretboardProps struct {
	ActiveMarkers []lesson.Marker

	// ĐỔI TỪ map[string]int SANG map[string]UpcomingItem
	UpcomingMarkers map[string]UpcomingItem

	Tuning      []theory.Note
	ScaleConfig *lesson.GeneratorConfig
	ShowAll     bool
	FretCount   int
	ShowFingers bool
}

func RenderFretboard(props FretboardProps) string {
	var b strings.Builder
	grid := make(map[string]cellData)

	// --- LAYER 0: SCALE GHOST NOTES (Giữ nguyên) ---
	if props.ShowAll && props.ScaleConfig != nil {
		root := parseNoteSimple(props.ScaleConfig.Root)
		for s := 0; s < 6; s++ {
			for f := 0; f <= props.FretCount; f++ {
				note := theory.CalculateNote(props.Tuning[s], f)
				if theory.IsNoteInScale(note, root, props.ScaleConfig.Scale) {
					key := fmt.Sprintf("%d_%d", s, f)
					color := theory.NoteColors[note]
					grid[key] = cellData{
						text:     fmt.Sprintf("%-3s", theory.NoteNames[note]),
						style:    lipgloss.NewStyle().Foreground(color),
						priority: 1,
					}
				}
			}
		}
	}

	// --- LAYER 1: UPCOMING NOTES (SỬA LOGIC MÀU THEO NGÓN) ---
	for key, item := range props.UpcomingMarkers {
		if item.Distance > 3 {
			continue
		}

		// 1. Chọn Symbol theo Distance (Gần to, xa nhỏ)
		symbol := upcomingPatterns[item.Distance-1]

		// 2. Chọn Màu theo Finger
		var style lipgloss.Style
		if s, ok := fingerFgStyles[item.Finger]; ok {
			style = s
		} else {
			style = fingerFgStyles[0] // Fallback
		}

		// Nếu là nốt sắp đánh ngay lập tức (Dist 1), cho Bold lên
		if item.Distance == 1 {
			style = style.Copy().Bold(true)
		} else {
			// Nốt xa hơn thì làm mờ đi tí (Faint) để tạo chiều sâu
			style = style.Copy().Faint(true)
		}

		grid[key] = cellData{
			text:     symbol,
			style:    style,
			priority: 2,
		}
	}

	// --- LAYER 2: ACTIVE NOTES (Giữ nguyên) ---
	for _, m := range props.ActiveMarkers {
		key := fmt.Sprintf("%d_%d", m.StringIndex, m.Fret)
		var displayText string
		var style lipgloss.Style

		if props.ShowFingers {
			if m.Finger > 0 {
				displayText = fmt.Sprintf(" %d ", m.Finger)
				if s, ok := fingerStyles[m.Finger]; ok {
					style = s
				} else {
					style = activeNoteStyle
				}
			} else {
				displayText = " 0 "
				style = fingerStyles[0]
			}
		} else {
			displayText = fmt.Sprintf(" %-2s", theory.NoteNames[m.Note])
			style = activeNoteStyle
			if m.Fret == 0 {
				style = activeOpenStringStyle
			}
		}

		grid[key] = cellData{text: displayText, style: style, priority: 3}
	}

	// --- RENDER (Giữ nguyên) ---
	b.WriteString("     ")
	for f := 0; f <= props.FretCount; f++ {
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
