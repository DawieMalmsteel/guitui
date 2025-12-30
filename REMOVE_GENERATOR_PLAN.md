# Plan: Remove Auto-Generated Lessons

## Vấn đề

Generator tự động sinh lessons **SAI HOÀN TOÀN**:
- Chỉ đúng với scale đầu tiên
- Các scales khác thiếu nốt hoặc sai vị trí
- Logic generator quá phức tạp và không chuẩn guitar

## Giải pháp

**LOẠI BỎ hoàn toàn auto-generation, chỉ dùng manual JSON lessons**

## Changes Required

### 1. Remove Generator Code
- ❌ Delete `internal/lesson/generator.go`
- ❌ Remove `GeneratorConfig` từ model
- ✅ Keep `Lesson` struct nhưng chỉ load từ JSON

### 2. Simplify Lesson Model
```go
type Lesson struct {
    Title    string  `json:"title"`
    Category string  `json:"category"`
    BPM      int     `json:"bpm"`
    KeyStr   string  `json:"key"`
    Steps    []Step  `json:"steps"`  // Load trực tiếp từ JSON
}
```

### 3. Update lessons.json
Tạo lessons với steps thủ công, viết tay đúng guitar scale.

### 4. Remove UI Features
- ❌ Remove "N" key (switch position)
- ❌ Remove "P" key (toggle position type)
- ✅ Keep "S", "H", "Tab", "U" keys (display modes)

### 5. Simplify Rendering
- Remove position-based rendering
- Chỉ render notes từ Steps array

Bạn đồng ý với plan này không?
