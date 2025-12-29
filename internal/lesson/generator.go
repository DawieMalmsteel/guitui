package lesson

import (
	"strings"

	"guitui/internal/theory"
)

// GenerateSteps tạo ra các nốt dựa trên cấu hình Scale và Box
func GenerateSteps(config *GeneratorConfig) ([]Step, error) {
	// 1. Parse Root Note
	root := parseNote(config.Root)

	var allMarkers []Marker

	// 2. Quét từ dây Trầm (0) lên dây Cao (5)
	// (Hoặc ngược lại tùy mày muốn bài tập chạy thế nào, ở đây tao quét 0->5)
	for s := 0; s < 6; s++ {
		openNote := theory.StandardTuning[s]

		// Quét qua vùng phím (Box)
		for f := config.StartFret; f <= config.EndFret; f++ {
			noteVal := theory.CalculateNote(openNote, f)

			// Check nếu nốt thuộc Scale
			if theory.IsNoteInScale(noteVal, root, config.Scale) {
				// Logic xếp ngón đơn giản: Phím thấp nhất Box -> Ngón 1
				finger := f - config.StartFret + 1
				if finger < 1 {
					finger = 1
				} // Safety
				if finger > 4 {
					finger = 4
				} // Safety (pentatonic box thường 4 ngón)

				allMarkers = append(allMarkers, Marker{
					StringIndex: s,
					Fret:        f,
					Finger:      finger,
					Note:        noteVal,
				})
			}
		}
	}

	// 3. Xử lý Direction (Đảo ngược nếu cần)
	if config.Direction == "descending" {
		for i, j := 0, len(allMarkers)-1; i < j; i, j = i+1, j-1 {
			allMarkers[i], allMarkers[j] = allMarkers[j], allMarkers[i]
		}
	}

	// 4. Đóng gói vào Steps (Mỗi nốt 1 beat)
	var steps []Step
	for i, m := range allMarkers {
		steps = append(steps, Step{
			Beat:    (i % 4) + 1, // Loop beat 1-2-3-4
			Markers: []Marker{m},
		})
	}

	return steps, nil
}

// Helper parse string -> Note
func parseNote(n string) theory.Note {
	n = strings.TrimSpace(n)
	// Map tay hoặc loop qua NoteNames
	for i, name := range theory.NoteNames {
		if strings.EqualFold(name, n) {
			return theory.Note(i)
		}
	}
	return theory.C // Default ngu
}
