package lesson

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"guitui/internal/theory"
)

// TabParser parses ASCII guitar tab files
type TabParser struct {
	metadata map[string]string
	tabLines map[int]string // stringIndex (0-5) -> tab line
}

// String names to index mapping
var stringNames = map[string]int{
	"e": 5, // High E (string 1)
	"B": 4, // B string (string 2)
	"G": 3, // G string (string 3)
	"D": 2, // D string (string 4)
	"A": 1, // A string (string 5)
	"E": 0, // Low E (string 6)
}

// LoadTabFile loads and parses a .tab file
func LoadTabFile(path string) (*Lesson, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open tab file: %w", err)
	}
	defer file.Close()

	parser := &TabParser{
		metadata: make(map[string]string),
		tabLines: make(map[int]string),
	}

	scanner := bufio.NewScanner(file)
	inTabSection := false

	for scanner.Scan() {
		line := scanner.Text()

		// Parse metadata
		if strings.Contains(line, ":") && !inTabSection {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				parser.metadata[key] = value
			}
		}

		// Parse tab lines
		if strings.Contains(line, "|") {
			inTabSection = true
			parser.parseTabLine(line)
		}

		// Stop at NOTES section
		if strings.HasPrefix(line, "NOTES:") || strings.HasPrefix(line, "LEGEND:") {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Build lesson from parsed data
	return parser.buildLesson()
}

// parseTabLine extracts tab notation from a line
func (p *TabParser) parseTabLine(line string) {
	// Format: "e|-----5f1-----7f3-----|"
	// Extract string name and content
	parts := strings.SplitN(line, "|", 2)
	if len(parts) < 2 {
		return
	}

	stringName := strings.TrimSpace(parts[0])
	content := parts[1]

	// Remove trailing |
	content = strings.TrimSuffix(content, "|")

	if idx, ok := stringNames[stringName]; ok {
		// Append to existing line (multi-bar support)
		if existing, exists := p.tabLines[idx]; exists {
			p.tabLines[idx] = existing + content
		} else {
			p.tabLines[idx] = content
		}
	}
}

// buildLesson converts parsed tab to Lesson structure
func (p *TabParser) buildLesson() (*Lesson, error) {
	lesson := &Lesson{
		Title:    p.metadata["TITLE"],
		Category: strings.ToLower(p.metadata["CATEGORY"]),
		KeyStr:   p.metadata["KEY"],
		Steps:    []Step{},
	}

	// Parse BPM
	if bpmStr := p.metadata["BPM"]; bpmStr != "" {
		if bpm, err := strconv.Atoi(bpmStr); err == nil {
			lesson.BPM = bpm
		}
	}

	// Parse actual key note
	lesson.ActualKey = parseNote(lesson.KeyStr)

	// Parse tab lines into steps
	steps := p.parseSteps()
	lesson.Steps = steps

	return lesson, nil
}

