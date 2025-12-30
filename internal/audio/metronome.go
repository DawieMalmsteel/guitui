package audio

import (
	"fmt"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/speaker"
)

var (
	speakerInitOnce sync.Once
	speakerInitErr  error
)

type MetronomeConfig struct {
	BPM           int
	TimeSignature TimeSignature
	AccentFirst   bool
	Volume        int
	SoundType     string
}

type TimeSignature int

const (
	TimeSig4_4 TimeSignature = 4
	TimeSig3_4 TimeSignature = 3
	TimeSig6_8 TimeSignature = 6
	TimeSig2_4 TimeSignature = 2
)

type MetronomePlayer struct {
	config       *MetronomeConfig
	sampleRate   beep.SampleRate
	beatDuration time.Duration
	isPlaying    bool
	mu           sync.RWMutex
	currentBeat  int
	stopChan     chan struct{}
	resetChan    chan struct{} // Signal to reset ticker
}

func NewMetronomePlayer(config *MetronomeConfig) (*MetronomePlayer, error) {
	sampleRate := beep.SampleRate(44100)
	
	// Initialize speaker only once
	speakerInitOnce.Do(func() {
		speakerInitErr = speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	})
	
	if speakerInitErr != nil {
		return nil, fmt.Errorf("khởi tạo speaker: %w", speakerInitErr)
	}

	metronome := &MetronomePlayer{
		config:       config,
		sampleRate:   sampleRate,
		beatDuration: time.Minute / time.Duration(config.BPM),
		isPlaying:    false,
		currentBeat:  0,
		stopChan:     make(chan struct{}),
		resetChan:    make(chan struct{}),
	}

	go metronome.run()

	return metronome, nil
}

func (m *MetronomePlayer) createSound(isAccent bool) beep.Streamer {
	switch m.config.SoundType {
	case "mechanical":
		return m.createMechanicalClick(isAccent)
	case "digital":
		return m.createDigitalBeep(isAccent)
	case "wood":
		return m.createWoodBlock(isAccent)
	default: // "click" or anything else
		return m.createWoodBlock(isAccent)
	}
}

// createWoodBlock creates a short, percussive wood block sound
func (m *MetronomePlayer) createWoodBlock(isAccent bool) beep.Streamer {
	var baseFreq, harmonicFreq float64
	var attackDuration, decayDuration time.Duration
	
	if isAccent {
		baseFreq = 1800.0
		harmonicFreq = 3600.0
		attackDuration = time.Millisecond * 2
		decayDuration = time.Millisecond * 40
	} else {
		baseFreq = 1400.0
		harmonicFreq = 2800.0
		attackDuration = time.Millisecond * 1
		decayDuration = time.Millisecond * 30
	}

	baseTone, err1 := generators.SineTone(m.sampleRate, baseFreq)
	harmonicTone, err2 := generators.SineTone(m.sampleRate, harmonicFreq)
	
	if err1 != nil || err2 != nil {
		return beep.Silence(0)
	}

	mixed := beep.Mix(baseTone, &volumeStreamer{Streamer: harmonicTone, Volume: 0.3})

	totalDuration := attackDuration + decayDuration
	enveloped := &envelopeStreamer{
		Streamer:       mixed,
		sampleRate:     m.sampleRate,
		attackSamples:  m.sampleRate.N(attackDuration),
		decaySamples:   m.sampleRate.N(decayDuration),
		currentSample:  0,
	}

	limited := beep.Take(m.sampleRate.N(totalDuration), enveloped)
	volume := float64(m.config.Volume) / 100.0
	
	return &volumeStreamer{Streamer: limited, Volume: volume}
}

// createMechanicalClick creates a sharp mechanical click (shorter, crisper)
func (m *MetronomePlayer) createMechanicalClick(isAccent bool) beep.Streamer {
	var freq float64
	var duration time.Duration
	
	if isAccent {
		freq = 2400.0
		duration = time.Millisecond * 15
	} else {
		freq = 2000.0
		duration = time.Millisecond * 12
	}

	tone, err := generators.SineTone(m.sampleRate, freq)
	if err != nil {
		return beep.Silence(0)
	}

	enveloped := &envelopeStreamer{
		Streamer:       tone,
		sampleRate:     m.sampleRate,
		attackSamples:  m.sampleRate.N(time.Millisecond * 1),
		decaySamples:   m.sampleRate.N(duration),
		currentSample:  0,
	}

	limited := beep.Take(m.sampleRate.N(duration), enveloped)
	volume := float64(m.config.Volume) / 100.0
	
	return &volumeStreamer{Streamer: limited, Volume: volume}
}

