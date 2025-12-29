package theory

import "github.com/charmbracelet/lipgloss"

// Note là số nguyên từ 0-11
type Note int

const (
	C Note = iota
	Cs
	D
	Ds
	E
	F
	Fs
	G
	Gs
	A
	As
	B
)

var NoteColors = map[Note]lipgloss.Color{
	C:  lipgloss.Color("196"), // Red
	Cs: lipgloss.Color("203"), // Light Red
	D:  lipgloss.Color("208"), // Orange
	Ds: lipgloss.Color("215"), // Light Orange
	E:  lipgloss.Color("226"), // Yellow
	F:  lipgloss.Color("46"),  // Green
	Fs: lipgloss.Color("83"),  // Light Green
	G:  lipgloss.Color("33"),  // Blue
	Gs: lipgloss.Color("45"),  // Light Blue
	A:  lipgloss.Color("129"), // Purple
	As: lipgloss.Color("135"), // Light Purple
	B:  lipgloss.Color("201"), // Pink
}

// Dùng cái này để hiển thị ra màn hình
var NoteNames = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

// Tuning chuẩn (Low E -> High E)
// Quan trọng: Index 0 là dây to nhất (E2), Index 5 là dây nhỏ nhất (E4)
var StandardTuning = []Note{E, A, D, G, B, E}

// Tính nốt dựa trên Dây (buông) và Phím
func CalculateNote(openStringNote Note, fret int) Note {
	return Note((int(openStringNote) + fret) % 12)
}
