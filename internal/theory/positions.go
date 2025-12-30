package theory

type PositionType string

const (
	PositionTypeCAGED PositionType = "caged"
	PositionType3NPS  PositionType = "3nps"
)

type NotePattern struct {
	RelativeFrets []int
}

type Position struct {
	Index         int
	Type          PositionType
	FretSpan      int
	RootStrings   []int
	StartOffset   int
	NotePatterns  [6]NotePattern
	FingerPattern [6][]int
}

type ScalePositions struct {
	ScaleName string
	CAGED     []Position
	ThreeNPS  []Position
}

var AllScalePositions = map[string]ScalePositions{
	"minor_pentatonic": {
		ScaleName: "Minor Pentatonic",
		CAGED: []Position{
			{
				Index:       1,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{6, 4, 1},
				StartOffset: 0,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 3}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 3}},
					{RelativeFrets: []int{0, 3}},
				},
				FingerPattern: [6][]int{
					{1, 4}, {1, 3}, {1, 3}, {1, 3}, {1, 4}, {1, 4},
				},
			},
			{
				Index:       2,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 4},
				StartOffset: 3,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
				},
				FingerPattern: [6][]int{
					{1, 3}, {1, 2, 4}, {1, 3}, {1, 3}, {1, 3}, {1, 3},
				},
			},
			{
				Index:       3,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 5},
				StartOffset: 5,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
				},
				FingerPattern: [6][]int{
					{1, 3}, {1, 3}, {1, 2, 4}, {1, 3}, {1, 3}, {1, 3},
				},
			},
			{
				Index:       4,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{3, 5},
				StartOffset: 7,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
				},
				FingerPattern: [6][]int{
					{1, 3}, {1, 3}, {1, 3}, {1, 2, 4}, {1, 3}, {1, 3},
				},
			},
			{
				Index:       5,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{1, 3, 6},
				StartOffset: 10,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2}},
				},
				FingerPattern: [6][]int{
					{1, 3}, {1, 3}, {1, 3}, {1, 3}, {1, 2, 4}, {1, 3},
				},
			},
		},
		ThreeNPS: []Position{},
	},
	"major_pentatonic": {
		ScaleName: "Major Pentatonic",
		CAGED: []Position{
			{
				Index:       1,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{1, 3, 6},
				StartOffset: 0,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
				},
				FingerPattern: [6][]int{
					{1, 3}, {1, 3}, {1, 3}, {1, 2, 4}, {1, 3}, {1, 3},
				},
			},
			{
				Index:       2,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 4},
				StartOffset: 2,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
				},
				FingerPattern: [6][]int{
					{1, 3}, {1, 3}, {1, 2, 4}, {1, 3}, {1, 3}, {1, 3},
				},
			},
			{
				Index:       3,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 5},
				StartOffset: 4,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
				},
				FingerPattern: [6][]int{
					{1, 3}, {1, 2, 4}, {1, 3}, {1, 3}, {1, 3}, {1, 3},
				},
			},
			{
				Index:       4,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{3, 5},
				StartOffset: 5,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
				},
				FingerPattern: [6][]int{
					{1, 3}, {1, 3}, {1, 3}, {1, 2, 4}, {1, 3}, {1, 3},
				},
			},
			{
				Index:       5,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{1, 4, 6},
				StartOffset: 7,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 3}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 2}},
					{RelativeFrets: []int{0, 3}},
					{RelativeFrets: []int{0, 3}},
				},
				FingerPattern: [6][]int{
					{1, 4}, {1, 3}, {1, 3}, {1, 3}, {1, 4}, {1, 4},
				},
			},
		},
		ThreeNPS: []Position{},
	},
	"blues": {
		ScaleName: "Blues",
		CAGED: []Position{
			{
				Index:       1,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{6, 4, 1},
				StartOffset: 0,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
				},
				FingerPattern: [6][]int{
					{1, 3, 4}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 3, 4}, {1, 3, 4},
				},
			},
			{
				Index:       2,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 4},
				StartOffset: 3,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2, 3}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
				},
				FingerPattern: [6][]int{
					{1, 2, 3}, {1, 2, 3, 4}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3},
				},
			},
			{
				Index:       3,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 5},
				StartOffset: 5,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2, 3}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
				},
				FingerPattern: [6][]int{
					{1, 2, 3}, {1, 2, 3}, {1, 2, 3, 4}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3},
				},
			},
			{
				Index:       4,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{3, 5},
				StartOffset: 7,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2, 3}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
				},
				FingerPattern: [6][]int{
					{1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3, 4}, {1, 2, 3}, {1, 2, 3},
				},
			},
			{
				Index:       5,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{1, 3, 6},
				StartOffset: 10,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2}},
					{RelativeFrets: []int{0, 1, 2, 3}},
					{RelativeFrets: []int{0, 1, 2}},
				},
				FingerPattern: [6][]int{
					{1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3, 4}, {1, 2, 3},
				},
			},
		},
		ThreeNPS: []Position{},
	},
	"major": {
		ScaleName: "Major",
		CAGED: []Position{
			{
				Index:       1,
				Type:        PositionTypeCAGED,
				FretSpan:    5,
				RootStrings: []int{6, 4, 1},
				StartOffset: 0,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{1, 3, 4}},
				},
				FingerPattern: [6][]int{
					{1, 2, 4}, {1, 2, 4}, {1, 3, 4}, {1, 2, 4}, {1, 2, 4}, {1, 3, 4},
				},
			},
			{
				Index:       2,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 4},
				StartOffset: 2,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
				},
				FingerPattern: [6][]int{
					{1, 2, 4}, {1, 3, 4}, {1, 2, 4}, {1, 2, 3}, {1, 2, 4}, {1, 2, 4},
				},
			},
			{
				Index:       3,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 5},
				StartOffset: 4,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
				},
				FingerPattern: [6][]int{
					{1, 2, 4}, {1, 2, 4}, {1, 3, 4}, {1, 2, 4}, {1, 2, 4}, {1, 2, 4},
				},
			},
			{
				Index:       4,
				Type:        PositionTypeCAGED,
				FretSpan:    5,
				RootStrings: []int{3, 5},
				StartOffset: 5,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
				},
				FingerPattern: [6][]int{
					{1, 2, 4}, {1, 2, 4}, {1, 2, 4}, {1, 3, 4}, {1, 2, 4}, {1, 2, 4},
				},
			},
			{
				Index:       5,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{1, 3, 6},
				StartOffset: 7,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{0, 2, 4}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 4}},
				},
				FingerPattern: [6][]int{
					{1, 2, 4}, {1, 2, 4}, {1, 2, 4}, {1, 2, 4}, {1, 3, 4}, {1, 2, 4},
				},
			},
		},
		ThreeNPS: []Position{},
	},
	"minor": {
		ScaleName: "Minor",
		CAGED: []Position{
			{
				Index:       1,
				Type:        PositionTypeCAGED,
				FretSpan:    5,
				RootStrings: []int{6, 4, 1},
				StartOffset: 0,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{1, 3, 4}},
				},
				FingerPattern: [6][]int{
					{1, 2, 3}, {1, 2, 3}, {1, 3, 4}, {1, 2, 3}, {1, 2, 3}, {1, 3, 4},
				},
			},
			{
				Index:       2,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 4},
				StartOffset: 2,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 1, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
				},
				FingerPattern: [6][]int{
					{1, 2, 3}, {1, 3, 4}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3},
				},
			},
			{
				Index:       3,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{2, 5},
				StartOffset: 4,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
				},
				FingerPattern: [6][]int{
					{1, 2, 3}, {1, 2, 3}, {1, 3, 4}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3},
				},
			},
			{
				Index:       4,
				Type:        PositionTypeCAGED,
				FretSpan:    5,
				RootStrings: []int{3, 5},
				StartOffset: 5,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
				},
				FingerPattern: [6][]int{
					{1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 3, 4}, {1, 2, 3}, {1, 2, 3},
				},
			},
			{
				Index:       5,
				Type:        PositionTypeCAGED,
				FretSpan:    4,
				RootStrings: []int{1, 3, 6},
				StartOffset: 7,
				NotePatterns: [6]NotePattern{
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{0, 2, 3}},
					{RelativeFrets: []int{1, 3, 4}},
					{RelativeFrets: []int{0, 2, 3}},
				},
				FingerPattern: [6][]int{
					{1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 3, 4}, {1, 2, 3},
				},
			},
		},
		ThreeNPS: []Position{},
	},
}

