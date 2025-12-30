package lesson

import (
	"guitui/internal/theory"
)

// Marker: Một điểm trên cần đàn
type Marker struct {
	StringIndex int         `json:"string"` // 0-5 (0 = String 6 low E, 5 = String 1 high E)
	Fret        int         `json:"fret"`
	Finger      int         `json:"finger"` // 0: Open, 1-4: Ngón tay
	Note        theory.Note `json:"-"`      // Calculated at runtime
}

// Step: Một bước trong bài học (ví dụ 1 beat đánh 1 nốt hoặc 1 hợp âm)
type Step struct {
	Beat    int      `json:"beat"`
	Markers []Marker `json:"markers"`
}

// Lesson: Cấu trúc bài học tổng thể (load từ JSON)
type Lesson struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	BPM      int    `json:"bpm"`
	KeyStr   string `json:"key"`
	
	// Steps được define thủ công trong JSON
	Steps []Step `json:"steps"`
	
	// Runtime data
	ActualKey theory.Note `json:"-"`
}
