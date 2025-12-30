package ui

import (
	"fmt"
	"strings"
	"time"

	"guitui/internal/audio"
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
	tuning        []theory.Note
	width, height int
	fretCount     int

	// Display Modes
	showAll        bool // Tab Mode (Note Names) - Phím Tab
	showFingers    bool // Finger Helper Mode - Phím H
	showScaleShape bool // Sequence/Shape Mode - Phím S
	showUpcoming   bool // Toggle upcoming markers - Phím U

	// Metronome State
	metronomeActive    bool
	metroPlayer        *audio.MetronomePlayer
	metronomeUIMode    bool // Toggle metronome settings UI
	metroBPM           int
	metroTimeSignature audio.TimeSignature
	metroSoundType     string
}

func NewModel() Model {
	// 1. Load Data from both JSON and TAB files
	loadedLessons, err := lesson.LoadLessonsFromMultipleSources("lessons.json", "lessons_tab")
	if err != nil {
		fmt.Println("Lỗi load lessons:", err)
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
	// 3. Default Lesson (first lesson from list)
	firstLesson := lesson.Lesson{}
	if len(loadedLessons) > 0 {
		firstLesson = loadedLessons[0]
	}

	metroPlayer, err := audio.NewMetronomePlayer(&audio.MetronomeConfig{
		BPM:           120,
		TimeSignature: audio.TimeSig4_4,
		AccentFirst:   true,
		Volume:        80,
		SoundType:     "wood",
	})
	if err != nil {
		fmt.Println("Lỗi khởi tạo metronome:", err)
	}

	return Model{
		lessons:         loadedLessons,
		currentLesson:   firstLesson,
		list:            l,
		tuning:          theory.StandardTuning,
		fretCount:       12,
		metronomeActive: false,
		metroPlayer:     metroPlayer,
		metronomeUIMode: false,
		metroBPM:        120,
		metroTimeSignature: audio.TimeSig4_4,
		metroSoundType:     "wood",

		// Default States
		showAll:        false,
		showFingers:    false,
		showScaleShape: false,
		showUpcoming:   true,
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
			// Tab mode có thể combine với S (Scale Shape)
			// Nhưng tắt Upcoming và Fingers
			if m.showAll {
				m.showUpcoming = false
				m.showFingers = false
			} else {
				m.showUpcoming = true
			}

		case "h", "H":
			m.showFingers = !m.showFingers
			// Nếu bật Fingers thì tắt Tab mode và Scale Shape
			if m.showFingers {
				m.showAll = false
				m.showScaleShape = false
			}

		case " ": // Space: Toggle Play/Pause
			m.metronomeActive = !m.metronomeActive
			if m.metronomeActive {
				if m.metroPlayer != nil {
					m.metroPlayer.Play()
				}
				cmds = append(cmds, tick(m.metroBPM))
			} else {
				if m.metroPlayer != nil {
					m.metroPlayer.Pause()
				}
			}

		case "m", "M": // Toggle metronome UI settings
			m.metronomeUIMode = !m.metronomeUIMode

		case "+", "=": // Increase BPM
			if m.metronomeUIMode {
				m.metroBPM += 5
				if m.metroBPM > 240 {
					m.metroBPM = 240
				}
				if m.metroPlayer != nil {
					m.metroPlayer.SetBPM(m.metroBPM)
				}
			}

		case "-", "_": // Decrease BPM
			if m.metronomeUIMode {
				m.metroBPM -= 5
				if m.metroBPM < 40 {
					m.metroBPM = 40
				}
				if m.metroPlayer != nil {
					m.metroPlayer.SetBPM(m.metroBPM)
				}
			}

		case "1": // Set 4/4 time signature
			if m.metronomeUIMode {
				m.metroTimeSignature = audio.TimeSig4_4
				if m.metroPlayer != nil {
					m.metroPlayer.SetTimeSignature(audio.TimeSig4_4)
				}
			}

		case "2": // Set 3/4 time signature
			if m.metronomeUIMode {
				m.metroTimeSignature = audio.TimeSig3_4
				if m.metroPlayer != nil {
					m.metroPlayer.SetTimeSignature(audio.TimeSig3_4)
				}
			}

		case "3": // Set 6/8 time signature
			if m.metronomeUIMode {
				m.metroTimeSignature = audio.TimeSig6_8
				if m.metroPlayer != nil {
					m.metroPlayer.SetTimeSignature(audio.TimeSig6_8)
				}
			}

		case "4": // Set 2/4 time signature
			if m.metronomeUIMode {
				m.metroTimeSignature = audio.TimeSig2_4
				if m.metroPlayer != nil {
					m.metroPlayer.SetTimeSignature(audio.TimeSig2_4)
				}
			}

		case "s": // Cycle sound types when in metronome mode
			if m.metronomeUIMode {
				soundTypes := []string{"wood", "mechanical", "digital"}
				currentIdx := 0
				for i, st := range soundTypes {
					if st == m.metroSoundType {
						currentIdx = i
						break
					}
				}
				nextIdx := (currentIdx + 1) % len(soundTypes)
				m.metroSoundType = soundTypes[nextIdx]
				if m.metroPlayer != nil {
					m.metroPlayer.SetSoundType(m.metroSoundType)
				}
			} else {
				// S key behavior for scale shape - CAN combine with Tab
				m.showScaleShape = !m.showScaleShape
				if m.showScaleShape {
					// Tắt Fingers và Upcoming, nhưng GIỮ Tab
					m.showFingers = false
					m.showUpcoming = false
				}
			}

		case "S": // Keep uppercase S for scale shape always
			m.showScaleShape = !m.showScaleShape
			if m.showScaleShape {
				// Tắt Fingers và Upcoming, nhưng GIỮ Tab
				m.showFingers = false
				m.showUpcoming = false
			}

		case "enter": // Chọn bài
			if selectedItem, ok := m.list.SelectedItem().(item); ok {
				m.currentLesson = selectedItem.lesson
				m.currentStep = 0

				if !m.metronomeActive {
					m.metronomeActive = true
					cmds = append(cmds, tick(m.metroBPM))
				}
			}

		case "u", "U": // Toggle upcoming markers
			m.showUpcoming = !m.showUpcoming
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
			cmds = append(cmds, tick(m.metroBPM))
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
	if m.showUpcoming && len(steps) > 0 {
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

	// Metronome display
	var metroDisplay string
	totalBeats := 4
	switch m.metroTimeSignature {
	case audio.TimeSig3_4:
		totalBeats = 3
	case audio.TimeSig6_8:
		totalBeats = 6
	case audio.TimeSig2_4:
		totalBeats = 2
	}
	
	if m.metronomeUIMode {
		// Show metronome settings UI
		metroDisplay = components.RenderMetronomeSettings(m.metroBPM, m.metroTimeSignature, m.metroSoundType, m.metronomeActive)
	} else {
		// Show simple metronome bar
		currentBeat := 0
		if m.metroPlayer != nil {
			currentBeat = m.metroPlayer.GetCurrentBeat()
		}
		metroDisplay = components.RenderMetronome(currentBeat, totalBeats, m.metroBPM)
	}

	// Status Bar
	playStatus := "Play "
	if m.metronomeActive {
		playStatus = "Pause"
	}
	helpText := fmt.Sprintf("[Space] %s  [M] Metro  [H] Fing(%s)  [S] Seq(%s)  [Tab] Note(%s)  [U] Upc(%s)  [F] Fret(%d)",
		playStatus, status(m.showFingers), status(m.showScaleShape), status(m.showAll), status(m.showUpcoming), m.fretCount)
	helpView := lipgloss.NewStyle().Foreground(theory.CatSubtext1).Render(helpText)

	infoBar := lipgloss.NewStyle().
		Foreground(theory.CatSky).Bold(true).
		Render(fmt.Sprintf("PLAYING: %s (Step %d/%d)", m.currentLesson.Title, m.currentStep+1, len(steps)))

	// Build bottom section (fretboard + metronome bar)
	bottomSection := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Padding(0, 1).Render(infoBar),
		lipgloss.NewStyle().Render(fretboardView),
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().PaddingLeft(2).Render(metroDisplay),
			lipgloss.NewStyle().PaddingLeft(4).Render(helpView),
		),
	)

	mainView := lipgloss.JoinVertical(lipgloss.Left, topContainer, bottomSection)

	// If metronome settings mode, overlay the settings panel centered
	if m.metronomeUIMode {
		// Center the settings panel on screen
		return lipgloss.Place(m.width, m.height, 
			lipgloss.Center, lipgloss.Center,
			metroDisplay)
	}

	return mainView
}

func tick(bpm int) tea.Cmd {
	return tea.Tick(time.Duration(60000/bpm)*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
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

func status(b bool) string {
	if b {
		return "ON"
	}
	return "off"
}
