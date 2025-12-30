package ui

import (
	"fmt"

	"guitui/internal/lesson"
	"guitui/internal/ui/components"
)

// FretboardDataBuilder builds data for fretboard rendering
type FretboardDataBuilder struct {
	lesson      *lesson.Lesson
	currentStep int
}

// NewFretboardDataBuilder creates a new builder
func NewFretboardDataBuilder(l *lesson.Lesson, step int) *FretboardDataBuilder {
	return &FretboardDataBuilder{
		lesson:      l,
		currentStep: step,
	}
}

// BuildActiveItems returns active notes for current step
func (b *FretboardDataBuilder) BuildActiveItems() []components.ActiveItem {
	if b.lesson == nil || len(b.lesson.Steps) == 0 {
		return []components.ActiveItem{}
	}

	if b.currentStep < 0 || b.currentStep >= len(b.lesson.Steps) {
		return []components.ActiveItem{}
	}

	activeItems := []components.ActiveItem{}
	step := b.lesson.Steps[b.currentStep]

	for _, marker := range step.Markers {
		activeItems = append(activeItems, components.ActiveItem{
			Marker: marker,
			Order:  b.currentStep + 1, // 1-based
		})
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
