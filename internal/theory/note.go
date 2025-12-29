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

// catppuccin color palette
var (
	CatRosewater = lipgloss.Color("#f5e0dc")
	CatFlamingo  = lipgloss.Color("#f2cdcd")
	CatPink      = lipgloss.Color("#f5c2e7")
	CatMauve     = lipgloss.Color("#cba6f7")
	CatRed       = lipgloss.Color("#f38ba8")
	CatMaroon    = lipgloss.Color("#eba0ac")
	CatPeach     = lipgloss.Color("#fab387")
	CatYellow    = lipgloss.Color("#f9e2af")
	CatGreen     = lipgloss.Color("#a6e3a1")
	CatTeal      = lipgloss.Color("#94e2d5")
	CatSky       = lipgloss.Color("#89dceb")
	CatSapphire  = lipgloss.Color("#74c7ec")
	CatBlue      = lipgloss.Color("#89b4fa")
	CatLavender  = lipgloss.Color("#b4befe")
	CatText      = lipgloss.Color("#cdd6f4")
	CatSubtext1  = lipgloss.Color("#bac2de")
	CatOverlay1  = lipgloss.Color("#7f849c") // Màu xám cho dây đàn
	CatSurface1  = lipgloss.Color("#45475a") // Màu nền nhẹ
	CatBase      = lipgloss.Color("#1e1e2e") // Màu nền tối
	CatCrust     = lipgloss.Color("#11111b") // Màu siêu tối
)

var NoteColors = map[Note]lipgloss.Color{
	C:  CatRed,      // C = Đỏ pastel
	Cs: CatMaroon,   // C# = Đỏ thẫm
	D:  CatPeach,    // D = Cam đào
	Ds: CatFlamingo, // D# = Hồng cam
	E:  CatYellow,   // E = Vàng
	F:  CatGreen,    // F = Xanh lá
	Fs: CatTeal,     // F# = Xanh cổ vịt
	G:  CatBlue,     // G = Xanh dương
	Gs: CatSapphire, // G# = Xanh dương đậm
	A:  CatMauve,    // A = Tím
	As: CatLavender, // A# = Tím nhạt
	B:  CatPink,     // B = Hồng
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
