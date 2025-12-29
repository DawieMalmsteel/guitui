package ui

import (
	"time"

	"guitui/internal/lesson"
	"guitui/internal/theory"
	"guitui/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TickMsg time.Time

type Model struct {
	// State bài học
	currentLesson lesson.Lesson
	currentStep   int // Index của step hiện tại

	// Config
	tuning []theory.Note
	width  int
	height int
}

func NewModel() Model {
	// 1. Setup một bài học mẫu (Hardcode test logic)
	config := &lesson.GeneratorConfig{
		Root:      "A",
		Scale:     "minor_pentatonic",
		StartFret: 5,
		EndFret:   8,
		Direction: "ascending",
	}

	// 2. Generate steps
	steps, _ := lesson.GenerateSteps(config)

	l := lesson.Lesson{
		Title:     "A Minor Pentatonic (Pos 1)",
		BPM:       120, // Nhanh tí cho máu
		Steps:     steps,
		Generator: config,
	}

	return Model{
		currentLesson: l,
		currentStep:   0,
		tuning:        theory.StandardTuning,
	}
}

func (m Model) Init() tea.Cmd {
	// Bắt đầu Tick (Metronome chạy)
	return tick(m.currentLesson.BPM)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case TickMsg:
		// Tăng step lên, hết bài thì quay lại 0 (Loop)
		m.currentStep = (m.currentStep + 1) % len(m.currentLesson.Steps)
		return m, tick(m.currentLesson.BPM)
	}

	return m, nil
}

func (m Model) View() string {
	// Lấy markers của step hiện tại
	markers := m.currentLesson.Steps[m.currentStep].Markers

	// Render Fretboard
	fretboard := components.RenderFretboard(markers, m.tuning, m.width)

	// Trang trí tí
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).Padding(1)

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("🎸 GUITUI - "+m.currentLesson.Title),
		lipgloss.NewStyle().Padding(1).Render(fretboard),
		lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("240")).Render("Press 'q' to quit"),
	)
}

// Helper Tick Metronome
func tick(bpm int) tea.Cmd {
	duration := time.Duration(60000/bpm) * time.Millisecond
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
