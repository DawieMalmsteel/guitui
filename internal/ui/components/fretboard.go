package components

import (
	"fmt"
	"strings"

	"guitui/internal/lesson"
	"guitui/internal/theory"

	"github.com/charmbracelet/lipgloss"
)

// --- CONSTANTS ---
const (
	fretCellWidth = 3 // Each fret cell is 3 characters wide
)

// --- STYLES ---
var (
	fretLineStyle = lipgloss.NewStyle().Foreground(theory.CatOverlay1)
	nutStyle      = lipgloss.NewStyle().Foreground(theory.CatText).Bold(true)
	inlayStyle    = lipgloss.NewStyle().Foreground(theory.CatSurface1)
	
	// Upcoming note patterns (distance 1, 2, 3)
	upcomingPatterns = []string{" ● ", " : ", " ∴ "}

	// Active note styles
	activeNoteStyle       = lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatPeach)
	activeOpenStringStyle = activeNoteStyle.Copy().Background(theory.CatRed)

	// Finger color mapping (background)
	fingerBgStyles = map[int]lipgloss.Style{
		0: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatRed),
		1: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatTeal),
		2: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatYellow),
		3: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatPeach),
		4: lipgloss.NewStyle().Bold(true).Foreground(theory.CatCrust).Background(theory.CatMauve),
	}

	// Finger color mapping (foreground)
	fingerFgStyles = map[int]lipgloss.Style{
		0: lipgloss.NewStyle().Bold(true).Foreground(theory.CatRed).Background(theory.CatBase),
		1: lipgloss.NewStyle().Bold(true).Foreground(theory.CatTeal).Background(theory.CatBase),
		2: lipgloss.NewStyle().Bold(true).Foreground(theory.CatYellow).Background(theory.CatBase),
		3: lipgloss.NewStyle().Bold(true).Foreground(theory.CatPeach).Background(theory.CatBase),
		4: lipgloss.NewStyle().Bold(true).Foreground(theory.CatMauve).Background(theory.CatBase),
	}
)

// Fret inlay positions (standard guitar)
var (
	singleDotFrets = map[int]bool{3: true, 5: true, 7: true, 9: true, 15: true, 17: true, 19: true, 21: true}
	doubleDotFrets = map[int]bool{12: true, 24: true}
)

// String name labels
var stringLabels = []string{"E", "A", "D", "G", "B", "e"}

// --- DATA STRUCTURES ---

// cellData represents a single fretboard cell with display info
type cellData struct {
	text     string
	style    lipgloss.Style
	priority int // Higher priority overwrites lower
}

// SequenceItem represents a note in scale sequence
type SequenceItem struct {
	Order  int // Beat/step order (1-based)
	Finger int // Finger number (0-4)
}

// ActiveItem represents currently playing note
type ActiveItem struct {
	Marker lesson.Marker
	Order  int // For display
}

// UpcomingItem represents upcoming note
type UpcomingItem struct {
	Distance int // Steps ahead (1, 2, 3)
	Finger   int
}

// FretboardProps contains all data needed to render fretboard
type FretboardProps struct {
	ActiveItems     []ActiveItem          // Currently playing notes
	UpcomingMarkers map[string]UpcomingItem // Upcoming notes
	ScaleSequence   map[string]SequenceItem // All notes in lesson
	Tuning          []theory.Note          // String tuning
	FretCount       int                    // Number of frets to show
	
	// Display modes
	ShowAll        bool // Tab mode: show all notes
	ShowScaleShape bool // S key: show scale pattern
	ShowFingers    bool // H key: show finger numbers
}

// --- HELPER FUNCTIONS ---

// formatOrder3Chars formats a number to 3-char width
func formatOrder3Chars(n int) string {
	if n < 10 {
		return fmt.Sprintf(" %d ", n)
	}
	return fmt.Sprintf(" %d", n)
}

// getFingerStyle returns the appropriate style for a finger
func getFingerStyle(finger int, background bool) lipgloss.Style {
	if background {
		if style, ok := fingerBgStyles[finger]; ok {
			return style
		}
		return fingerBgStyles[0]
	}
	if style, ok := fingerFgStyles[finger]; ok {
		return style
	}
	return fingerFgStyles[0]
}

// --- MAIN RENDER FUNCTION ---

// RenderFretboard renders the guitar fretboard with all display modes
func RenderFretboard(props FretboardProps) string {
	grid := make(map[string]cellData)

	// Build display grid in layers (lower priority first)
	buildBackgroundLayer(grid, props)
	buildUpcomingLayer(grid, props)
	buildActiveLayer(grid, props)

	// Render the grid to string
	output := renderGrid(grid, props)
	
	// Add technique and picking lines below fretboard
	techLine := renderTechniqueLine(props)
	pickLine := renderPickingLine(props)
	
	if techLine != "" {
		output += techLine + "\n"
	}
	if pickLine != "" {
		output += pickLine + "\n"
	}
	
	return output
}