// createDigitalBeep creates a clean digital beep sound
func (m *MetronomePlayer) createDigitalBeep(isAccent bool) beep.Streamer {
	var freq float64
	var duration time.Duration
	
	if isAccent {
		freq = 880.0 // A5
		duration = time.Millisecond * 80
	} else {
		freq = 440.0 // A4
		duration = time.Millisecond * 60
	}

	tone, err := generators.SineTone(m.sampleRate, freq)
	if err != nil {
		return beep.Silence(0)
	}

	enveloped := &envelopeStreamer{
		Streamer:       tone,
		sampleRate:     m.sampleRate,
		attackSamples:  m.sampleRate.N(time.Millisecond * 5),
		decaySamples:   m.sampleRate.N(time.Millisecond * 20),
		currentSample:  0,
	}

	limited := beep.Take(m.sampleRate.N(duration), enveloped)
	volume := float64(m.config.Volume) / 100.0 * 0.6 // Digital beep is a bit quieter
	
	return &volumeStreamer{Streamer: limited, Volume: volume}
}

// envelopeStreamer applies an attack-decay envelope to create a percussive sound
type envelopeStreamer struct {
	beep.Streamer
	sampleRate     beep.SampleRate
	attackSamples  int
	decaySamples   int
	currentSample  int
}

func (e *envelopeStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = e.Streamer.Stream(samples)
	
	for i := 0; i < n; i++ {
		var envelope float64
		
		if e.currentSample < e.attackSamples {
			// Attack phase: 0 -> 1
			envelope = float64(e.currentSample) / float64(e.attackSamples)
		} else if e.currentSample < e.attackSamples+e.decaySamples {
			// Decay phase: 1 -> 0
			decayProgress := float64(e.currentSample-e.attackSamples) / float64(e.decaySamples)
			envelope = 1.0 - decayProgress
		} else {
			envelope = 0
		}
		
		samples[i][0] *= envelope
		samples[i][1] *= envelope
		e.currentSample++
	}
	
	return n, ok
}

type volumeStreamer struct {
	beep.Streamer
	Volume float64
}

func (v *volumeStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = v.Streamer.Stream(samples)
	for i := 0; i < n; i++ {
		samples[i][0] *= v.Volume
		samples[i][1] *= v.Volume
	}
	return n, ok
}

func getBeatsPerMeasure(ts TimeSignature) int {
	switch ts {
	case TimeSig6_8:
		return 6
	case TimeSig3_4:
		return 3
	case TimeSig2_4:
		return 2
	case TimeSig4_4:
		return 4
	default:
		return 4
	}
}

func (m *MetronomePlayer) run() {
	ticker := time.NewTicker(m.beatDuration)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-m.resetChan:
			// Reset ticker with new beat duration
			ticker.Stop()
			m.mu.RLock()
			newDuration := m.beatDuration
			m.mu.RUnlock()
			ticker = time.NewTicker(newDuration)
		case <-ticker.C:
			m.mu.RLock()
			if !m.isPlaying {
				m.mu.RUnlock()
				continue
			}
			
			beatsPerMeasure := getBeatsPerMeasure(m.config.TimeSignature)
			isAccent := m.config.AccentFirst && m.currentBeat == 0
			m.mu.RUnlock()

			sound := m.createSound(isAccent)
			speaker.Play(sound)

			m.mu.Lock()
			m.currentBeat = (m.currentBeat + 1) % beatsPerMeasure
			m.mu.Unlock()
		}
	}
}

func (m *MetronomePlayer) Play() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isPlaying = true
}

func (m *MetronomePlayer) Pause() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isPlaying = false
}

func (m *MetronomePlayer) Stop() {
	m.mu.Lock()
	m.isPlaying = false
	m.currentBeat = 0
	m.mu.Unlock()
}

func (m *MetronomePlayer) SetBPM(bpm int) {
	m.mu.Lock()
	m.config.BPM = bpm
	m.beatDuration = time.Minute / time.Duration(bpm)
	m.mu.Unlock()
	
	// Signal to reset ticker with new duration
	select {
	case m.resetChan <- struct{}{}:
	default:
		// Non-blocking, skip if reset already pending
	}
}

func (m *MetronomePlayer) SetTimeSignature(ts TimeSignature) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.TimeSignature = ts
	m.currentBeat = 0
}

func (m *MetronomePlayer) SetAccentFirst(accent bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.AccentFirst = accent
}

func (m *MetronomePlayer) SetVolume(volume int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.Volume = volume
}

func (m *MetronomePlayer) SetSoundType(soundType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.SoundType = soundType
}

func (m *MetronomePlayer) GetCurrentBeat() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentBeat
}

func (m *MetronomePlayer) GetTotalBeats() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return getBeatsPerMeasure(m.config.TimeSignature)
}

func (m *MetronomePlayer) Close() {
	close(m.stopChan)
	speaker.Close()
}
