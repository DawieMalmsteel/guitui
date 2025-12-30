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

// parseSteps parses tab lines column by column to extract beats
func (p *TabParser) parseSteps() []Step {
	if len(p.tabLines) == 0 {
		return []Step{}
	}

	// Find max length
	maxLen := 0
	for _, line := range p.tabLines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Scan column by column
	steps := []Step{}
	beatNumber := 1

	for col := 0; col < maxLen; col++ {
		markers := []Marker{}
		skip := false

		// Check each string at this column
		for stringIdx := 0; stringIdx < 6; stringIdx++ {
			line, exists := p.tabLines[stringIdx]
			if !exists || col >= len(line) {
				continue
			}

			// Check for skip marker
			if col+1 < len(line) && line[col:col+2] == "sk" {
				skip = true
				continue
			}

			// Check if note starts at this column
			if marker := p.parseNoteAt(stringIdx, line, col); marker != nil {
				markers = append(markers, *marker)
			}
		}

		// Create step if we have markers or skip
		if len(markers) > 0 {
			step := Step{
				Beat:    beatNumber,
				Markers: markers,
			}
			steps = append(steps, step)
			beatNumber++
		} else if skip {
			// Skip beat (rest) - create empty step
			step := Step{
				Beat:    beatNumber,
				Markers: []Marker{},
			}
			steps = append(steps, step)
			beatNumber++
		}
	}

	return steps
}

// parseNoteAt checks if a note starts at the given column position
func (p *TabParser) parseNoteAt(stringIdx int, line string, col int) *Marker {
	if col >= len(line) {
		return nil
	}

	char := line[col]

	// Skip if it's a dash, space, or pipe
	if char == '-' || char == ' ' || char == '|' {
		return nil
	}

	// Check if it's a digit (start of fret number)
	if char >= '0' && char <= '9' {
		// Don't parse if previous char was also a digit (we're in middle of number)
		if col > 0 && line[col-1] >= '0' && line[col-1] <= '9' {
			return nil
		}
		
		// Don't parse if previous char was 'f' (we're the finger number)
		if col > 0 && line[col-1] == 'f' {
			return nil
		}

		// Extract full fret number and optional finger
		fret, finger, _ := p.extractFretAndFinger(line, col)

		// Calculate note
		openNote := theory.StandardTuning[stringIdx]
		note := theory.CalculateNote(openNote, fret)

		return &Marker{
			StringIndex: stringIdx,
			Fret:        fret,
			Finger:      finger,
			Note:        note,
		}
	}

	// Check for muted note
	if char == 'x' || char == 'X' {
		return &Marker{
			StringIndex: stringIdx,
			Fret:        -1, // Special value for muted
			Finger:      0,
			Note:        theory.C, // Placeholder
		}
	}

	return nil
}

// extractFretAndFinger parses fret number and optional finger from position
func (p *TabParser) extractFretAndFinger(line string, start int) (fret, finger int, end int) {
	fret = 0
	finger = 0
	end = start

	// Parse fret number (can be multi-digit: 10, 12, etc.)
	fretStr := ""
	i := start
	for i < len(line) && line[i] >= '0' && line[i] <= '9' {
		fretStr += string(line[i])
		i++
	}

	if fretStr != "" {
		fret, _ = strconv.Atoi(fretStr)
	}

	end = i

	// Check for finger notation: "f1" or "(1)"
	if i < len(line) {
		if line[i] == 'f' && i+1 < len(line) && line[i+1] >= '0' && line[i+1] <= '9' {
			// Format: 5f1
			finger, _ = strconv.Atoi(string(line[i+1]))
			end = i + 2
		} else if line[i] == '(' && i+2 < len(line) && line[i+1] >= '0' && line[i+1] <= '9' && line[i+2] == ')' {
			// Format: 5(1)
			finger, _ = strconv.Atoi(string(line[i+1]))
			end = i + 3
		}
	}

	return fret, finger, end
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