// buildBackgroundLayer builds Layer 0: background display modes
func buildBackgroundLayer(grid map[string]cellData, props FretboardProps) {
	// Mode 1: Tab Mode - Show ALL notes on entire fretboard
	if props.ShowAll {
		buildTabMode(grid, props)
		return
	}

	// Mode 2: Scale Shape Mode - Show scale sequence numbers
	if props.ShowScaleShape {
		buildScaleShapeMode(grid, props)
		return
	}

	// Mode 3: Finger Mode - Show finger numbers
	if props.ShowFingers {
		buildFingerMode(grid, props)
		return
	}
}

// buildTabMode displays all notes on fretboard
func buildTabMode(grid map[string]cellData, props FretboardProps) {
	for s := 0; s < 6; s++ {
		for f := 0; f <= props.FretCount; f++ {
			note := theory.CalculateNote(props.Tuning[s], f)
			key := fmt.Sprintf("%d_%d", s, f)

			// Check if note is in scale sequence
			if seqItem, inScale := props.ScaleSequence[key]; inScale {
				// S + Tab mode: colored background with note name
				style := getFingerStyle(seqItem.Finger, true)
				style = style.Copy().Foreground(lipgloss.Color("#000000")) // Black text
				grid[key] = cellData{
					text:     fmt.Sprintf("%-3s", theory.NoteNames[note]),
					style:    style,
					priority: 1,
				}
			} else {
				// Tab only: note name with note color
				grid[key] = cellData{
					text:     fmt.Sprintf("%-3s", theory.NoteNames[note]),
					style:    lipgloss.NewStyle().Foreground(theory.NoteColors[note]),
					priority: 1,
				}
			}
		}
	}
}

// buildScaleShapeMode displays scale sequence numbers
func buildScaleShapeMode(grid map[string]cellData, props FretboardProps) {
	for key, seqItem := range props.ScaleSequence {
		text := formatOrder3Chars(seqItem.Order)
		style := getFingerStyle(seqItem.Finger, true)
		grid[key] = cellData{
			text:     text,
			style:    style,
			priority: 1,
		}
	}
}

// buildFingerMode displays finger numbers
func buildFingerMode(grid map[string]cellData, props FretboardProps) {
	for key, seqItem := range props.ScaleSequence {
		fingerText := " 0 "
		if seqItem.Finger > 0 {
			fingerText = fmt.Sprintf(" %d ", seqItem.Finger)
		}

		style := getFingerStyle(seqItem.Finger, true)
		grid[key] = cellData{
			text:     fingerText,
			style:    style,
			priority: 1,
		}
	}
}

// buildUpcomingLayer builds Layer 1: upcoming notes
func buildUpcomingLayer(grid map[string]cellData, props FretboardProps) {
	// Don't show upcoming in scale/tab modes
	if props.ShowScaleShape || props.ShowAll {
		return
	}

	for key, item := range props.UpcomingMarkers {
		if item.Distance > 3 {
			continue
		}

		symbol := upcomingPatterns[item.Distance-1]
		style := getFingerStyle(item.Finger, false)

		if item.Distance == 1 {
			style = style.Copy().Bold(true)
		} else {
			style = style.Copy().Faint(true)
		}

		grid[key] = cellData{
			text:     symbol,
			style:    style,
			priority: 2,
		}
	}
}

// buildActiveLayer builds Layer 2: currently playing notes
func buildActiveLayer(grid map[string]cellData, props FretboardProps) {
	for _, item := range props.ActiveItems {
		m := item.Marker
		key := fmt.Sprintf("%d_%d", m.StringIndex, m.Fret)

		var displayText string
		var style lipgloss.Style

		if props.ShowScaleShape {
			// Scale mode: show sequence number with underline
			displayText = formatOrder3Chars(item.Order)
			style = getFingerStyle(m.Finger, false)
			style = style.Copy().Underline(true).Bold(true)
		} else if props.ShowFingers {
			// Finger mode: show finger number with background
			if m.Finger > 0 {
				displayText = fmt.Sprintf(" %d ", m.Finger)
			} else {
				displayText = " 0 "
			}
			style = getFingerStyle(m.Finger, true)
			style = style.Copy().Bold(true).Underline(true)
		} else if props.ShowAll {
			// Tab mode: show fret number only (no inline technique)
			displayText = fmt.Sprintf("%-3d", m.Fret)
			// Inverted colors: background = note color, text = dark
			note := theory.CalculateNote(props.Tuning[m.StringIndex], m.Fret)
			style = lipgloss.NewStyle().
				Bold(true).
				Foreground(theory.CatBase).        // Dark text
				Background(theory.NoteColors[note]) // Note color background
		} else {
			// Default mode: show note name with finger color
			displayText = fmt.Sprintf(" %-2s", theory.NoteNames[m.Note])
			style = getFingerStyle(m.Finger, true)
			style = style.Copy().Bold(true)
		}

		grid[key] = cellData{
			text:     displayText,
			style:    style,
			priority: 3, // Highest priority
		}
	}
}

