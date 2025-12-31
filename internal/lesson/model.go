package lesson

import (
	"guitui/internal/theory"
)

// TechniqueType represents the type of guitar technique
type TechniqueType string

const (
	TechNone     TechniqueType = ""
	TechBend     TechniqueType = "bend"
	TechPreBend  TechniqueType = "prebend"
	TechSlide    TechniqueType = "slide"
	TechHammer   TechniqueType = "hammer"
	TechPullOff  TechniqueType = "pulloff"
	TechVibrato  TechniqueType = "vibrato"
	TechTap      TechniqueType = "tap"
	TechHarmonic TechniqueType = "harmonic"
	TechPinch    TechniqueType = "pinch"
	TechTrill    TechniqueType = "trill"
)

// PickingType represents picking hand technique
type PickingType string

const (
	PickNone      PickingType = ""
	PickDown      PickingType = "down"       // d - down stroke
	PickUp        PickingType = "up"         // u - up stroke
	PickAlternate PickingType = "alternate"  // alternating down/up
	PickTremolo   PickingType = "tremolo"    // very fast repeated
	PickSweep     PickingType = "sweep"      // sweep picking
	PickEconomy   PickingType = "economy"    // economy picking
)

// TechniqueParams contains parameters for guitar techniques
type TechniqueParams struct {
	TargetFret   int    // For slides, hammer-ons, pull-offs
	BendSteps    string // Bend amount: "1", "½", "1½", "2" etc.
	BendRelease  bool   // True if bend has release (r suffix)
	VibratoWidth string // "normal", "wide" for vibrato
	SlideType    string // "up", "down", "in", "out" for slides
}

// Marker: Một điểm trên cần đàn
type Marker struct {
	StringIndex int         `json:"string"` // 0-5 (0 = String 6 low E, 5 = String 1 high E)
	Fret        int         `json:"fret"`
	Finger      int         `json:"finger"` // 0: Open, 1-4: Ngón tay
	Note        theory.Note `json:"-"`      // Calculated at runtime
	Beat        int         `json:"-"`      // Optional manual beat number (used during parsing)
	Duration    int         `json:"duration,omitempty"` // Number of beats to hold (default 1)
	
	// Technique information
	Technique TechniqueType   `json:"technique,omitempty"`
	TechParams TechniqueParams `json:"tech_params,omitempty"`
	
	// Picking information
	Picking PickingType `json:"picking,omitempty"`
}

// Step: Một bước trong bài học (ví dụ 1 beat đánh 1 nốt hoặc 1 hợp âm)
type Step struct {
	Beat    int      `json:"beat"`
	Markers []Marker `json:"markers"`
	
	// Step-level annotations (apply to all markers in this beat)
	PickingPattern string `json:"picking_pattern,omitempty"` // e.g., "d u d u"
	Accent         bool   `json:"accent,omitempty"`          // Accent this beat
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
