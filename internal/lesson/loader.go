package lesson

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"guitui/internal/theory"
)

// LoadLessons đọc file JSON và trả về danh sách bài học
func LoadLessons(path string) ([]Lesson, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("không đọc được file %s: %w", path, err)
	}

	var lessons []Lesson
	if err := json.Unmarshal(data, &lessons); err != nil {
		return nil, fmt.Errorf("json format error: %w", err)
	}

	// Calculate notes for each marker
	for i := range lessons {
		l := &lessons[i]
		l.ActualKey = parseNote(l.KeyStr)
		
		// Calculate note for each marker based on string + fret
		for j := range l.Steps {
			for k := range l.Steps[j].Markers {
				marker := &l.Steps[j].Markers[k]
				if marker.StringIndex >= 0 && marker.StringIndex < 6 {
					openNote := theory.StandardTuning[marker.StringIndex]
					marker.Note = theory.CalculateNote(openNote, marker.Fret)
				}
			}
		}
	}

	return lessons, nil
}

func parseNote(n string) theory.Note {
	n = strings.TrimSpace(n)
	for i, name := range theory.NoteNames {
		if strings.EqualFold(name, n) {
			return theory.Note(i)
		}
	}
	return theory.C
}
