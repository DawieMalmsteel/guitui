package ui

import (
	"fmt"

	"guitui/internal/lesson"
	"guitui/internal/ui/components"
)

// FretboardDataBuilder builds data for fretboard rendering
type FretboardDataBuilder struct {
	lesson      *lesson.Lesson
	currentBeat int
	currentStep int // Index of current step (derived from currentBeat)
}

// NewFretboardDataBuilder creates a new builder
func NewFretboardDataBuilder(l *lesson.Lesson, beat int) *FretboardDataBuilder {
	// Calculate currentStep from currentBeat
	currentStep := -1
	if l != nil {
		for i, step := range l.Steps {
			if step.Beat <= beat {
				currentStep = i
			} else {
				break
			}
		}
	}
	return &FretboardDataBuilder{
		lesson:      l,
		currentBeat: beat,
		currentStep: currentStep,
	}
}

// BuildActiveItems returns active notes for current beat
// Includes notes from previous steps that are still holding (duration > 1)
func (b *FretboardDataBuilder) BuildActiveItems() []components.ActiveItem {
	if b.lesson == nil || len(b.lesson.Steps) == 0 {
		return []components.ActiveItem{}
	}

	if b.currentBeat < 1 {
		return []components.ActiveItem{}
	}

	activeItems := []components.ActiveItem{}

	// Find all steps that should be active at current beat
	// A step is active if: step.Beat <= currentBeat < step.Beat + duration
	for i, step := range b.lesson.Steps {
		for _, marker := range step.Markers {
			stepBeat := step.Beat
			// Duration defaults to 1 if not specified or 0
			duration := marker.Duration
			if duration <= 0 {
				duration = 1
			}
			stepEndBeat := stepBeat + duration - 1 // Duration includes start beat

			// Check if this marker is active at current beat
			if b.currentBeat >= stepBeat && b.currentBeat <= stepEndBeat {
				activeItems = append(activeItems, components.ActiveItem{
					Marker: marker,
					Order:  i + 1, // 1-based step order
				})
			}
		}
	}

	return activeItems
}

// BuildUpcomingMarkers returns upcoming notes (lookahead)
func (b *FretboardDataBuilder) BuildUpcomingMarkers(lookAhead int) map[string]components.UpcomingItem {
	upcoming := make(map[string]components.UpcomingItem)

	if b.lesson == nil || len(b.lesson.Steps) == 0 {
		return upcoming
	}

	for i := 1; i <= lookAhead; i++ {
		nextIdx := (b.currentStep + i) % len(b.lesson.Steps)
		step := b.lesson.Steps[nextIdx]

		for _, marker := range step.Markers {
			key := fmt.Sprintf("%d_%d", marker.StringIndex, marker.Fret)
			// Only save first occurrence
			if _, exists := upcoming[key]; !exists {
				upcoming[key] = components.UpcomingItem{
					Distance: i,
					Finger:   marker.Finger,
				}
			}
		}
	}

	return upcoming
}

// BuildScaleSequence returns all unique notes in lesson with order
func (b *FretboardDataBuilder) BuildScaleSequence() map[string]components.SequenceItem {
	scaleSeq := make(map[string]components.SequenceItem)

	if b.lesson == nil || len(b.lesson.Steps) == 0 {
		return scaleSeq
	}

	for i, step := range b.lesson.Steps {
		for _, marker := range step.Markers {
			key := fmt.Sprintf("%d_%d", marker.StringIndex, marker.Fret)
			// Only save first occurrence
			if _, exists := scaleSeq[key]; !exists {
				scaleSeq[key] = components.SequenceItem{
					Order:  i + 1, // 1-based
					Finger: marker.Finger,
				}
			}
		}
	}

	return scaleSeq
}

// BuildAll builds all fretboard data at once
func (b *FretboardDataBuilder) BuildAll(showUpcoming bool, lookAhead int) (
	activeItems []components.ActiveItem,
	upcoming map[string]components.UpcomingItem,
	scaleSeq map[string]components.SequenceItem,
) {
	activeItems = b.BuildActiveItems()

	if showUpcoming {
		upcoming = b.BuildUpcomingMarkers(lookAhead)
	} else {
		upcoming = make(map[string]components.UpcomingItem)
	}

	scaleSeq = b.BuildScaleSequence()

	return activeItems, upcoming, scaleSeq
}
