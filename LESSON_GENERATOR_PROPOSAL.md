# Lesson Generator - Đề xuất cải tiến

## 🔍 Vấn đề hiện tại

### 1. Logic Generator phức tạp
```go
// Có 2 paths khác nhau
- generateStepsWithPosition() // Dùng position data
- generateStepsLegacy()       // Scan toàn bộ fretboard
```

**Vấn đề:**
- Không rõ khi nào dùng path nào
- Legacy path không chuẩn (scan tất cả notes trong range)
- Không follow guitar scale patterns

### 2. Thứ tự notes không đúng
```go
// Hiện tại: Duyệt theo string sau đó theo fret
for s := 0; s < 6; s++ {
    for f := startFret; f <= endFret; f++ {
        // Add note
    }
}
```

**Kết quả:** Notes không theo pattern tự nhiên của guitar scale box

### 3. Finger pattern không chính xác
```go
finger := f - config.StartFret + 1  // Tính toán đơn giản
if finger > 4 {
    finger = 4
}
```

**Vấn đề:** Không match với finger pattern thực tế trong positions data

### 4. Beat assignment vô nghĩa
```go
Beat: (i % 4) + 1  // 1,2,3,4,1,2,3,4...
```

**Vấn đề:** Không có ý nghĩa âm nhạc, chỉ là modulo

## 💡 Đề xuất giải pháp

### Approach 1: Scale Box Pattern (Recommended)

#### Concept
- Mỗi position có một "box pattern" chuẩn
- Follow thứ tự từ string thấp → cao, note thấp → cao
- Sử dụng finger pattern có sẵn trong positions data

#### Ví dụ: A Minor Pentatonic Position 1
```
E |---5---8---| 
B |---5---8---| 
G |---5---7---| 
D |---5---7---| 
A |---5---7---| 
E |---5---8---|  ← Start here (root on 6th string)
   Fret 5-8
```

**Thứ tự chơi (ascending):**
1. String 6, Fret 5 (Root - A) - Finger 1
2. String 6, Fret 8 (C)        - Finger 4
3. String 5, Fret 5 (D)        - Finger 1
4. String 5, Fret 7 (E)        - Finger 3
5. String 4, Fret 5 (G)        - Finger 1
6. String 4, Fret 7 (A)        - Finger 3
... continue pattern

#### Implementation
```go
func GenerateSteps(config *GeneratorConfig) ([]Step, error) {
    // 1. Get position data
    position := theory.GetPosition(config.Scale, config.PosType, config.Position)
    
    // 2. Calculate actual frets based on root note
    rootFret := theory.FindRootFretOn6thString(root)
    startFret := rootFret + position.StartOffset
    
    // 3. Build markers following scale box pattern
    markers := buildScaleBoxPattern(position, startFret, root, config.Scale)
    
    // 4. Apply direction
    if config.Direction == "descending" {
        reverse(markers)
    }
    
    // 5. Create steps (1 note per step)
    steps := createStepsFromMarkers(markers)
    
    return steps
}

func buildScaleBoxPattern(pos Position, startFret int, root Note, scale string) []Marker {
    var markers []Marker
    
    // Traverse from low string to high string
    for stringIdx := 5; stringIdx >= 0; stringIdx-- {
        pattern := pos.NotePatterns[stringIdx]
        fingers := pos.FingerPattern[stringIdx]
        
        // Sort frets ascending
        frets := make([]int, len(pattern.RelativeFrets))
        for i, rel := range pattern.RelativeFrets {
            frets[i] = startFret + rel
        }
        
        // Add notes in order
        for i, fret := range frets {
            note := calculateNote(stringIdx, fret)
            
            // Only include if in scale
            if isNoteInScale(note, root, scale) {
                markers = append(markers, Marker{
                    StringIndex: stringIdx,
                    Fret:        fret,
                    Finger:      fingers[i],
                    Note:        note,
                })
            }
        }
    }
    
    return markers
}
```

### Approach 2: Musical Sequence Pattern

#### Concept
- Tạo sequence âm nhạc có ý nghĩa
- Patterns: Scale runs, arpeggios, licks

#### Ví dụ patterns:
```
1. Scale Run (ascending/descending)
2. Triplet pattern: 1-2-3, 2-3-4, 3-4-5...
3. Interval pattern: 1-3, 2-4, 3-5... (thirds)
4. Sequence: 1-2-3-2, 2-3-4-3, 3-4-5-4...
```

### Approach 3: Exercise Patterns

#### Concept
- Technical exercises cho guitar
- Finger dexterity, speed building

#### Ví dụ:
```
1. Chromatic: 1-2-3-4 on each string
2. String skipping: 6-4-5-3-4-2...
3. Hammer-on/Pull-off patterns
```

## �� Recommendation

**Implement Approach 1 first** vì:
1. ✅ Đơn giản, rõ ràng
2. ✅ Chuẩn theo guitar scale patterns
3. ✅ Sử dụng position data có sẵn
4. ✅ Dễ maintain và extend

**Sau đó có thể thêm:**
- Approach 2 cho musical exercises
- Approach 3 cho technical exercises

## 📝 Proposed Changes

### File: `internal/lesson/generator.go`

#### 1. Simplify entry point
```go
func GenerateSteps(config *GeneratorConfig) ([]Step, error) {
    root := parseNote(config.Root)
    
    // Always use position-based generation
    return generateScaleBoxSteps(config, root)
}
```

#### 2. Main generator function
```go
func generateScaleBoxSteps(config *GeneratorConfig, root Note) ([]Step, error) {
    // Get position
    posType := theory.PositionTypeCAGED
    if config.PosType == "3nps" {
        posType = theory.PositionType3NPS
    }
    
    position, exists := theory.GetPosition(config.Scale, posType, config.Position)
    if !exists {
        return nil, fmt.Errorf("position not found")
    }
    
    // Calculate start fret
    rootFret := theory.FindRootFretOn6thString(root)
    startFret := rootFret + position.StartOffset
    
    // Build pattern
    markers := buildScaleBoxPattern(position, startFret, root, config.Scale)
    
    // Apply direction
    if config.Direction == "descending" {
        reverseMarkers(markers)
    }
    
    // Create steps
    return createSteps(markers), nil
}
```

#### 3. Helper functions
```go
func buildScaleBoxPattern(pos Position, startFret int, root Note, scaleName string) []Marker
func reverseMarkers(markers []Marker)
func createSteps(markers []Marker) []Step
```

## 🎸 Example Output

### Before (rối):
```
String 6, Fret 5
String 6, Fret 6  ❌ Not in scale
String 6, Fret 7  ❌ Not in scale
String 6, Fret 8
String 5, Fret 5
...
```

### After (chuẩn):
```
String 6, Fret 5 (A - Root) Finger 1
String 6, Fret 8 (C)        Finger 4
String 5, Fret 5 (D)        Finger 1
String 5, Fret 7 (E)        Finger 3
String 4, Fret 5 (G)        Finger 1
String 4, Fret 7 (A)        Finger 3
...
```

## ❓ Questions for Review

1. **Pattern Order**: Low string → High string OK? (Hoặc ngược lại?)
2. **Beat Assignment**: Giữ simple (1 note per beat) hay tạo patterns phức tạp hơn?
3. **Root Note Highlighting**: Có cần đánh dấu root notes đặc biệt không?
4. **Skip Notes**: Có cho phép skip notes trong pattern không? (Ví dụ: chỉ play 3rds)

## 🚀 Next Steps

1. Review proposal này
2. Quyết định approach
3. Implement new generator
4. Test với các scales khác nhau
5. Update lessons.json nếu cần

Bạn muốn tôi implement approach nào?