// parseSteps parses tab lines by splitting on | delimiter
func (p *TabParser) parseSteps() []Step {
	if len(p.tabLines) == 0 {
		return []Step{}
	}

	// Split each line by | delimiter
	beatCells := make(map[int][]string) // stringIdx -> array of beat cells
	maxBeats := 0

	for stringIdx := 0; stringIdx < 6; stringIdx++ {
		line, exists := p.tabLines[stringIdx]
		if !exists {
			continue
		}

		// Split by | delimiter - each cell between | is one beat
		cells := strings.Split(line, "|")
		
		// Keep ALL cells between pipes, including empty ones (they are rest beats)
		// Only remove the very first and very last if they're from string start/end
		var filteredCells []string
		for i, cell := range cells {
			// Skip first cell if it's from before the first |
			if i == 0 && strings.TrimSpace(cell) == "" {
				continue
			}
			// Skip last cell if it's from after the final |
			if i == len(cells)-1 && strings.TrimSpace(cell) == "" {
				continue
			}
			// Keep everything else, including empty cells (rest beats)
			filteredCells = append(filteredCells, cell)
		}

		beatCells[stringIdx] = filteredCells
		if len(filteredCells) > maxBeats {
			maxBeats = len(filteredCells)
		}
	}

	// Process each beat (column of cells)
	steps := []Step{}
	beatNumber := 1

	for beatIdx := 0; beatIdx < maxBeats; beatIdx++ {
		// Check if this is a skip beat (all cells empty)
		allEmpty := true
		beatMarkers := []Marker{}

		for stringIdx := 0; stringIdx < 6; stringIdx++ {
			cells, exists := beatCells[stringIdx]
			if !exists {
				// String line doesn't exist, treat as empty for all beats
				continue
			}
			
			var cell string
			if beatIdx >= len(cells) {
				// This string has fewer beats, treat as empty
				cell = ""
			} else {
				cell = strings.TrimSpace(cells[beatIdx])
			}
			
			// Remove all dashes (visual only)
			cell = strings.ReplaceAll(cell, "-", "")
			cell = strings.TrimSpace(cell)

			if cell != "" {
				allEmpty = false
				// Parse notes in this cell
				marker := p.parseCell(stringIdx, cell)
				if marker != nil {
					beatMarkers = append(beatMarkers, *marker)
				}
			}
		}

		// If all cells are empty, it's a skip beat
		if allEmpty {
			// Create empty step (rest)
			step := Step{
				Beat:    beatNumber,
				Markers: []Marker{},
			}
			steps = append(steps, step)
			beatNumber++
		} else if len(beatMarkers) > 0 {
			// Create ONE step with all markers in this beat
			step := Step{
				Beat:    beatNumber,
				Markers: beatMarkers,
			}
			steps = append(steps, step)
			beatNumber++
		}
	}

	return steps
}

// parseCell parses a single beat cell for one string
// Cell format examples: "5(f1)", "7b9", "5/7", "5h7", "12t", "<12>", "x"
func (p *TabParser) parseCell(stringIdx int, cell string) *Marker {
	cell = strings.TrimSpace(cell)
	
	if cell == "" {
		return nil
	}

	// Check for muted note
	if cell == "x" || cell == "X" {
		return &Marker{
			StringIndex: stringIdx,
			Fret:        -1, // Special value for muted
			Finger:      0,
			Note:        theory.C, // Placeholder
			Technique:   TechNone,
		}
	}

	// Check for harmonic notation: <12>
	if strings.HasPrefix(cell, "<") && strings.Contains(cell, ">") {
		fretStr := strings.TrimPrefix(cell, "<")
		fretStr = strings.TrimSuffix(fretStr, ">")
		fret, _ := strconv.Atoi(fretStr)
		
		openNote := theory.StandardTuning[stringIdx]
		note := theory.CalculateNote(openNote, fret)
		
		return &Marker{
			StringIndex: stringIdx,
			Fret:        fret,
			Finger:      0,
			Note:        note,
			Technique:   TechHarmonic,
		}
	}

	// Check if starts with a digit (fret number)
	if len(cell) > 0 && cell[0] >= '0' && cell[0] <= '9' {
		// Extract fret, finger, and technique
		fret, finger := p.extractFretFinger(cell)
		technique, params := p.extractTechnique(cell)
		picking := p.extractPicking(cell)

		// Calculate note
		openNote := theory.StandardTuning[stringIdx]
		note := theory.CalculateNote(openNote, fret)

		return &Marker{
			StringIndex: stringIdx,
			Fret:        fret,
			Finger:      finger,
			Note:        note,
			Technique:   technique,
			TechParams:  params,
			Picking:     picking,
		}
	}

	return nil
}

