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
	showAll         bool
	metronomeActive bool // Bấm Space để chạy/dừng
}

func NewModel() Model {
	// 1. Load Data
	loadedLessons, err := lesson.LoadLessons("lessons.json")
	if err != nil {
		fmt.Println("Lỗi load lessons.json (Tạo file chưa tml?):", err)
		loadedLessons = []lesson.Lesson{} // Empty fallback
	}

	// 2. Setup List Component
	var items []list.Item
	for _, l := range loadedLessons {
		items = append(items, item{lesson: l})
	}

	// Custom Delegate để hiển thị list đẹp kiểu Catppuccin
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

	// 3. Default Lesson (Bài đầu tiên)
	firstLesson := lesson.Lesson{}
	if len(loadedLessons) > 0 {
		firstLesson = loadedLessons[0]
		// Cần generate steps cho bài đầu tiên ngay
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
		metronomeActive: false, // Mặc định dừng
	}
}

func (m Model) Init() tea.Cmd {
	return nil // Chưa chạy metronome vội
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
			m.fretCount = 36 - m.fretCount // Toggle 12 <-> 24 (Trick: 12->24, 24->12)
			if m.fretCount != 24 {
				m.fretCount = 12
			} // Safety

		case "tab":
			m.showAll = !m.showAll

		case " ": // Space: Toggle Play/Pause
			m.metronomeActive = !m.metronomeActive
			if m.metronomeActive {
				cmds = append(cmds, tick(m.currentLesson.BPM))
			}

		case "enter": // Chọn bài từ list
			if selectedItem, ok := m.list.SelectedItem().(item); ok {
				m.currentLesson = selectedItem.lesson

				// Generate Steps nếu chưa có
				if len(m.currentLesson.Steps) == 0 && m.currentLesson.Generator != nil {
					steps, err := lesson.GenerateSteps(m.currentLesson.Generator)
					if err == nil {
						m.currentLesson.Steps = steps
					}
				}

				// Reset trạng thái
				m.currentStep = 0

				// --- BUG FIX START ---
				// Kiểm tra: Chỉ start tick mới nếu trước đó nó ĐANG DỪNG.
				// Nếu đang chạy rồi thì thôi, để cái loop cũ nó tự handle bài mới.
				if !m.metronomeActive {
					m.metronomeActive = true
					cmds = append(cmds, tick(m.currentLesson.BPM))
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Resize List: Nằm bên phải Circle
		listWidth := msg.Width - circleWidth - 4
		if listWidth < 20 {
			listWidth = 20
		}
		listHeight := circleHeight - listHeaderHeight // Khớp chiều cao với Circle

		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)

	case TickMsg:
		if m.metronomeActive && len(m.currentLesson.Steps) > 0 {
			m.currentStep = (m.currentStep + 1) % len(m.currentLesson.Steps)
			cmds = append(cmds, tick(m.currentLesson.BPM))
		}
	}

	// Update List Component (để nó xử lý phím lên xuống)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	} // Chờ WindowSizeMsg

	// --- 1. TOP SECTION (Circle | List) ---

	// Circle
	rawCircle := strings.TrimSuffix(components.RenderCircle(m.currentLesson.ActualKey), "\n")
	circleBox := lipgloss.NewStyle().
		Width(circleWidth).
		Height(circleHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.NormalBorder(), false, true, false, false). // Border phải
		BorderForeground(theory.CatOverlay1).
		Render(rawCircle)

	// List
	listBox := lipgloss.NewStyle().
		PaddingLeft(1).
		Render(m.list.View())

	topSection := lipgloss.JoinHorizontal(lipgloss.Top, circleBox, listBox)

	// Container Top
	topContainer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false). // Border dưới
		BorderForeground(theory.CatOverlay1).
		Width(m.width).
		Render(topSection)

	// --- 2. BOTTOM SECTION (Fretboard + Metronome) ---

	// Prepare Props
	var activeMarkers []lesson.Marker
	if len(m.currentLesson.Steps) > 0 {
		activeMarkers = m.currentLesson.Steps[m.currentStep].Markers
	}

	// Upcoming Logic
	upcoming := make(map[string]int)
	lookAhead := 3
	if len(m.currentLesson.Steps) > 0 {
		for i := 1; i <= lookAhead; i++ {
			nextIdx := (m.currentStep + i) % len(m.currentLesson.Steps)
			for _, marker := range m.currentLesson.Steps[nextIdx].Markers {
				key := fmt.Sprintf("%d_%d", marker.StringIndex, marker.Fret)
				if _, exists := upcoming[key]; !exists {
					upcoming[key] = i
				}
			}
		}
	}

	fretProps := components.FretboardProps{
		ActiveMarkers:   activeMarkers,
		UpcomingMarkers: upcoming,
		Tuning:          m.tuning,
		ScaleConfig:     m.currentLesson.Generator,
		ShowAll:         m.showAll,
		FretCount:       m.fretCount,
	}

	// Render Components
	fretboardView := components.RenderFretboard(fretProps)

	// Tính total beats để vẽ metronome (Giả sử 4/4 hoặc đếm theo steps)
	// Để đơn giản, vẽ 4 cục đại diện cho nhịp điệu, active cái (currentStep % 4)
	metroView := components.RenderMetronome(m.currentStep%4, 4, m.currentLesson.BPM)

	playStatus := "Play "
	if m.metronomeActive {
		playStatus = "Pause"
	}

	helpText := fmt.Sprintf("[Space] %s  [F] Fret(%d)  [Tab] Scale  [Enter] Select",
		playStatus, m.fretCount)

	helpView := lipgloss.NewStyle().Foreground(theory.CatSubtext1).Render(helpText)

	// Info Bar
	infoBar := lipgloss.NewStyle().
		Foreground(theory.CatSky).Bold(true).
		Render(fmt.Sprintf("PLAYING: %s", m.currentLesson.Title))

	bottomSection := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Padding(0, 1).Render(infoBar),
		lipgloss.NewStyle().Render(fretboardView), // Sát border trên
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