// formatFretWithTechnique formats fret number with technique notation inline using Unicode
func formatFretWithTechnique(m lesson.Marker) string {
	var result string
	fretStr := fmt.Sprintf("%d", m.Fret)
	
	// Add technique suffix with Unicode symbols
	switch m.Technique {
	case "bend":
		result = fmt.Sprintf("%s↗%d", fretStr, m.TechParams.TargetFret)
	case "slide":
		if m.TechParams.SlideType == "up" {
			result = fmt.Sprintf("%s→%d", fretStr, m.TechParams.TargetFret)
		} else if m.TechParams.SlideType == "down" {
			result = fmt.Sprintf("%s←%d", fretStr, m.TechParams.TargetFret)
		} else {
			result = fretStr
		}
	case "hammer":
		result = fmt.Sprintf("%sʰ%d", fretStr, m.TechParams.TargetFret)
	case "pulloff":
		result = fmt.Sprintf("%sᵖ%d", fretStr, m.TechParams.TargetFret)
	case "vibrato":
		result = fretStr + "~"
	case "tap":
		result = fretStr + "ᵀ"
	case "harmonic":
		result = fmt.Sprintf("%s◊", fretStr)
	case "pinch":
		result = fretStr + "*"
	case "trill":
		result = fmt.Sprintf("%s≈%d", fretStr, m.TechParams.TargetFret)
	default:
		// No technique, just fret number
		result = fretStr
	}
	
	// Pad to 4 chars (new cell width)
	return fmt.Sprintf("%-4s", result)
}

