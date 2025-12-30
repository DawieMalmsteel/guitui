# Fix: S Key Render After N Key

## Vấn đề mô tả

Khi nhấn `S` để bật Scale Shape mode, sau đó nhấn `N` để switch position, Scale Shape không render lại đúng.

## Root Cause Analysis

### Code Flow

1. **Nhấn S**: `showScaleShape = true`
2. **View() renders**: Build `scaleSequence` từ `m.currentLesson.Steps`
3. **Nhấn N**: 
   - Update `currentPosition`
   - Regenerate `steps` ✅
   - Reset `currentStep` ✅
   - Update `availablePositions` ✅ (đã thêm)
4. **View() renders again**: Rebuild `scaleSequence` từ steps MỚI ✅

### Logic trong View()

```go
scaleSequence := make(map[string]components.SequenceItem)
if m.showScaleShape && len(steps) > 0 {
    for i, step := range steps {
        for _, marker := range step.Markers {
            key := fmt.Sprintf("%d_%d", marker.StringIndex, marker.Fret)
            if _, exists := scaleSequence[key]; !exists {
                scaleSequence[key] = components.SequenceItem{
                    Order:  i + 1,
                    Finger: marker.Finger,
                }
            }
        }
    }
}
```

**Điều này ĐÚNG!** ScaleSequence được rebuild mỗi frame từ `m.currentLesson.Steps`.

## Possible Issues

### 1. ✅ Steps không regenerate
**Fix:** Đã có rồi ở line 325-329

### 2. ✅ availablePositions không update
**Fix:** Đã thêm `m.availablePositions = maxPos`

### 3. Có thể: ScaleSequence map key conflict
Nếu Position 1 và Position 2 có cùng notes ở cùng vị trí (string, fret), map sẽ chỉ lưu lần xuất hiện đầu tiên.

**Nhưng** map được rebuild HOÀN TOÀN mỗi frame, nên không có vấn đề này.

## Expected Behavior

### Scenario 1: S → N
```
1. Nhấn S → showScaleShape = true
2. View renders → ScaleSequence built từ Position 1
3. Nhấn N → Position changes to 2, steps regenerated
4. View renders → ScaleSequence rebuilt từ Position 2 ✅
```

### Scenario 2: N → S  
```
1. Nhấn N → Position changes, steps regenerated
2. Nhấn S → showScaleShape = true
3. View renders → ScaleSequence built từ current position ✅
```

### Scenario 3: S → N → N → S → N
```
All transitions should work correctly ✅
```

## Actual Fix Applied

```go
case "n", "N": // Switch position
    if m.currentLesson.Generator != nil {
        maxPos := theory.GetPositionCount(...)
        if maxPos > 0 {
            m.currentPosition = (m.currentPosition % maxPos) + 1
            m.currentLesson.Generator.Position = m.currentPosition
            m.availablePositions = maxPos  // ← THÊM DÒNG NÀY
            
            // Regenerate steps with new position
            steps, err := lesson.GenerateSteps(m.currentLesson.Generator)
            if err == nil {
                m.currentLesson.Steps = steps
                m.currentStep = 0
            }
        }
    }
```

## Testing Steps

1. **Start app**: `./guitui`
2. **Select lesson**: Arrow keys + Enter
3. **Test S key**:
   - Nhấn `S` → Scale shape numbers hiện
   - Nhấn `S` lại → Scale shape tắt
4. **Test S → N**:
   - Nhấn `S` → Scale shape ON
   - Nhấn `N` → Position switches
   - ✅ Scale shape vẫn ON với numbers từ position mới
5. **Test N → S**:
   - Nhấn `N` nhiều lần → Switch positions
   - Nhấn `S` → Scale shape hiện với current position
   - ✅ Numbers đúng cho position hiện tại

## Additional Notes

