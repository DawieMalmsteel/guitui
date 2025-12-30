# Metronome Update - Migration Summary

## Packages Updated

Đã migrate tất cả packages lên version mới nhất:

- **Go version**: 1.21.0 → 1.25.5
- **charmbracelet/bubbles**: v0.21.0 (added)
- **charmbracelet/colorprofile**: v0.2.3 → v0.4.1
- **charmbracelet/x/ansi**: v0.10.1 → v0.11.3
- **charmbracelet/x/cellbuf**: v0.0.13 → v0.0.14
- **charmbracelet/x/term**: v0.2.1 → v0.2.2
- **ebitengine/oto/v3**: v3.3.2 → v3.4.0
- **ebitengine/purego**: v0.8.0 → v0.9.1
- **gopxl/beep/v2**: v2.1.0 → v2.1.1
- **lucasb-eyer/go-colorful**: v1.2.0 → v1.3.0
- **mattn/go-runewidth**: v0.0.16 → v0.0.19
- **golang.org/x/sys**: v0.36.0 → v0.39.0
- **golang.org/x/text**: v0.3.8 → v0.32.0

## Metronome Features Implemented

### 1. Audio Playback
- Sử dụng gopxl/beep/v2 API mới với `generators.SineTone()`
- Speaker được init một lần duy nhất với `sync.Once`
- Metronome phát âm thanh thực sự khi active
- Accent beat (phách đầu) có âm cao hơn

### 2. Interactive UI
- **[M]** - Toggle metronome settings panel
- **[Space]** - Play/Pause metronome
- **[+/-]** - Tăng/giảm BPM (5 bpm mỗi lần, range: 40-240)
- **[1]** - Set time signature 4/4
- **[2]** - Set time signature 3/4
- **[3]** - Set time signature 6/8
- **[4]** - Set time signature 2/4

### 3. Display
- Metronome bar hiển thị beat hiện tại với màu sắc:
  - Đỏ: Down beat (phách đầu)
  - Vàng: Up beat (phách thường)
  - Xám: Inactive beats
- Settings panel hiển thị:
  - BPM hiện tại
  - Time signature hiện tại
  - Play/Pause status
  - Keyboard shortcuts

## Files Modified

### 1. internal/audio/metronome.go
- Viết lại hoàn toàn để sử dụng gopxl/beep/v2 API
- Sử dụng `sync.Once` cho speaker initialization
- Thêm volume control với `volumeStreamer`
- Thread-safe với `sync.RWMutex`
- Metronome beat tracking với `currentBeat`

### 2. internal/ui/model.go
- Thêm fields: `metronomeUIMode`, `metroBPM`, `metroTimeSignature`
- Tích hợp metronome player controls
- Keyboard bindings cho metronome settings
- Update tick logic để sync với metronome

### 3. internal/ui/components/metronome.go
- Thêm `RenderMetronomeSettings()` function
- UI styles cho metronome panel
- Display time signature và controls

## Usage

1. **Start/Stop Metronome**: Press `[Space]`
2. **Open Settings**: Press `[M]`
3. **Adjust BPM**: Press `[+]` or `[-]` (khi settings mở)
4. **Change Time Signature**: Press `[1]`, `[2]`, `[3]`, hoặc `[4]` (khi settings mở)

## Technical Details

### Audio Implementation
```go
// Tạo sine wave tone
tone, _ := generators.SineTone(sampleRate, frequency)

// Giới hạn duration
limited := beep.Take(duration, tone)

// Volume control
volumeStreamer{Streamer: limited, Volume: 0.8}

// Play through speaker
speaker.Play(sound)
```

### Speaker Initialization
```go
var (
    speakerInitOnce sync.Once
    speakerInitErr  error
)

speakerInitOnce.Do(func() {
    speakerInitErr = speaker.Init(sampleRate, bufferSize)
})
```

## Testing

Build thành công:
```bash
go build ./...
```

All packages updated và code compile without errors.
