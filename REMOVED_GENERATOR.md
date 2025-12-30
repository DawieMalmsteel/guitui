# Removed Auto-Generator - Manual Lessons Only

## Vấn đề

Auto-generator **SAI HOÀN TOÀN**:
- Chỉ đúng với scale đầu tiên  
- Các scales khác thiếu nốt hoặc sai vị trí
- Logic quá phức tạp, không chuẩn guitar thực tế

## Giải pháp: MANUAL LESSONS

Loại bỏ hoàn toàn auto-generation, chỉ dùng **lessons được viết tay trong JSON**.

## Changes Made

### 1. ❌ Removed Files/Code
- `internal/lesson/generator.go` - DELETED
- `GeneratorConfig` struct - REMOVED
- Position switching (N key) - REMOVED
- Position type toggle (P key) - REMOVED
- All generator-based rendering logic - REMOVED

### 2. ✅ Simplified Model

```go
// Before
type Lesson struct {
    Generator *GeneratorConfig
    Steps     []Step `json:"-"`  // Generated
}

// After
type Lesson struct {
    Steps []Step `json:"steps"`  // Load trực tiếp từ JSON
}
```

### 3. ✅ Manual JSON Format

```json
{
  "title": "A Minor Pentatonic - Box 1",
  "category": "scale",
  "bpm": 80,
  "key": "A",
  "steps": [
    {"beat": 1, "markers": [{"string": 0, "fret": 5, "finger": 1}]},
    {"beat": 2, "markers": [{"string": 0, "fret": 8, "finger": 4}]},
    ...
  ]
}
```

**String index**: 0 = String 6 (low E), 5 = String 1 (high E)

### 4. ✅ Simplified UI

**Removed keys:**
- ❌ `N` - Switch position
- ❌ `P` - Toggle position type

**Kept keys:**
- ✅ `Space` - Play/Pause
- ✅ `M` - Metronome settings
- ✅ `S` - Scale shape (sequence numbers)
- ✅ `H` - Finger numbers
- ✅ `Tab` - Note names
- ✅ `U` - Upcoming markers
- ✅ `F` - Toggle fret count

### 5. ✅ Simplified Rendering

**Before:**
- Generator calculates positions
- Complex position-based rendering
- Multiple code paths

**After:**
- Load steps from JSON
- Render directly from steps
- Single simple code path

## How to Create Lessons

### Example: A Minor Pentatonic Box 1

```json
{
  "title": "A Minor Pentatonic - Box 1 (Fret 5-8)",
  "category": "scale",
  "bpm": 80,
  "key": "A",
  "steps": [
    // String 6 (index 0): Low E string
    {"beat": 1, "markers": [{"string": 0, "fret": 5, "finger": 1}]},  // A
    {"beat": 2, "markers": [{"string": 0, "fret": 8, "finger": 4}]},  // C
    
    // String 5 (index 1): A string
    {"beat": 3, "markers": [{"string": 1, "fret": 5, "finger": 1}]},  // D
    {"beat": 4, "markers": [{"string": 1, "fret": 7, "finger": 3}]},  // E
    
    // ... continue for all strings
  ]
}
```

### Guitar String Reference

| String Index (JSON) | String Number | Open Note |
|---------------------|---------------|-----------|
| 0 | 6 (thấp nhất) | E |
| 1 | 5 | A |
| 2 | 4 | D |
| 3 | 3 | G |
| 4 | 2 | B |
| 5 | 1 (cao nhất) | E |

## Benefits

### ✅ Advantages
1. **100% Accurate**: Lessons được viết tay, chuẩn guitar thực tế
2. **Simple**: Không có logic phức tạp
3. **Flexible**: Có thể tạo bất kỳ pattern nào
4. **Maintainable**: Dễ fix bugs, dễ hiểu code
5. **Reliable**: Không có surprise, không có auto-generation sai

### ❌ Trade-offs
1. Phải viết lessons thủ công
2. Không có auto-switching positions
3. Cần tạo nhiều JSON files hơn

## Migration Guide

### Old lessons.json (with generator)
```json
{
  "generator": {
    "root": "A",
    "scale": "minor_pentatonic",
    "position": 1
  }
}
```

### New lessons.json (manual steps)
```json
{
  "steps": [
    {"beat": 1, "markers": [{"string": 0, "fret": 5, "finger": 1}]},
    ...
  ]
}
```

## Files Modified

1. `internal/lesson/model.go` - Simplified Lesson struct
2. `internal/lesson/loader.go` - Calculate notes from markers
3. `internal/ui/model.go` - Removed position state, N/P keys
4. `internal/ui/components/fretboard.go` - Simplified rendering

## Result

**Simpler, more reliable, 100% accurate guitar lessons!** 🎸✅

Bây giờ bạn có thể tạo lessons chính xác bằng tay trong JSON.
