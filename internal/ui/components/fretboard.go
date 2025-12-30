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
	fretLineStyle    = lipgloss.NewStyle().Foreground(theory.CatOverlay1)
	nutStyle         = lipgloss.NewStyle().Foreground(theory.CatText).Bold(true)
	upcomingPatterns = []string{" ● ", " : ", " ∴ "}

	// Active Styles (Giữ nguyên)
	activeNoteStyle       = lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatPeach)
	activeOpenStringStyle = activeNoteStyle.Copy().Background(theory.CatRed)

	// Finger Styles (Giữ nguyên)
	fingerBgStyles = map[int]lipgloss.Style{
		0: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatRed),
		1: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatTeal),
		2: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatYellow),
		3: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatPeach),
		4: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatRed),
	}

	fingerFgStyles = map[int]lipgloss.Style{
		0: lipgloss.NewStyle().Bold(true).Foreground(theory.CatRed).Background(theory.CatBase),
		1: lipgloss.NewStyle().Bold(true).Foreground(theory.CatTeal).Background(theory.CatBase),
		2: lipgloss.NewStyle().Bold(true).Foreground(theory.CatYellow).Background(theory.CatBase),
		3: lipgloss.NewStyle().Bold(true).Foreground(theory.CatPeach).Background(theory.CatBase),
		4: lipgloss.NewStyle().Bold(true).Foreground(theory.CatRed).Background(theory.CatBase),
	}

	// --- NEW STYLE: FRET MARKER (INLAY) ---
	// Màu Surface2 (Xám sáng hơn nền chút) để không tranh chấp với nốt nhạc
	inlayStyle = lipgloss.NewStyle().Foreground(theory.CatSurface1)
)

// List các phím có chấm (Standard Guitar)
var (
	singleDotFrets = map[int]bool{3: true, 5: true, 7: true, 9: true, 15: true, 17: true, 19: true, 21: true}
	doubleDotFrets = map[int]bool{12: true, 24: true}
)

// ... (Struct cellData, SequenceItem, ActiveItem, UpcomingItem, FretboardProps GIỮ NGUYÊN) ...

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
	ShowAll         bool
	FretCount       int
	ShowScaleShape  bool
	ShowFingers     bool
}

// Helper format số
func formatOrder3Chars(n int) string {
	if n < 10 {
		return fmt.Sprintf(" %d ", n)
	}
	return fmt.Sprintf(" %d", n)
}

func RenderFretboard(props FretboardProps) string {
	var b strings.Builder
	grid := make(map[string]cellData)

	// ==========================================================
	// LAYER 0: BACKGROUND MODES
	// ==========================================================
	
	// Mode 1: Tab Mode (Tab key) - Show ALL note names on entire fretboard
	if props.ShowAll {
		for s := 0; s < 6; s++ {
			for f := 0; f <= props.FretCount; f++ {
				note := theory.CalculateNote(props.Tuning[s], f)
				key := fmt.Sprintf("%d_%d", s, f)
				
				// Check if this note is in the scale (ScaleSequence)
				if seqItem, inScale := props.ScaleSequence[key]; inScale && props.ShowScaleShape {
					// S + Tab mode: Show note name with finger background color
					var style lipgloss.Style
					if st, ok := fingerBgStyles[seqItem.Finger]; ok {
						style = st
					} else {
						style = fingerBgStyles[0]
					}
					// Black text on colored background for readability
					style = style.Copy().Foreground(lipgloss.Color("#000000"))
					grid[key] = cellData{
						text:     fmt.Sprintf("%-3s", theory.NoteNames[note]),
						style:    style,
						priority: 1,
					}
				} else {
					// Tab only: Show note name with note color
					grid[key] = cellData{
						text:     fmt.Sprintf("%-3s", theory.NoteNames[note]),
						style:    lipgloss.NewStyle().Foreground(theory.NoteColors[note]),
						priority: 1,
					}
				}
			}
		}
	} else if props.ShowScaleShape {
		// Mode 2: Scale Sequence (S key) - Show sequence numbers
		for key, seqItem := range props.ScaleSequence {
			text := formatOrder3Chars(seqItem.Order)
			
			var style lipgloss.Style
			if st, ok := fingerBgStyles[seqItem.Finger]; ok {
				style = st
			} else {
				style = fingerBgStyles[0]
			}
			grid[key] = cellData{text: text, style: style, priority: 1}
		}
	} else if props.ShowFingers {
		// Mode 3: Finger Mode (H key) - Show finger numbers for lesson notes
		for key, seqItem := range props.ScaleSequence {
			fingerText := "   "
			if seqItem.Finger > 0 {
				fingerText = fmt.Sprintf(" %d ", seqItem.Finger)
			} else {
				fingerText = " 0 " // Open string
			}
			
			var style lipgloss.Style
			if st, ok := fingerBgStyles[seqItem.Finger]; ok {
				style = st
			} else {
				style = fingerBgStyles[0]
			}
			grid[key] = cellData{text: fingerText, style: style, priority: 1}
		}
	}

	// ==========================================================
	// LAYER 1: UPCOMING
	// ==========================================================
	if !props.ShowScaleShape && !props.ShowAll {
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
	// LAYER 2: ACTIVE
	// ==========================================================
	for _, item := range props.ActiveItems {
		m := item.Marker
		key := fmt.Sprintf("%d_%d", m.StringIndex, m.Fret)
		var displayText string
		var style lipgloss.Style

		if props.ShowScaleShape {
			// Scale Shape mode: show sequence number
			displayText = formatOrder3Chars(item.Order)
			if s, ok := fingerFgStyles[m.Finger]; ok {
				style = s
			} else {
				style = fingerFgStyles[0]
			}
			style = style.Copy().Underline(true).Bold(true)
		} else if props.ShowFingers {
			// Finger mode: show finger number with bold/underline
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
			style = style.Copy().Bold(true).Underline(true)
		} else {
			// Default mode: show note name
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

	// Loop vẽ dây (5 -> 0)
	for s := 5; s >= 0; s-- {
		b.WriteString(nutStyle.Render(fmt.Sprintf(" %s ║", stringLabels[s])))

		for f := 0; f <= props.FretCount; f++ {
			key := fmt.Sprintf("%d_%d", s, f)

			if cell, exists := grid[key]; exists {
				// Có dữ liệu -> Vẽ ô có màu
				b.WriteString(cell.style.Render(cell.text))
				b.WriteString(fretLineStyle.Render("|"))
			} else {
				// Ô Trống -> Xử lý Inlay (Chấm tròn) ở đây
				content := "---" // Mặc định dây

				// Logic vẽ Inlay
				isSingle := singleDotFrets[f]
				isDouble := doubleDotFrets[f]

				if isDouble {
					// 2 Chấm: Vẽ ở dây 1, 2, 3, 4 (A, D, G, B) - Chừa 2 dây ngoài
					if s >= 1 && s <= 4 {
						content = " ○ "
					}
				} else if isSingle {
					// 1 Chấm: Vẽ ở dây 2, 3 (D, G) - Giữa cần đàn
					if s == 2 || s == 3 {
						content = " ○ "
					}
				}

				// Render inlay với màu mờ (inlayStyle) hoặc màu dây (fretLineStyle)
				if content == "-○-" {
					b.WriteString(inlayStyle.Render(content))
				} else {
					b.WriteString(fretLineStyle.Render(content))
				}

				b.WriteString(fretLineStyle.Render("|"))
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
