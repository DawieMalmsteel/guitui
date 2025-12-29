package theory

// ScaleFormula lưu khoảng cách các nốt (intervals)
type ScaleFormula []int

var Scales = map[string]ScaleFormula{
	"chromatic":        {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
	"major":            {0, 2, 4, 5, 7, 9, 11},
	"minor":            {0, 2, 3, 5, 7, 8, 10},
	"minor_pentatonic": {0, 3, 5, 7, 10},
	"major_pentatonic": {0, 2, 4, 7, 9},
	"blues":            {0, 3, 5, 6, 7, 10},
	"dorian":           {0, 2, 3, 5, 7, 9, 10},
	"phrygian":         {0, 1, 3, 5, 7, 8, 10},
	"lydian":           {0, 2, 4, 6, 7, 9, 11},
	"mixolydian":       {0, 2, 4, 5, 7, 9, 10},
	"locrian":          {0, 1, 3, 5, 6, 8, 10},
}

// Kiểm tra xem nốt target có nằm trong Scale (gốc Root) không
func IsNoteInScale(target, root Note, scaleName string) bool {
	intervals, ok := Scales[scaleName]
	if !ok {
		return false
	}
	// Tính khoảng cách relative từ Root -> Target
	diff := (int(target) - int(root) + 12) % 12

	for _, interval := range intervals {
		if diff == interval {
			return true
		}
	}
	return false
}
