# BPM Adjustment Fix

## Vấn đề

Khi user nhấn `+/-` để thay đổi BPM, giá trị BPM được update nhưng metronome vẫn chạy với tốc độ cũ.

## Nguyên nhân

```go
func (m *MetronomePlayer) run() {
    ticker := time.NewTicker(m.beatDuration) // Tạo ticker 1 lần
    defer ticker.Stop()
    
    for {
        case <-ticker.C:  // Ticker này không bao giờ thay đổi!
            // Play beat
    }
}
```

**Vấn đề**: `time.Ticker` được tạo 1 lần trong `run()` với `beatDuration` ban đầu. Khi `SetBPM()` update `beatDuration`, ticker cũ vẫn tiếp tục chạy với interval cũ.

## Giải pháp

### 1. Thêm Reset Channel

```go
type MetronomePlayer struct {
    // ...
    resetChan chan struct{} // Signal để reset ticker
}
```

### 2. Update run() để lắng nghe reset signal

```go
func (m *MetronomePlayer) run() {
    ticker := time.NewTicker(m.beatDuration)
    defer ticker.Stop()

    for {
        select {
        case <-m.stopChan:
            return
        case <-m.resetChan:  // Nhận reset signal
            ticker.Stop()
            m.mu.RLock()
            newDuration := m.beatDuration
            m.mu.RUnlock()
            ticker = time.NewTicker(newDuration)  // Tạo ticker mới
        case <-ticker.C:
            // Play beat
        }
    }
}
```

### 3. SetBPM() gửi reset signal

```go
func (m *MetronomePlayer) SetBPM(bpm int) {
    m.mu.Lock()
    m.config.BPM = bpm
    m.beatDuration = time.Minute / time.Duration(bpm)
    m.mu.Unlock()
    
    // Trigger reset ticker
    select {
    case m.resetChan <- struct{}{}:
    default:
        // Non-blocking
    }
}
```

## Cách hoạt động

```
User nhấn [+]
    ↓
SetBPM(125) được gọi
    ↓
Update beatDuration = 60s/125 = 480ms
    ↓
Gửi signal qua resetChan
    ↓
run() nhận signal
    ↓
Stop ticker cũ
    ↓
Tạo ticker mới với duration 480ms
    ↓
Metronome chạy với BPM mới!
```

## Test

1. Run app: `./guitui`
2. Select lesson, nhấn `Enter`
3. Nhấn `M` để mở metronome settings
4. Nhấn `Space` để play
5. Nhấn `+` nhiều lần - BPM tăng, tempo nhanh hơn ✅
6. Nhấn `-` nhiều lần - BPM giảm, tempo chậm lại ✅

## Technical Notes

- **Thread-safe**: Sử dụng `sync.RWMutex` để protect `beatDuration`
- **Non-blocking**: Reset signal dùng `select` với `default` để tránh blocking
- **Immediate effect**: Ticker được reset ngay lập tức khi BPM thay đổi
- **No race conditions**: Lock/unlock đúng cách để tránh race

## Files Modified

- `internal/audio/metronome.go`:
  - Added `resetChan` field
  - Updated `run()` with reset case
  - Updated `SetBPM()` to send reset signal