// extractFretFinger parses fret number and optional finger + picking from cell string
// Supports: "5", "5(f1)", "5(f1:d)", "5(d)", "12(f3:u)"
func (p *TabParser) extractFretFinger(cell string) (fret, finger int) {
	fret = 0
	finger = 0

	// Parse fret number (can be multi-digit: 10, 12, etc.)
	fretStr := ""
	i := 0
	for i < len(cell) && cell[i] >= '0' && cell[i] <= '9' {
		fretStr += string(cell[i])
		i++
	}

	if fretStr != "" {
		fret, _ = strconv.Atoi(fretStr)
	}

	// Check for finger notation in parentheses
	if i < len(cell) && cell[i] == '(' {
		// Find closing parenthesis
		closeIdx := strings.Index(cell[i:], ")")
		if closeIdx != -1 {
			closeIdx += i
			content := cell[i+1 : closeIdx]

			// Parse content: "f1:d", "f1", "d", or just "1" (legacy)
			if len(content) == 1 && content[0] >= '0' && content[0] <= '9' {
				// Legacy format: (1) means finger 1
				finger, _ = strconv.Atoi(content)
			} else if strings.HasPrefix(content, "f") {
				// Modern format: (f1) or (f1:d)
				// Split by : to check for picking
				if strings.Contains(content, ":") {
					// Format: (f1:d)
					parts := strings.Split(content, ":")
					if len(parts) == 2 {
						finger, _ = strconv.Atoi(parts[0][1:]) // Skip 'f' prefix
						// Picking will be parsed separately in extractPicking()
					}
				} else {
					// Format: (f1)
					finger, _ = strconv.Atoi(content[1:])
				}
			}
			// If no 'f' prefix, check if it's just picking: (d) or (u)
			// Will be handled in extractPicking()
		}
	}

	return fret, finger
}

// extractTechnique parses technique notation from cell string
// Supports: 7b9 (bend), 5/7 (slide up), 7\5 (slide down), 5h7 (hammer), 7p5 (pull), 5~ (vibrato), 12t (tap), 5l7 (trill)
func (p *TabParser) extractTechnique(cell string) (TechniqueType, TechniqueParams) {
	params := TechniqueParams{}
	
	// Remove finger notation first to analyze technique
	cellClean := cell
	if idx := strings.Index(cell, "("); idx != -1 {
		cellClean = cell[:idx]
	}
	
	// Check for bend: 7b9
	if strings.Contains(cellClean, "b") {
		parts := strings.Split(cellClean, "b")
		if len(parts) == 2 {
			target, err := strconv.Atoi(parts[1])
			if err == nil {
				params.TargetFret = target
				return TechBend, params
			}
		}
	}
	
	// Check for slide up: 5/7
	if strings.Contains(cellClean, "/") {
		parts := strings.Split(cellClean, "/")
		if len(parts) == 2 {
			target, err := strconv.Atoi(parts[1])
			if err == nil {
				params.TargetFret = target
				params.SlideType = "up"
				return TechSlide, params
			}
		} else if len(parts) == 2 && parts[1] == "" {
			// Slide out up: 5/
			params.SlideType = "out_up"
			return TechSlide, params
		}
	}
	
	// Check for slide down: 7\5
	if strings.Contains(cellClean, "\\") {
		parts := strings.Split(cellClean, "\\")
		if len(parts) == 2 {
			target, err := strconv.Atoi(parts[1])
			if err == nil {
				params.TargetFret = target
				params.SlideType = "down"
				return TechSlide, params
			}
		} else if len(parts) == 2 && parts[1] == "" {
			// Slide out down: 7\
			params.SlideType = "out_down"
			return TechSlide, params
		}
	}
	
	// Check for hammer-on: 5h7
	if strings.Contains(cellClean, "h") {
		parts := strings.Split(cellClean, "h")
		if len(parts) == 2 {
			target, err := strconv.Atoi(parts[1])
			if err == nil {
				params.TargetFret = target
				return TechHammer, params
			}
		}
	}
	
	// Check for pull-off: 7p5
	if strings.Contains(cellClean, "p") {
		parts := strings.Split(cellClean, "p")
		if len(parts) == 2 {
			target, err := strconv.Atoi(parts[1])
			if err == nil {
				params.TargetFret = target
				return TechPullOff, params
			}
		}
	}
	
	// Check for trill: 5l7
	if strings.Contains(cellClean, "l") {
		parts := strings.Split(cellClean, "l")
		if len(parts) == 2 {
			target, err := strconv.Atoi(parts[1])
			if err == nil {
				params.TargetFret = target
				return TechTrill, params
			}
		}
	}
	
	// Check for vibrato: 5~, 5~~, 5~~~
	if strings.Contains(cellClean, "~") {
		tildeCount := strings.Count(cellClean, "~")
		if tildeCount >= 2 {
			params.VibratoWidth = "wide"
		} else {
			params.VibratoWidth = "normal"
		}
		return TechVibrato, params
	}
	
	// Check for tap: 12t
	if strings.HasSuffix(cellClean, "t") {
		return TechTap, params
	}
	
	// Check for pinch harmonic: 7*
	if strings.HasSuffix(cellClean, "*") {
		return TechPinch, params
	}
	
	return TechNone, params
}

