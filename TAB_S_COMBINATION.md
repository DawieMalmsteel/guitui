# Tab + S Combination Mode

## 🎯 Feature: Tab + Scale Shape Combination

**Behavior:** Khi bật cả Tab và S cùng lúc, hiển thị:
- Tên nốt nhạc (A, C, D...) trên **toàn bộ fretboard**
- Notes trong scale có **màu nền** theo ngón tay đã quy định
- **Chữ đen** để dễ đọc trên nền màu

## 🎨 Visual Example

### Tab Only (Tab ON, S OFF)
```
Tất cả notes với màu chữ theo pitch:
E |--E--F--F#-G--G#-A--A#-B--C--C#-D--D#-E--|
   (màu chữ: E=yellow, F=green, F#=teal...)
```

### S Only (S ON, Tab OFF)
```
Chỉ notes trong scale với sequence numbers:
E |---1---------4--------------------------|
   (background màu theo ngón, số thứ tự)
```

### Tab + S (BOTH ON) ✨
```
Tất cả notes, nhưng scale notes có background màu:
E |--E--F--F#-G--G#-[A]--A#-B--[C]--C#-D--D#-E--|
   
Notes không có trong scale: Màu chữ theo pitch
Notes có trong scale: Chữ ĐEN trên background màu ngón tay
  - [A] - Chữ đen trên nền BLUE (ngón 1)
  - [C] - Chữ đen trên nền RED (ngón 4)
```

## 🎹 Keyboard Controls

### Bật Tab + S

**Option 1:**
1. Nhấn `Tab` → Tab mode ON
2. Nhấn `S` → Scale shape ON
3. **Result:** Tab + S combination mode

**Option 2:**
1. Nhấn `S` → Scale shape ON
2. Nhấn `Tab` → Tab mode ON
3. **Result:** Tab + S combination mode

### Tắt một trong hai

- Nhấn `Tab` lại → Chỉ còn S (sequence numbers)
- Nhấn `S` lại → Chỉ còn Tab (all note names)

## 🔧 Implementation

### Rendering Logic (fretboard.go)

```go
if props.ShowAll {
    // Tab mode
    for all frets on fretboard {
        note := CalculateNote(...)
        
        // Check if note is in scale
        if seqItem, inScale := props.ScaleSequence[key]; inScale && props.ShowScaleShape {
            // S + Tab mode: Note name with finger background
            style = fingerBgStyles[seqItem.Finger]
            style = style.Foreground(BLACK)  // Chữ đen
            text = noteName
        } else {
            // Tab only: Note name with note color
            style = NoteColors[note]
            text = noteName
        }
    }
}
```

### Auto-Disable Rules

**Tab key:**
- Disables: Upcoming (U), Fingers (H)
- Keeps: Scale Shape (S) ✅

**S key:**
- Disables: Upcoming (U), Fingers (H)
- Keeps: Tab ✅

**H key:**
- Disables: Tab, S, Upcoming
- Standalone mode

## 📊 Mode Combinations

| Tab | S | H | U | Result |
|-----|---|---|---|--------|
| ON | OFF | OFF | OFF | All note names (colored text) |
| OFF | ON | OFF | OFF | Sequence numbers (colored bg) |
| ON | ON | OFF | OFF | **Note names with finger backgrounds** |
| OFF | OFF | ON | OFF | Finger numbers for lesson notes |
| OFF | OFF | OFF | ON | Upcoming preview only |
| OFF | ON | OFF | ON | Sequence + Upcoming ✅ |
| OFF | OFF | ON | ON | Fingers + Upcoming ✅ |
| ON | ON | - | - | **Tab+S combo** ✅ |
| ON | - | ON | - | ❌ Conflict (H disables Tab) |
| - | ON | ON | - | ❌ Conflict (H disables S) |

## 🎯 Use Cases

### 1. Learning Scale Patterns
```
Mode: Tab + S
Purpose: See where scale notes are on entire fretboard
         với finger colors để biết ngón nào
```

### 2. Understanding Note Relationships
```
Mode: Tab only
Purpose: See all notes, find intervals and patterns
```

### 3. Practicing Fingering
```
Mode: S only
Purpose: Focus on sequence and finger positions
```

### 4. Finger Guide
```
Mode: H only
Purpose: See finger numbers for lesson notes
```

## ✅ Benefits

**Tab + S Combination:**
- ✅ Xem được TẤT CẢ notes trên fretboard
- ✅ Scale notes nổi bật với màu nền ngón tay
- ✅ Chữ đen dễ đọc trên nền màu
- ✅ Dễ nhận biết notes nào trong scale, notes nào ngoài scale
- ✅ Học positions và relationships giữa các notes

**Example: A Minor Pentatonic**
```
Notes in scale: A, C, D, E, G
  - A at fret 5: Chữ "A" đen trên nền BLUE (ngón 1)
  - C at fret 8: Chữ "C" đen trên nền RED (ngón 4)
  - D at fret 5: Chữ "D" đen trên nền BLUE (ngón 1)
  - E at fret 7: Chữ "E" đen trên nền YELLOW (ngón 3)
  - G at fret 5: Chữ "G" đen trên nền BLUE (ngón 1)

Notes not in scale: F, F#, G#, A#, B...
  - F: Màu GREEN (note color)
  - F#: Màu TEAL (note color)
  - etc.
```

**Rất dễ nhìn và học patterns!** 🎸✨