func GetScalePositions(scaleName string) (ScalePositions, bool) {
	positions, exists := AllScalePositions[scaleName]
	return positions, exists
}

func GetPosition(scaleName string, posType PositionType, posIndex int) (Position, bool) {
	scalePos, exists := GetScalePositions(scaleName)
	if !exists {
		return Position{}, false
	}

	var positions []Position
	switch posType {
	case PositionTypeCAGED:
		positions = scalePos.CAGED
	case PositionType3NPS:
		positions = scalePos.ThreeNPS
	}

	if len(positions) == 0 {
		return Position{}, false
	}

	if posIndex < 1 || posIndex > len(positions) {
		return Position{}, false
	}

	return positions[posIndex-1], true
}

func GetPositionCount(scaleName string, posType PositionType) int {
	scalePos, exists := GetScalePositions(scaleName)
	if !exists {
		return 0
	}

	switch posType {
	case PositionTypeCAGED:
		return len(scalePos.CAGED)
	case PositionType3NPS:
		return len(scalePos.ThreeNPS)
	default:
		return 0
	}
}

func CalculateFretRange(position Position, rootFret int) (start, end int) {
	start = rootFret + position.StartOffset
	end = start + position.FretSpan - 1
	return
}

func FindRootFretOn6thString(root Note) int {
	diff := int(root) - int(E)
	for diff < 0 {
		diff += 12
	}
	return diff
}
