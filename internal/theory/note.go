package theory

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

// Dùng cái này để hiển thị ra màn hình
var NoteNames = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

// Tuning chuẩn (Low E -> High E)
// Quan trọng: Index 0 là dây to nhất (E2), Index 5 là dây nhỏ nhất (E4)
var StandardTuning = []Note{E, A, D, G, B, E}

// Tính nốt dựa trên Dây (buông) và Phím
func CalculateNote(openStringNote Note, fret int) Note {
	return Note((int(openStringNote) + fret) % 12)
}
