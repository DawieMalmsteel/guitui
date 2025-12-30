package components

import (
	"fmt"
	"strings"

	"guitui/internal/lesson"
	"guitui/internal/theory"

	"github.com/charmbracelet/lipgloss"
)

// TechniqueDisplayProps contains data for rendering technique info
type TechniqueDisplayProps struct {
	CurrentStep  lesson.Step
	CurrentIndex int
	TotalSteps   int
}

var (
	techBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theory.CatMauve).
		Padding(0, 1)
	
	techTitleStyle = lipgloss.NewStyle().
		Foreground(theory.CatMauve).
		Bold(true)
	
	techLabelStyle = lipgloss.NewStyle().
		Foreground(theory.CatGreen).
		Bold(true)
	
	techValueStyle = lipgloss.NewStyle().
		Foreground(theory.CatYellow)
	
	pickingDownStyle = lipgloss.NewStyle().
		Foreground(theory.CatRed).
		Bold(true)
	
	pickingUpStyle = lipgloss.NewStyle().
		Foreground(theory.CatBlue).
		Bold(true)
)

// RenderTechniqueInfo renders technique information panel
func RenderTechniqueInfo(props TechniqueDisplayProps) string {
	if len(props.CurrentStep.Markers) == 0 {
		return ""
	}

	var lines []string
	
	// Title
	title := techTitleStyle.Render(fmt.Sprintf("TECHNIQUE INFO - Beat %d/%d", 
		props.CurrentIndex+1, props.TotalSteps))
	lines = append(lines, title)
	
	// Collect all techniques in current step
	techniques := make(map[lesson.TechniqueType]int)
	pickings := make(map[lesson.PickingType]int)
	hasLeftHand := false
	
	for _, marker := range props.CurrentStep.Markers {
		if marker.Technique != lesson.TechNone {
			techniques[marker.Technique]++
			hasLeftHand = true
		}
		if marker.Picking != lesson.PickNone {
			pickings[marker.Picking]++
		}
	}
	
	// Display left hand techniques
	if hasLeftHand {
		lines = append(lines, "")
		lines = append(lines, techLabelStyle.Render("Left Hand:"))
		
		for tech, count := range techniques {
			var desc string
			// Get first marker with this technique for details
			var params lesson.TechniqueParams
			for _, m := range props.CurrentStep.Markers {
				if m.Technique == tech {
					params = m.TechParams
					break
				}
			}
			
			switch tech {
			case lesson.TechBend:
				desc = fmt.Sprintf("Bend → fret %d", params.TargetFret)
			case lesson.TechSlide:
				if params.SlideType == "up" {
					desc = fmt.Sprintf("Slide UP → fret %d", params.TargetFret)
				} else if params.SlideType == "down" {
					desc = fmt.Sprintf("Slide DOWN → fret %d", params.TargetFret)
				} else {
					desc = "Slide"
				}
			case lesson.TechHammer:
				desc = fmt.Sprintf("Hammer-on → fret %d", params.TargetFret)
			case lesson.TechPullOff:
				desc = fmt.Sprintf("Pull-off → fret %d", params.TargetFret)
			case lesson.TechVibrato:
				width := params.VibratoWidth
				if width == "" {
					width = "normal"
				}
				desc = fmt.Sprintf("Vibrato (%s)", width)
			case lesson.TechTap:
				desc = "Tap with right hand"
			case lesson.TechHarmonic:
				desc = "Natural harmonic"
			case lesson.TechPinch:
				desc = "Pinch harmonic"
			case lesson.TechTrill:
				desc = fmt.Sprintf("Trill with fret %d", params.TargetFret)
			}
			
			if count > 1 {
				desc += fmt.Sprintf(" (×%d)", count)
			}
			
			lines = append(lines, "  "+techValueStyle.Render("• "+desc))
		}
	}
	
	// Display right hand (picking) techniques
	if len(pickings) > 0 {
		lines = append(lines, "")
		lines = append(lines, techLabelStyle.Render("Right Hand:"))
		
		for pick, count := range pickings {
			var desc string
			var style lipgloss.Style
			
			switch pick {
			case lesson.PickDown:
				desc = "Down stroke ∏"
				style = pickingDownStyle
			case lesson.PickUp:
				desc = "Up stroke V"
				style = pickingUpStyle
			case lesson.PickAlternate:
				desc = "Alternate picking"
				style = techValueStyle
			case lesson.PickTremolo:
				desc = "Tremolo picking (fast)"
				style = techValueStyle
			case lesson.PickSweep:
				desc = "Sweep picking"
				style = techValueStyle
			case lesson.PickEconomy:
				desc = "Economy picking"
				style = techValueStyle
			}
			
			if count > 1 {
				desc += fmt.Sprintf(" (×%d)", count)
			}
			
			lines = append(lines, "  "+style.Render("• "+desc))
		}
	}
	
	// Display picking pattern if exists
	if props.CurrentStep.PickingPattern != "" {
		lines = append(lines, "")
		lines = append(lines, techLabelStyle.Render("Pattern:"))
		
		// Format pattern with colors
		pattern := props.CurrentStep.PickingPattern
		var formatted []string
		for _, char := range strings.Fields(pattern) {
			switch char {
			case "d":
				formatted = append(formatted, pickingDownStyle.Render("∏"))
			case "u":
				formatted = append(formatted, pickingUpStyle.Render("V"))
			default:
				formatted = append(formatted, char)
			}
		}
		lines = append(lines, "  "+strings.Join(formatted, " "))
	}
	
	// If no techniques, show basic info
	if !hasLeftHand && len(pickings) == 0 && props.CurrentStep.PickingPattern == "" {
		lines = append(lines, "")
		lines = append(lines, techValueStyle.Render("No special techniques"))
	}
	
	content := strings.Join(lines, "\n")
	return techBoxStyle.Render(content)
}
