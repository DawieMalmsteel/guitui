package components

import (
	"fmt"
	"strings"

	"guitui/internal/lesson"
	"guitui/internal/theory"

	"github.com/charmbracelet/lipgloss"
)

// Config màu mè
var (
	// Màu cho dây đàn và phím rỗng
	fretLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Xám tối

	// Màu cho nốt đang đánh (Active)
	noteStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("232")). // Chữ đen
			Background(lipgloss.Color("208"))  // Nền Cam (nổi bật)

	// Màu cho Finger (Helper)
	fingerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("51")) // Cyan
)

// String Names cho UI (Hiển thị kiểu Tab: e B G D A E)
// Index 5 (High E) -> Index 0 (Low E)
var stringLabels = []string{"E", "A", "D", "G", "B", "e"}

// RenderFretboard vẽ cần đàn
// activeSteps: Các nốt cần hiển thị (thường là steps[currentBeat])
// tuning: Để biết dây buông là nốt gì (cho Drop D sau này)
func RenderFretboard(markers []lesson.Marker, tuning []theory.Note, width int) string {
	var b strings.Builder

	// 1. Tối ưu hóa: Chuyển Slice Markers sang Map để tra cứu O(1)
	// Key của map sẽ là "stringIdx_fretIdx" (VD: "0_5")
	activeMap := make(map[string]lesson.Marker)
	for _, m := range markers {
		key := fmt.Sprintf("%d_%d", m.StringIndex, m.Fret)
		activeMap[key] = m
	}

	// 2. Vẽ Header (Số phím: 0 1 2 3...)
	// Giới hạn vẽ 16 phím thôi cho gọn, hoặc tùy width
	maxFret := 16
	b.WriteString("     ") // Padding cho tên dây
	for f := 0; f <= maxFret; f++ {
		b.WriteString(fmt.Sprintf("%-4d", f))
	}
	b.WriteString("\n")

	// 3. Vẽ từng dây (QUAN TRỌNG: Vẽ từ Dây 5 (High E) xuống Dây 0 (Low E))
	// Vì Tablature thì dây nhỏ nằm trên cùng.
	for s := 5; s >= 0; s-- {
		// Tên dây
		label := stringLabels[s]
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(fmt.Sprintf(" %s ║ ", label)))

		// Vẽ từng phím trên dây này
		for f := 0; f <= maxFret; f++ {
			// Check xem có marker ở vị trí này không?
			key := fmt.Sprintf("%d_%d", s, f)
			if m, ok := activeMap[key]; ok {
				// -- CÓ NỐT --
				noteName := theory.NoteNames[m.Note]

				// Render cục nốt (Tròn trịa tí: " A ")
				display := fmt.Sprintf(" %-2s", noteName)
				b.WriteString(noteStyle.Render(display) + fretLineStyle.Render("|"))
			} else {
				// -- KHÔNG CÓ NỐT (Dây rỗng) --
				// Vẽ ---|
				b.WriteString(fretLineStyle.Render("---|"))
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}
