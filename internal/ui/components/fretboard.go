package components

import (
	"fmt"
	"strings"

	"guitui/internal/lesson"
	"guitui/internal/theory"

	"github.com/charmbracelet/lipgloss"
)

// --- STYLES (Giữ nguyên) ---
var (
	fretLineStyle    = lipgloss.NewStyle().Foreground(theory.CatOverlay1)
	nutStyle         = lipgloss.NewStyle().Foreground(theory.CatText).Bold(true)
	upcomingPatterns = []string{" ● ", " : ", " ∴ "}

	activeNoteStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theory.CatCrust).
			Background(theory.CatPeach)
	activeOpenStringStyle = activeNoteStyle.Copy().Background(theory.CatRed)

	// Styles Ngón Tay (Background màu, Chữ đen)
	fingerBgStyles = map[int]lipgloss.Style{
		0: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatRed),
		1: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatTeal),
		2: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatYellow),
		3: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatPeach),
		4: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatRed),
	}

	// Styles Ngón Tay Đảo Ngược (Background tối, Chữ màu)
	fingerFgStyles = map[int]lipgloss.Style{
		0: lipgloss.NewStyle().Bold(true).Foreground(theory.CatRed).Background(theory.CatBase),
		1: lipgloss.NewStyle().Bold(true).Foreground(theory.CatTeal).Background(theory.CatBase),
		2: lipgloss.NewStyle().Bold(true).Foreground(theory.CatYellow).Background(theory.CatBase),
		3: lipgloss.NewStyle().Bold(true).Foreground(theory.CatPeach).Background(theory.CatBase),
		4: lipgloss.NewStyle().Bold(true).Foreground(theory.CatRed).Background(theory.CatBase),
	}
)

type cellData struct {
	text     string
	style    lipgloss.Style
	priority int
}

type SequenceItem struct {
	Order  int
	Finger int
}

type ActiveItem struct {
	Marker lesson.Marker
	Order  int
}

type UpcomingItem struct {
	Distance int
	Finger   int
}

type FretboardProps struct {
	ActiveItems     []ActiveItem
	UpcomingMarkers map[string]UpcomingItem
	ScaleSequence   map[string]SequenceItem
	Tuning          []theory.Note
	ScaleConfig     *lesson.GeneratorConfig
	ShowAll         bool
	FretCount       int
	ShowScaleShape  bool
	ShowFingers     bool
}

// Helper format số để không vỡ layout (Luôn trả về string độ dài 3)
func formatOrder3Chars(n int) string {
	if n < 10 {
		return fmt.Sprintf(" %d ", n) // " 1 "
	}
	// Nếu >= 10, bỏ padding phải để vừa khít 3 ký tự
	return fmt.Sprintf(" %d", n) // " 10"
}

func RenderFretboard(props FretboardProps) string {
	var b strings.Builder
	grid := make(map[string]cellData)

	// ==========================================================
	// LAYER 0: MAP / BACKGROUND
	// ==========================================================

	if props.ShowScaleShape && props.ScaleConfig != nil {
		// MODE 'S': VẼ FULL SHAPE
		root := parseNoteSimple(props.ScaleConfig.Root)
		startBox := props.ScaleConfig.StartFret
		endBox := props.ScaleConfig.EndFret

		for s := 0; s < 6; s++ {
			for f := 0; f <= props.FretCount; f++ {
				note := theory.CalculateNote(props.Tuning[s], f)

				if theory.IsNoteInScale(note, root, props.ScaleConfig.Scale) {
					if f >= startBox && f <= endBox {
						finger := f - startBox + 1
						if finger < 1 {
							finger = 1
						}
						if finger > 4 {
							finger = 4
						}

						key := fmt.Sprintf("%d_%d", s, f)

						text := "   "
						if seqItem, exists := props.ScaleSequence[key]; exists {
							// FIX LAYOUT: Dùng hàm format 3 ký tự
							text = formatOrder3Chars(seqItem.Order)
						}

						var style lipgloss.Style
						if st, ok := fingerBgStyles[finger]; ok {
							style = st
						} else {
							style = fingerBgStyles[0]
						}

						grid[key] = cellData{text: text, style: style, priority: 1}
					}
				}
			}
		}
	} else if props.ShowAll && props.ScaleConfig != nil {
		// MODE 'TAB'
		root := parseNoteSimple(props.ScaleConfig.Root)
		for s := 0; s < 6; s++ {
			for f := 0; f <= props.FretCount; f++ {
				note := theory.CalculateNote(props.Tuning[s], f)
				if theory.IsNoteInScale(note, root, props.ScaleConfig.Scale) {
					key := fmt.Sprintf("%d_%d", s, f)
					grid[key] = cellData{
						text:     fmt.Sprintf("%-3s", theory.NoteNames[note]),
						style:    lipgloss.NewStyle().Foreground(theory.NoteColors[note]),
						priority: 1,
					}
				}
			}
		}
	}

	// ==========================================================
	// LAYER 1: UPCOMING (Dự báo)
	// ==========================================================

	// FIX LOGIC: Nếu đang bật ShowScaleShape (S) thì KHÔNG hiện Upcoming nữa cho đỡ rối
	if !props.ShowScaleShape {
		for key, item := range props.UpcomingMarkers {
			if item.Distance > 3 {
				continue
			}
			symbol := upcomingPatterns[item.Distance-1]

			var style lipgloss.Style
			if s, ok := fingerFgStyles[item.Finger]; ok {
				style = s
			} else {
				style = fingerFgStyles[0]
			}

			if item.Distance == 1 {
				style = style.Copy().Bold(true)
			} else {
				style = style.Copy().Faint(true)
			}

			grid[key] = cellData{text: symbol, style: style, priority: 2}
		}
	}

	// ==========================================================
	// LAYER 2: ACTIVE (Đang đánh)
	// ==========================================================
	for _, item := range props.ActiveItems {
		m := item.Marker
		key := fmt.Sprintf("%d_%d", m.StringIndex, m.Fret)

		var displayText string
		var style lipgloss.Style

		if props.ShowScaleShape {
			// MODE 'S': Active Note
			// FIX LAYOUT: Dùng hàm format 3 ký tự cho số thứ tự
			displayText = formatOrder3Chars(item.Order)

			if s, ok := fingerFgStyles[m.Finger]; ok {
				style = s
			} else {
				style = fingerFgStyles[0]
			}
			style = style.Copy().Underline(true)

		} else if props.ShowFingers {
			// MODE 'H'
			if m.Finger > 0 {
				displayText = fmt.Sprintf(" %d ", m.Finger)
				if s, ok := fingerBgStyles[m.Finger]; ok {
					style = s
				} else {
					style = activeNoteStyle
				}
			} else {
				displayText = " 0 "
				style = fingerBgStyles[0]
			}
		} else {
			// STANDARD
			displayText = fmt.Sprintf(" %-2s", theory.NoteNames[m.Note])
			style = activeNoteStyle
			if m.Fret == 0 {
				style = activeOpenStringStyle
			}
		}

		grid[key] = cellData{text: displayText, style: style, priority: 3}
	}

	// ==========================================================
	// RENDER
	// ==========================================================
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