// extractPicking parses picking notation from cell string
// Supports: "5(f1:d)", "5(d)", "7(u)", "3(f2:u)"
func (p *TabParser) extractPicking(cell string) PickingType {
	// Find content in parentheses
	startIdx := strings.Index(cell, "(")
	endIdx := strings.Index(cell, ")")
	
	if startIdx == -1 || endIdx == -1 {
		return PickNone
	}
	
	content := cell[startIdx+1 : endIdx]
	
	// Check if content has picking notation
	var pickStr string
	
	if strings.Contains(content, ":") {
		// Format: (f1:d) or (f2:u)
		parts := strings.Split(content, ":")
		if len(parts) == 2 {
			pickStr = strings.TrimSpace(parts[1])
		}
	} else if !strings.HasPrefix(content, "f") && len(content) > 0 {
		// Format: (d) or (u) - no finger, just picking
		pickStr = strings.TrimSpace(content)
	}
	
	// Map picking symbols to types
	switch pickStr {
	case "d":
		return PickDown
	case "u":
		return PickUp
	case "a":
		return PickAlternate
	case "t":
		return PickTremolo
	case "s":
		return PickSweep
	case "e":
		return PickEconomy
	default:
		return PickNone
	}
}

// LoadTabDirectory loads all .tab files from a directory
func LoadTabDirectory(dirPath string) ([]Lesson, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory: %w", err)
	}

	var lessons []Lesson

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".tab") || strings.HasSuffix(name, ".txt") {
			filePath := dirPath + "/" + name
			lesson, err := LoadTabFile(filePath)
			if err != nil {
				fmt.Printf("Warning: failed to load %s: %v\n", name, err)
				continue
			}
			lessons = append(lessons, *lesson)
		}
	}

	return lessons, nil
}

// LoadLessonsFromMultipleSources loads from both JSON and tab files
func LoadLessonsFromMultipleSources(jsonPath, tabDir string) ([]Lesson, error) {
	var allLessons []Lesson

	// Try loading JSON (backward compatibility)
	if jsonPath != "" {
		if _, err := os.Stat(jsonPath); err == nil {
			jsonLessons, err := LoadLessons(jsonPath)
			if err == nil {
				allLessons = append(allLessons, jsonLessons...)
			}
		}
	}

	// Load tab files
	if tabDir != "" {
		if _, err := os.Stat(tabDir); err == nil {
			tabLessons, err := LoadTabDirectory(tabDir)
			if err == nil {
				allLessons = append(allLessons, tabLessons...)
			}
		}
	}

	if len(allLessons) == 0 {
		return nil, fmt.Errorf("no lessons found")
	}

	return allLessons, nil
}

// detectTechnique attempts to parse technique notations (future enhancement)
func detectTechnique(line string, pos int) string {
	// Pattern matchers for techniques
	patterns := map[string]*regexp.Regexp{
		"hammer":    regexp.MustCompile(`\d+h\d+`),
		"pulloff":   regexp.MustCompile(`\d+p\d+`),
		"bend":      regexp.MustCompile(`\d+b\d+`),
		"slide":     regexp.MustCompile(`\d+[/\\]\d+`),
		"vibrato":   regexp.MustCompile(`\d+~+`),
		"trill":     regexp.MustCompile(`\d+l\d+`),
		"tap":       regexp.MustCompile(`\d+t`),
		"harmonic":  regexp.MustCompile(`<\d+>`),
		"pinch":     regexp.MustCompile(`\d+\*`),
	}

	// Check for patterns starting at pos
	substr := line[pos:]
	for technique, pattern := range patterns {
		if pattern.MatchString(substr) {
			return technique
		}
	}

	return ""
}