// renderGrid converts the grid to string output
func renderGrid(grid map[string]cellData, props FretboardProps) string {
	var b strings.Builder

	// Header: fret numbers
	b.WriteString("      ") // 6 spaces to align with string label
	for f := 0; f <= props.FretCount; f++ {
		b.WriteString(lipgloss.NewStyle().Foreground(theory.CatLavender).Render(fmt.Sprintf("%-4d", f)))
	}
	b.WriteString("\n")

	// Render each string (top to bottom: e, B, G, D, A, E)
	for s := 5; s >= 0; s-- {
		// String label (6 chars to match header)
		b.WriteString(nutStyle.Render(fmt.Sprintf("  %s ║", stringLabels[s])))

		// Render each fret
		for f := 0; f <= props.FretCount; f++ {
			key := fmt.Sprintf("%d_%d", s, f)

			if cell, exists := grid[key]; exists {
				// Cell has data: render with style
				b.WriteString(cell.style.Render(cell.text))
			} else {
				// Empty cell: show inlay or string
				content := renderEmptyCell(s, f)
				if strings.Contains(content, "○") {
					b.WriteString(inlayStyle.Render(content))
				} else {
					b.WriteString(fretLineStyle.Render(content))
				}
			}

			// Fret line separator
			b.WriteString(fretLineStyle.Render("|"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderEmptyCell returns display for empty fretboard cell
func renderEmptyCell(stringIdx, fret int) string {
	// Check for fret inlays
	isDouble := doubleDotFrets[fret]
	isSingle := singleDotFrets[fret]

	if isDouble {
		// Double dots: show on strings A, D, G, B (indices 1-4)
		if stringIdx >= 1 && stringIdx <= 4 {
			return " ○ "
		}
	} else if isSingle {
		// Single dot: show on strings D, G (indices 2-3)
		if stringIdx == 2 || stringIdx == 3 {
			return " ○ "
		}
	}

	// Default: string
	return "---"
}

// renderTechniqueLine renders the technique notation line below fretboard with Unicode symbols
func renderTechniqueLine(props FretboardProps) string {
	if len(props.ActiveItems) == 0 {
		return ""
	}
	
	// Build technique indicators for each fret position with finger info
	type TechInfo struct {
		Symbol string
		Finger int
	}
	techMap := make(map[int][]TechInfo) // fret -> array of techniques (for multi-note chords)
	
	for _, item := range props.ActiveItems {
		m := item.Marker
		if m.Technique == "" {
			continue
		}
		
		var symbol string
		switch m.Technique {
		case "bend":
			symbol = fmt.Sprintf("↗%d", m.TechParams.TargetFret)
		case "slide":
			if m.TechParams.SlideType == "up" {
				symbol = fmt.Sprintf("→%d", m.TechParams.TargetFret)
			} else if m.TechParams.SlideType == "down" {
				symbol = fmt.Sprintf("←%d", m.TechParams.TargetFret)
			} else {
				symbol = "→"
			}
		case "hammer":
			symbol = fmt.Sprintf("ʰ%d", m.TechParams.TargetFret)
		case "pulloff":
			symbol = fmt.Sprintf("ᵖ%d", m.TechParams.TargetFret)
		case "vibrato":
			symbol = "~"
		case "tap":
			symbol = "ᵀ"
		case "harmonic":
			symbol = "◊"
		case "pinch":
			symbol = "*"
		case "trill":
			symbol = fmt.Sprintf("≈%d", m.TechParams.TargetFret)
		}
		
		if symbol != "" {
			techMap[m.Fret] = append(techMap[m.Fret], TechInfo{
				Symbol: symbol,
				Finger: m.Finger,
			})
		}
	}
	
	if len(techMap) == 0 {
		return ""
	}
	
	// Build the line - NO left spacing
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Foreground(theory.CatGreen).Bold(true).Render("Tech: "))
	
	for f := 0; f <= props.FretCount; f++ {
		if techs, exists := techMap[f]; exists {
			// Multiple techniques on same fret: join with |
			var parts []string
			for _, tech := range techs {
				// Color by finger
				style := lipgloss.NewStyle()
				if fingerStyle, ok := fingerFgStyles[tech.Finger]; ok {
					style = fingerStyle.Copy().Background(lipgloss.NoColor{})
				} else {
					style = lipgloss.NewStyle().Foreground(theory.CatYellow)
				}
				parts = append(parts, style.Render(tech.Symbol))
			}
			
			// Join multiple techniques with |
			content := strings.Join(parts, lipgloss.NewStyle().Foreground(theory.CatOverlay1).Render("|"))
			b.WriteString(fmt.Sprintf("%-3s ", content))
		} else {
			b.WriteString("    ") // 4 spaces for empty fret
		}
	}
	
	return b.String()
}

// renderPickingLine renders the picking notation line below fretboard
func renderPickingLine(props FretboardProps) string {
	if len(props.ActiveItems) == 0 {
		return ""
	}
	
	// Build picking indicators for each fret position
	type PickInfo struct {
		Symbol string
		Finger int
	}
	pickMap := make(map[int][]PickInfo) // fret -> array of picking (for multi-note chords)
	
	for _, item := range props.ActiveItems {
		m := item.Marker
		if m.Picking == "" {
			continue
		}
		
		var symbol string
		switch m.Picking {
		case "down":
			symbol = "∏"
		case "up":
			symbol = "V"
		case "alternate":
			symbol = "d-u"
		case "tremolo":
			symbol = "≈"
		case "sweep":
			symbol = "→"
		case "economy":
			symbol = "e"
		}
		
		if symbol != "" {
			pickMap[m.Fret] = append(pickMap[m.Fret], PickInfo{
				Symbol: symbol,
				Finger: m.Finger,
			})
		}
	}
	
	if len(pickMap) == 0 {
		return ""
	}
	
	// Build the line - NO left spacing
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().Foreground(theory.CatGreen).Bold(true).Render("Pick: "))
	
	for f := 0; f <= props.FretCount; f++ {
		if picks, exists := pickMap[f]; exists {
			// Multiple picks on same fret: join with |
			var parts []string
			for _, pick := range picks {
				// Color by picking type
				style := lipgloss.NewStyle()
				if pick.Symbol == "∏" {
					style = lipgloss.NewStyle().Foreground(theory.CatRed).Bold(true)
				} else if pick.Symbol == "V" {
					style = lipgloss.NewStyle().Foreground(theory.CatBlue).Bold(true)
				} else {
					style = lipgloss.NewStyle().Foreground(theory.CatYellow)
				}
				parts = append(parts, style.Render(pick.Symbol))
			}
			
			// Join multiple picks with |
			content := strings.Join(parts, lipgloss.NewStyle().Foreground(theory.CatOverlay1).Render("|"))
			b.WriteString(fmt.Sprintf("%-3s ", content))
		} else {
			b.WriteString("    ") // 4 spaces for empty fret
		}
	}
	
	return b.String()
}
