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
	fretLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	nutStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Bold(true)

	// Style cho Upcoming (Dự báo)
	// 1. Ngay lập tức (Dist 1) -> Màu nổi nhất (Cyan/Aqua)
	upcoming1Style = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	// 2. Tiếp theo (Dist 2) -> Màu cảnh báo (Yellow)
	upcoming2Style = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	// 3. Xa hơn (Dist 3) -> Màu chìm (Grey)
	upcoming3Style = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	activeNoteStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("232")).
			Background(lipgloss.Color("208"))

	// Pattern hiển thị theo khoảng cách: 1 -> 2 -> 3
	// Mày thích chấm với 2 chấm đúng không? Đây:
	// Dist 1: ● (To tổ bố)
	// Dist 2: : (Hai chấm)
	// Dist 3: ∴ (Ba chấm - Kiểu toán học vì tao thích thế, hoặc mày đổi thành " . " cũng được)
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
		for s := 0; s < 6; s++ {
			for f := 0; f <= props.FretCount; f++ {
				note := theory.CalculateNote(props.Tuning[s], f)

				// Nếu nốt thuộc Scale
				if theory.IsNoteInScale(note, root, props.ScaleConfig.Scale) {
					key := fmt.Sprintf("%d_%d", s, f)

					// 1. Lấy màu tương ứng với nốt (C=Đỏ, E=Vàng...)
					color := theory.NoteColors[note]

					// 2. Tạo style chỉ có Foreground (Chữ), không có Background
					// Tao thêm Faint(true) cho nó dịu lại một tí, đỡ tranh chấp với nốt chính
					// Nếu thích rực rỡ thì bỏ .Faint(true) đi
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

	// --- LAYER 1: UPCOMING NOTES (Sửa lại theo ý mày) ---
	for key, dist := range props.UpcomingMarkers {
		// Chỉ render nếu khoảng cách <= 3 (theo config pattern)
		if dist > 3 {
			continue
		}

		var style lipgloss.Style
		// Lấy pattern tương ứng (Index = dist - 1)
		// Dist 1 -> patterns[0] -> ●
		// Dist 2 -> patterns[1] -> :
		// Dist 3 -> patterns[2] -> ∴
		symbol := upcomingPatterns[dist-1]

		switch dist {
		case 1:
			style = upcoming1Style
		case 2:
			style = upcoming2Style
		default:
			style = upcoming3Style
		}

		grid[key] = cellData{
			text:     symbol,
			style:    style,
			priority: 2, // Đè lên Ghost note
		}
	}

	// --- LAYER 2: ACTIVE NOTES ---
	for _, m := range props.ActiveMarkers {
		key := fmt.Sprintf("%d_%d", m.StringIndex, m.Fret)

		style := activeNoteStyle
		if m.Fret == 0 {
			style = style.Copy().Background(lipgloss.Color("196"))
		}

		grid[key] = cellData{
			text:     fmt.Sprintf(" %-2s", theory.NoteNames[m.Note]),
			style:    style,
			priority: 3, // Đè lên tất cả
		}
	}

	// --- RENDER ---
	// Header
	b.WriteString("     ")
	for f := 0; f <= props.FretCount; f++ {
		b.WriteString(fmt.Sprintf("%-4d", f))
	}
	b.WriteString("\n")

	// Body (High E -> Low E)
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

