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

// MetroBeatMsg is sent when metronome plays a beat
type MetroBeatMsg struct {
	Beat int // Current beat in measure (0-based from metronome)
}

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
	currentBeat   int // Current beat number (1-based)

	// UI State
	list          list.Model
	tuning        []theory.Note
	width, height int
	fretCount     int

	// Display Modes
	showAll        bool // Tab Mode (Note Names) - Phím Tab
	showFingers    bool // Finger Helper Mode - Phím H
	showScaleShape bool // Sequence/Shape Mode - Phím S
	showUpcoming   bool // Toggle upcoming markers - Phím U
	showHelp       bool // Toggle full help text - Phím ?

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
	l.SetShowHelp(false) // Disable built-in help, we'll add custom help
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
		lessons:            loadedLessons,
		currentLesson:      firstLesson,
		list:               l,
		tuning:             theory.StandardTuning,
		fretCount:          12,
		metronomeActive:    false,
		metroPlayer:        metroPlayer,
		metronomeUIMode:    false,
		metroBPM:           120,
		metroTimeSignature: audio.TimeSig4_4,
		metroSoundType:     "wood",

		// Default States
		showAll:        false,
		showFingers:    false,
		showScaleShape: false,
		showUpcoming:   true,
	}
}

// getTotalBeats returns total number of beats in current lesson
func (m Model) getTotalBeats() int {
	if len(m.currentLesson.Steps) == 0 {
		return 0
	}
	// Last step's beat number is the total beats
	lastStep := m.currentLesson.Steps[len(m.currentLesson.Steps)-1]
	return lastStep.Beat
}

// getCurrentStepIndex finds the step index for current beat
// Returns -1 if no step at current beat (hold/rest beat)
func (m Model) getCurrentStepIndex() int {
	for i, step := range m.currentLesson.Steps {
		if step.Beat == m.currentBeat {
			return i
		}
	}
	return -1 // No step at this beat (it's a hold or rest)
}

// getActiveStepIndex finds the most recent step at or before current beat
func (m Model) getActiveStepIndex() int {
	activeIdx := -1
	for i, step := range m.currentLesson.Steps {
		if step.Beat <= m.currentBeat {
			activeIdx = i
		} else {
			break
		}
	}
	return activeIdx
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
				// Start listening for metronome beats
				cmds = append(cmds, listenMetronomeBeat(m.metroPlayer))
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
				m.currentBeat = 1 // Start at beat 1
				// Set BPM from lesson
				if m.currentLesson.BPM > 0 {
					m.metroBPM = m.currentLesson.BPM
					if m.metroPlayer != nil {
						m.metroPlayer.SetBPM(m.metroBPM)
					}
				}
				// Don't auto-start - user will press Space to play
			}

		case "u", "U": // Toggle upcoming markers
			m.showUpcoming = !m.showUpcoming

		case "?": // Toggle help display
			m.showHelp = !m.showHelp
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
			totalBeats := m.getTotalBeats()
			if totalBeats > 0 {
				m.currentBeat = (m.currentBeat % totalBeats) + 1 // Loop through beats
			}
			// Continue listening for metronome beats
			cmds = append(cmds, listenMetronomeBeat(m.metroPlayer))
		}

	case MetroBeatMsg:
		// Sync UI beat with metronome beat
		if m.metronomeActive && len(m.currentLesson.Steps) > 0 {
			totalBeats := m.getTotalBeats()
			if totalBeats > 0 {
				m.currentBeat = (m.currentBeat % totalBeats) + 1 // Loop through beats
			}
			// Continue listening for next beat
			cmds = append(cmds, listenMetronomeBeat(m.metroPlayer))
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

	// --- 1. PREPARE FRETBOARD PROPS ---

	// Build fretboard data using optimized builder
	builder := NewFretboardDataBuilder(&m.currentLesson, m.currentBeat)

	var activeItems []components.ActiveItem
	var upcoming map[string]components.UpcomingItem
	var scaleSequence map[string]components.SequenceItem

	// Only build what we need based on display modes
	activeItems = builder.BuildActiveItems()

	if m.showUpcoming {
		upcoming = builder.BuildUpcomingMarkers(3)
	} else {
		upcoming = make(map[string]components.UpcomingItem)
	}

	if m.showScaleShape || m.showFingers || m.showAll {
		// Need scale sequence for these modes
		scaleSequence = builder.BuildScaleSequence()
	} else {
		scaleSequence = make(map[string]components.SequenceItem)
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

	// Help text below list
	var helpText string
	if m.showHelp {
		// Full help with current status - split into multiple lines
		playStatus := "Play"
		if m.metronomeActive {
			playStatus = "Pause"
		}
		line1 := fmt.Sprintf("[Space] %s  [M] Metro  [H] Fing(%s)  [S] Seq(%s)",
			playStatus, status(m.showFingers), status(m.showScaleShape))
		line2 := fmt.Sprintf("[Tab] Note(%s)  [U] Upc(%s)  [F] Fret(%d)  [?] less",
			status(m.showAll), status(m.showUpcoming), m.fretCount)
		helpText = line1 + "\n" + line2
	} else {
		// Short help
		helpText = "↑/k up • ↓/j down • q quit • ? more"
	}
	helpView := lipgloss.NewStyle().
		Foreground(theory.CatSubtext1).
		PaddingLeft(1).
		Render(helpText)

	listWithHelp := lipgloss.JoinVertical(lipgloss.Left, listBox, helpView)

	topSection := lipgloss.JoinHorizontal(lipgloss.Top, circleBox, listWithHelp)
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

	// Info Bar
	infoBar := lipgloss.NewStyle().
		Foreground(theory.CatSky).Bold(true).
		Render(fmt.Sprintf("PLAYING: %s (Beat %d/%d)", m.currentLesson.Title, m.currentBeat, m.getTotalBeats()))

	// Build bottom section (fretboard + metronome bar)
	bottomSection := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Padding(0, 1).Render(infoBar),
		lipgloss.NewStyle().Render(fretboardView),
		lipgloss.NewStyle().PaddingLeft(2).Render(metroDisplay),
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

// listenMetronomeBeat creates a command that waits for the next metronome beat
func listenMetronomeBeat(player *audio.MetronomePlayer) tea.Cmd {
	if player == nil {
		return nil
	}
	return func() tea.Msg {
		beat := <-player.OnBeatChannel()
		return MetroBeatMsg{Beat: beat}
	}
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
