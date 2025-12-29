package lesson

import (
	"guitui/internal/theory"
)

// Marker: Một điểm trên cần đàn
type Marker struct {
	StringIndex int // 0 (Low E) -> 5 (High E)
	Fret        int
	Finger      int         // 0: Open, 1-4: Ngón tay
	Note        theory.Note // Nốt thực tế (tính toán được)
}

// Step: Một bước trong bài học (ví dụ 1 beat đánh 1 nốt hoặc 1 hợp âm)
type Step struct {
	Beat    int
	Markers []Marker
}

// Config cho Generator (để trong JSON)
type GeneratorConfig struct {
	Root      string `json:"root"`       // "A", "C#"
	Scale     string `json:"scale"`      // "minor_pentatonic"
	StartFret int    `json:"start_fret"` // Phím bắt đầu Box
	EndFret   int    `json:"end_fret"`   // Phím kết thúc Box
	Direction string `json:"direction"`  // "ascending" | "descending"
}

// Lesson: Cấu trúc bài học tổng thể
type Lesson struct {
	Title     string           `json:"title"`
	Category  string           `json:"category"`
	BPM       int              `json:"bpm"`
	KeyStr    string           `json:"key"` // Key hiển thị
	Generator *GeneratorConfig `json:"generator,omitempty"`

	// Dữ liệu runtime (đã qua xử lý)
	Steps     []Step      `json:"-"`
	ActualKey theory.Note `json:"-"`
}
