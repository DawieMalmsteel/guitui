package ui

import (
	"fmt"
	"strings"
	"time"

	"guitui/internal/lesson"
	"guitui/internal/theory"
	"guitui/internal/ui/components"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- CONSTANTS ---
const (
	circleWidth      = 38
	circleHeight     = 13
	listHeaderHeight = 2
)

type TickMsg time.Time

// Wrapper cho list item
type item struct {
	lesson lesson.Lesson
}

func (i item) Title() string { return i.lesson.Title }
func (i item) Description() string {
	return fmt.Sprintf("Key: %s | BPM: %d", i.lesson.KeyStr, i.lesson.BPM)
}
func (i item) FilterValue() string { return i.lesson.Title }

type Model struct {
	// Logic Data
	lessons       []lesson.Lesson
	currentLesson lesson.Lesson
	currentStep   int

	// UI State
	list            list.Model
	tuning          []theory.Note
	width, height   int
	fretCount       int
	metronomeActive bool

	// Display Modes
	showAll        bool // Tab Mode (Note Names) - Phím Tab
	showFingers    bool // Finger Helper Mode - Phím H
	showScaleShape bool // Sequence/Shape Mode - Phím S
}

func NewModel() Model {
	// 1. Load Data
	loadedLessons, err := lesson.LoadLessons("lessons.json")
	if err != nil {
		fmt.Println("Lỗi load lessons.json:", err)
		loadedLessons = []lesson.Lesson{}
	}

	// 2. Setup List Component
	var items []list.Item
	for _, l := range loadedLessons {
		items = append(items, item{lesson: l})
	}

	// Custom Delegate hiển thị list kiểu Catppuccin
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(theory.CatRed).
		BorderLeftForeground(theory.CatRed)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(theory.CatFlamingo).
		BorderLeftForeground(theory.CatRed)

	l := list.New(items, delegate, 0, 0)
	l.Title = "GUITAR LESSONS"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Foreground(theory.CatMauve).Bold(true).Padding(0, 1)

	// 3. Default Lesson
	firstLesson := lesson.Lesson{}
	if len(loadedLessons) > 0 {
		firstLesson = loadedLessons[0]
		// Generate Steps ngay nếu có Generator
		if firstLesson.Generator != nil {
			firstLesson.Steps, _ = lesson.GenerateSteps(firstLesson.Generator)
		}
	}

	return Model{
		lessons:         loadedLessons,
		currentLesson:   firstLesson,
		list:            l,
		tuning:          theory.StandardTuning,
		fretCount:       12,
		metronomeActive: false,

		// Default States
		showAll:        false,
		showFingers:    false,
		showScaleShape: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "f":
			// Toggle 12 <-> 24
			m.fretCount = 36 - m.fretCount
			if m.fretCount != 24 {
				m.fretCount = 12
			}

		case "tab":
			m.showAll = !m.showAll
			// Nếu bật Tab (Note names) thì tắt Shape đi cho đỡ loạn
			if m.showAll {
				m.showScaleShape = false
			}

		case "h", "H":
			m.showFingers = !m.showFingers

		case "s", "S":
			m.showScaleShape = !m.showScaleShape
			// Nếu bật Shape (Sequence) thì tắt Tab đi
			if m.showScaleShape {
				m.showAll = false
			}

		case " ": // Space: Toggle Play/Pause
			m.metronomeActive = !m.metronomeActive
			if m.metronomeActive {
				cmds = append(cmds, tick(m.currentLesson.BPM))
			}

		case "enter": // Chọn bài
			if selectedItem, ok := m.list.SelectedItem().(item); ok {
				m.currentLesson = selectedItem.lesson

				// Generate Steps nếu chưa có
				if len(m.currentLesson.Steps) == 0 && m.currentLesson.Generator != nil {
					steps, err := lesson.GenerateSteps(m.currentLesson.Generator)
					if err == nil {
						m.currentLesson.Steps = steps
					}
				}

				// Reset
				m.currentStep = 0

				// FIX BUG TĂNG TỐC ĐỘ: Chỉ start tick mới nếu đang dừng
				if !m.metronomeActive {
					m.metronomeActive = true
					cmds = append(cmds, tick(m.currentLesson.BPM))
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Resize List
		listWidth := msg.Width - circleWidth - 4
		if listWidth < 20 {
			listWidth = 20
		}
		listHeight := circleHeight - listHeaderHeight

		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)

	case TickMsg:
		if m.metronomeActive && len(m.currentLesson.Steps) > 0 {
			m.currentStep = (m.currentStep + 1) % len(m.currentLesson.Steps)
			cmds = append(cmds, tick(m.currentLesson.BPM))
		}
	}

	// Update List
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	steps := m.currentLesson.Steps

	// --- 1. PREPARE FRETBOARD PROPS ---

	// A. Active Items (Thêm Order cho Mode S nháy)
	var activeItems []components.ActiveItem
	if len(steps) > 0 {
		markers := steps[m.currentStep].Markers
		for _, marker := range markers {
			activeItems = append(activeItems, components.ActiveItem{
				Marker: marker,
				Order:  m.currentStep + 1, // 1-based index
			})
		}
	}

	// B. Upcoming Markers (Lookahead 3 bước)
	upcoming := make(map[string]components.UpcomingItem)
	lookAhead := 3
	if len(steps) > 0 {
		for i := 1; i <= lookAhead; i++ {
			nextIdx := (m.currentStep + i) % len(steps)
			for _, marker := range steps[nextIdx].Markers {
				key := fmt.Sprintf("%d_%d", marker.StringIndex, marker.Fret)
				if _, exists := upcoming[key]; !exists {
					upcoming[key] = components.UpcomingItem{
						Distance: i,
						Finger:   marker.Finger,
					}
				}
			}
		}
	}

	// C. Scale Sequence (Map toàn bộ nốt trong bài để vẽ Layer 0 Mode S)
	scaleSequence := make(map[string]components.SequenceItem)
	if m.showScaleShape && len(steps) > 0 {
		for i, step := range steps {
			for _, marker := range step.Markers {
				key := fmt.Sprintf("%d_%d", marker.StringIndex, marker.Fret)
				// Chỉ lưu lần xuất hiện đầu tiên
				if _, exists := scaleSequence[key]; !exists {
					scaleSequence[key] = components.SequenceItem{
						Order:  i + 1,
						Finger: marker.Finger,
					}
				}
			}
		}
	}

	fretProps := components.FretboardProps{
		ActiveItems:     activeItems,
		UpcomingMarkers: upcoming,
		ScaleSequence:   scaleSequence,
		Tuning:          m.tuning,
		ScaleConfig:     m.currentLesson.Generator,
		ShowAll:         m.showAll,
		FretCount:       m.fretCount,
		ShowFingers:     m.showFingers,
		ShowScaleShape:  m.showScaleShape,
	}

	// --- 2. RENDER COMPONENTS ---

	// Top Section: Circle + List
	rawCircle := strings.TrimSuffix(components.RenderCircle(m.currentLesson.ActualKey), "\n")
	circleBox := lipgloss.NewStyle().
		Width(circleWidth).
		Height(circleHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(theory.CatOverlay1).
		Render(rawCircle)

	listBox := lipgloss.NewStyle().
		PaddingLeft(1).
		Render(m.list.View())

	topSection := lipgloss.JoinHorizontal(lipgloss.Top, circleBox, listBox)
	topContainer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(theory.CatOverlay1).
		Width(m.width).
		Render(topSection)

	// Bottom Section: Fretboard + Metronome
	fretboardView := components.RenderFretboard(fretProps)

	// Metronome (4 beat pattern)
	metroView := components.RenderMetronome(m.currentStep%4, 4, m.currentLesson.BPM)

	// Status Bar
	playStatus := "Play "
	if m.metronomeActive {
		playStatus = "Pause"
	}
	helpText := fmt.Sprintf("[Space] %s  [H] Fing(%s)  [S] Seq(%s)  [Tab] Note(%s)  [F] Fret(%d)",
		playStatus, status(m.showFingers), status(m.showScaleShape), status(m.showAll), m.fretCount)
	helpView := lipgloss.NewStyle().Foreground(theory.CatSubtext1).Render(helpText)

	infoBar := lipgloss.NewStyle().
		Foreground(theory.CatSky).Bold(true).
		Render(fmt.Sprintf("PLAYING: %s (Step %d/%d)", m.currentLesson.Title, m.currentStep+1, len(steps)))

	bottomSection := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Padding(0, 1).Render(infoBar),
		lipgloss.NewStyle().Render(fretboardView),
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().PaddingLeft(2).Render(metroView),
			lipgloss.NewStyle().PaddingLeft(4).Render(helpView),
		),
	)

	return lipgloss.JoinVertical(lipgloss.Left, topContainer, bottomSection)
}

func tick(bpm int) tea.Cmd {
	return tea.Tick(time.Duration(60000/bpm)*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func status(b bool) string {
	if b {
		return "ON"
	}
	return "off"
}
