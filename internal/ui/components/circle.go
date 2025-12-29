package components

import (
	"math"
	"strings"

	"guitui/internal/theory"

	"github.com/charmbracelet/lipgloss"
)

// Map Relative Minor (Giọng thứ song song)
var minorMap = map[theory.Note]string{
	theory.C: "Am", theory.G: "Em", theory.D: "Bm", theory.A: "F#m", theory.E: "C#m", theory.B: "G#m",
	theory.Fs: "D#m", theory.Cs: "Bbm", theory.Gs: "Fm", theory.Ds: "Cm", theory.As: "Gm", theory.F: "Dm",
}

var circleOrder = []theory.Note{
	theory.C, theory.G, theory.D, theory.A, theory.E, theory.B,
	theory.Fs, theory.Cs, theory.Gs, theory.Ds, theory.As, theory.F,
}

func RenderCircle(activeNote theory.Note) string {
	const (
		width   = 36
		height  = 13
		centerX = float64(width) / 2.0
		centerY = float64(height) / 2.0
		rOuterX = 13.0
		rOuterY = 6.0
		rInnerX = 9.5
		rInnerY = 4.2
	)

	canvas := make([][]string, height)
	for i := range canvas {
		canvas[i] = make([]string, width)
		for j := range canvas[i] {
			canvas[i][j] = " "
		}
	}

	angleStep := 2 * math.Pi / 12
	startAngle := -math.Pi / 2

	for i, note := range circleOrder {
		theta := startAngle + float64(i)*angleStep

		majName := theory.NoteNames[note]
		minName := minorMap[note]

		// Default Style
		majStyle := lipgloss.NewStyle().Foreground(theory.CatOverlay1).Bold(true)
		minStyle := lipgloss.NewStyle().Foreground(theory.CatSurface1)

		// Active Style (Highlight Key hiện tại)
		if note == activeNote {
			// Màu Active theo bảng màu nốt nhạc để dễ nhận diện
			color := theory.NoteColors[note]
			majStyle = majStyle.Foreground(color).Underline(true)
			minStyle = minStyle.Foreground(theory.CatText).Bold(true)
		}

		// Vẽ Outer (Major)
		renderAt(canvas, width, height, centerX, centerY, rOuterX, rOuterY, theta, majStyle.Render(majName))
		// Vẽ Inner (Minor)
		renderAt(canvas, width, height, centerX, centerY, rInnerX, rInnerY, theta, minStyle.Render(minName))
	}

	var sb strings.Builder
	for _, row := range canvas {
		sb.WriteString(strings.Join(row, "") + "\n")
	}
	return sb.String() // Đã trim suffix ở logic cũ nếu cần thì trim ở Model
}

func renderAt(canvas [][]string, w, h int, cx, cy, rx, ry, theta float64, text string) {
	txtLen := lipgloss.Width(text)
	x := int(cx+rx*math.Cos(theta)) - txtLen/2
	y := int(cy + ry*math.Sin(theta))

	if y >= 0 && y < h && x >= 0 && x+txtLen < w {
		canvas[y][x] = text
		for k := 1; k < txtLen; k++ {
			if x+k < w {
				canvas[y][x+k] = ""
			}
		}
	}
}
