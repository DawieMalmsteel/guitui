package ui

import (
	"fmt"
	"time"

	"guitui/internal/lesson"
	"guitui/internal/theory"
	"guitui/internal/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TickMsg time.Time

type Model struct {
	currentLesson lesson.Lesson
	currentStep   int

	// State UI
	tuning    []theory.Note
	width     int
	height    int
	fretCount int  // 12 hoặc 24
	showAll   bool // Show all notes in scale
}

func NewModel() Model {
	// Setup bài học mẫu (A Minor Pentatonic)
	config := &lesson.GeneratorConfig{
		Root:      "A",
		Scale:     "minor_pentatonic",
		StartFret: 5,
		EndFret:   8,
		Direction: "ascending",
	}
	steps, _ := lesson.GenerateSteps(config)

	l := lesson.Lesson{
		Title:     "A Minor Pentatonic Box 1",
		BPM:       100, // Tốc độ vừa phải
		Steps:     steps,
		Generator: config,
	}

	return Model{
		currentLesson: l,
		currentStep:   0,
		tuning:        theory.StandardTuning,
		fretCount:     12,    // Mặc định 12 phím
		showAll:       false, // Mặc định tắt
	}
}

func (m Model) Init() tea.Cmd {
	return tick(m.currentLesson.BPM)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "f": // Toggle Fret 12/24
			if m.fretCount == 12 {
				m.fretCount = 24
			} else {
				m.fretCount = 12
			}

		case "tab": // Toggle Show All Notes
			m.showAll = !m.showAll
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case TickMsg:
		m.currentStep = (m.currentStep + 1) % len(m.currentLesson.Steps)
		return m, tick(m.currentLesson.BPM)
	}

	return m, nil
}

func (m Model) View() string {
	steps := m.currentLesson.Steps

	// 1. Lấy Active Markers
	activeMarkers := steps[m.currentStep].Markers

	// 2. Tính toán Upcoming Markers (Nhìn trước 4 bước)
	upcoming := make(map[string]int)
	lookAhead := 3
	for i := 1; i <= lookAhead; i++ {
		// Tính index tương lai (loop vòng tròn)
		nextIdx := (m.currentStep + i) % len(steps)

		for _, marker := range steps[nextIdx].Markers {
			key := fmt.Sprintf("%d_%d", marker.StringIndex, marker.Fret)

			// Chỉ lưu khoảng cách nhỏ nhất (nếu nốt đó lặp lại nhiều lần)
			if _, exists := upcoming[key]; !exists {
				upcoming[key] = i
			}
		}
	}

	// 3. Chuẩn bị Props cho Fretboard
	props := components.FretboardProps{
		ActiveMarkers:   activeMarkers,
		UpcomingMarkers: upcoming,
		Tuning:          m.tuning,
		ScaleConfig:     m.currentLesson.Generator, // Truyền config scale để nó tính showAll
		ShowAll:         m.showAll,
		FretCount:       m.fretCount,
	}

	fretboard := components.RenderFretboard(props)

	// UI Layout
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Render(fmt.Sprintf("🎸 %s (Step: %d/%d)", m.currentLesson.Title, m.currentStep+1, len(steps)))

	help := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).
		Render(fmt.Sprintf("[F] Toggle Frets (%d)  |  [Tab] Show Scale (%v)  |  [Q] Quit", m.fretCount, m.showAll))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Padding(1).Render(title),
		lipgloss.NewStyle().Padding(0, 2).Render(fretboard),
		lipgloss.NewStyle().Padding(1, 2).Render(help),
	)
}

func tick(bpm int) tea.Cmd {
	return tea.Tick(time.Duration(60000/bpm)*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