### Display Modes Persist Across Position Changes
- `showScaleShape` (S key) - PERSIST
- `showFingers` (H key) - PERSIST  
- `showAll` (Tab key) - PERSIST
- `showUpcoming` (U key) - PERSIST

**Đây là behavior mong muốn!** User có thể giữ display mode khi explore các positions khác nhau.

### If Still Having Issues

Có thể là do:
1. **Pattern type conflict**: Nếu lesson có `pattern: "sequence"` thay vì `"box"`, steps sẽ khác
2. **Scale không có position đó**: Check `maxPos > 0`
3. **Generator config bị null**: Check `m.currentLesson.Generator != nil`

Chạy app và test để xác nhận!

---

## UPDATE: Fixed Scale Colors Not Updating

### Vấn đề thực sự

Khi nhấn `S` để hiện màu scale notes, sau đó nhấn `N` để switch position, **màu sắc không update** theo position mới!

### Root Cause

#### Render Logic (fretboard.go)

```go
if props.ShowScaleShape && props.CurrentPosition.Type != "" {
    for s := 0; s < 6; s++ {
        pattern := props.CurrentPosition.NotePatterns[s]
        for i, relFret := range pattern.RelativeFrets {
            fret := props.ScaleConfig.StartFret + relFret  // ← VẤN ĐỀ!
            // Render colored note
        }
    }
}
```

**Vấn đề:** `props.ScaleConfig.StartFret` KHÔNG được update khi switch position!

#### Khi switch position (N key):

```go
// CHỈ update Position, KHÔNG update StartFret/EndFret
m.currentLesson.Generator.Position = m.currentPosition
```

**Kết quả:** 
- Position mới ✅
- Steps regenerated ✅  
- **StartFret/EndFret vẫn giữ giá trị cũ** ❌
- Màu render sai vị trí ❌

### Fix Applied

Khi nhấn N, **recalculate StartFret và EndFret** cho position mới:

```go
case "n", "N": // Switch position
    m.currentPosition = (m.currentPosition % maxPos) + 1
    m.currentLesson.Generator.Position = m.currentPosition
    
    // ← THÊM LOGIC NÀY
    if pos, exists := theory.GetPosition(..., m.currentPosition); exists {
        root := parseNote(m.currentLesson.Generator.Root)
        rootFret := theory.FindRootFretOn6thString(root)
        startFret, endFret := theory.CalculateFretRange(pos, rootFret)
        
        // Update config với fret range mới
        m.currentLesson.Generator.StartFret = startFret
        m.currentLesson.Generator.EndFret = endFret
    }
    
    // Regenerate steps
    steps, err := lesson.GenerateSteps(m.currentLesson.Generator)
    ...
```

### Ví dụ

**A Minor Pentatonic:**

Position 1:
- StartFret = 5 (root at fret 5)
- Colors at frets 5, 7, 8

Position 2:
- StartFret = 8 (root + offset 3)
- Colors at frets 8, 10, 12

**Trước fix:**
1. Chọn Position 1 → Màu ở frets 5, 7, 8 ✅
2. Nhấn N → Position 2 nhưng màu VẪN ở 5, 7, 8 ❌

**Sau fix:**
1. Chọn Position 1 → Màu ở frets 5, 7, 8 ✅
2. Nhấn N → Position 2, màu UPDATE đến 8, 10, 12 ✅

### Testing

```
1. Start app
2. Select A Minor Pentatonic lesson
3. Nhấn S → Scale colors hiện (position 1)
4. Nhấn N → Colors MOVE to position 2 ✅
5. Nhấn N lại → Colors MOVE to position 3 ✅
6. Nhấn S → Colors tắt
7. Nhấn S lại → Colors hiện lại đúng position hiện tại ✅
```

### Files Changed

- `internal/ui/model.go`:
  - Case "n", "N": Added fret range recalculation
  - Added `parseNote()` helper function

**Bây giờ scale colors update ĐÚNG khi switch position!** 🎨✅
