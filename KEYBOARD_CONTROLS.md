# Keyboard Controls - Final Implementation

## Display Modes (Background Layer)

### 🔢 S - Scale Sequence Mode
**Function**: Hiển thị số thứ tự (1, 2, 3...) của notes trong lesson

**Behavior:**
- Shows sequence numbers for all notes in lesson
- Each number colored by finger (finger 1 = blue, 2 = green, etc.)
- Active note: underlined + bold
- Auto disables: ShowAll (Tab), Upcoming

**Example:**
```
E |---1---4---|  (Số 1, 4 với màu theo ngón)
B |---2---5---|
G |---3---6---|
```

### 🎵 Tab - Note Names Mode
**Function**: Hiển thị tên nốt nhạc trên **TOÀN BỘ** fretboard

**Behavior:**
- Shows ALL note names on entire fretboard (not just lesson notes)
- Each note colored by its pitch (A = purple, C = red, E = yellow...)
- Displays notes from fret 0 to fret 12 (or FretCount)
- Auto disables: ShowScaleShape (S), ShowUpcoming (U)

**Example:**
```
E |--E--F--F#-G--G#-A--A#-B--C--C#-D--D#-E--|
B |--B--C--C#-D--D#-E--F--F#-G--G#-A--A#-B--|
G |--G--G#-A--A#-B--C--C#-D--D#-E--F--F#-G--|
D |--D--D#-E--F--F#-G--G#-A--A#-B--C--C#-D--|
A |--A--A#-B--C--C#-D--D#-E--F--F#-G--G#-A--|
E |--E--F--F#-G--G#-A--A#-B--C--C#-D--D#-E--|
```

### 👆 H - Finger Helper Mode
**Function**: Hiển thị số ngón tay (1, 2, 3, 4) cho notes trong lesson

**Behavior:**
- Shows finger numbers for lesson notes only
- Colored background by finger
  - Finger 1 (index) = Blue background
  - Finger 2 (middle) = Green background
  - Finger 3 (ring) = Yellow background
  - Finger 4 (pinky) = Red background
  - Finger 0 (open) = Gray background
- Active note: bold + underline
- Works independently from S and Tab

**Example:**
```
E |---1---4---|  (Ngón 1, 4 với background màu)
B |---1---3---|
G |---1---3---|
```

### 👁️ U - Upcoming Markers Mode
**Function**: Hiển thị preview của 3 notes sắp tới

**Behavior:**
- Shows next 3 upcoming notes with distance indicator
- Distance 1: Bold arrow (▶)
- Distance 2: Regular arrow (→)
- Distance 3: Faint arrow (⇒)
- Disabled when: ShowScaleShape (S) or ShowAll (Tab) is ON
- Can combine with: ShowFingers (H)

## Metronome Controls

### ⏯️ Space - Play/Pause
Toggle metronome play/pause

### 🎛️ M - Metronome Settings
Open/close metronome settings panel

**In settings panel:**
- `+/-` - Adjust BPM (40-240, step 5)
- `1-4` - Change time signature (4/4, 3/4, 6/8, 2/4)
- `S` - Cycle sound types (Wood → Mechanical → Digital)
- `M` - Close settings

## Other Controls

### 🎸 F - Toggle Fret Count
Cycle through fret counts: 12 → 15 → 22 → 12

### 📋 Enter - Select Lesson
Select highlighted lesson from list and start playing

### ❌ Q / Ctrl+C - Quit
Exit application

## Key Combinations & Priorities

### Display Mode Exclusivity

**Mutually Exclusive:**
- `S` (Scale Shape) ⟷ `Tab` (Note Names)
- When one is ON, the other is automatically OFF

**Can Combine:**
- `H` (Fingers) + `U` (Upcoming) ✅
- `H` (Fingers) + Default mode ✅

**Auto-Disable Rules:**
- Press `S` → Disables `Tab`, `U`
- Press `Tab` → Disables `S`, `U`
- Press `Tab` again → Re-enables `U`

### Rendering Priority (Low to High)

1. **Background Layer** (Priority 1)
   - Scale Sequence (S)
   - Note Names (Tab)
   - Finger Numbers (H)

2. **Upcoming Layer** (Priority 2)
   - Upcoming markers (U)

3. **Active Layer** (Priority 3)
   - Currently playing note
   - Always visible, highest priority

## Summary Table

| Key | Mode | Shows | Scope | Combines With |
|-----|------|-------|-------|---------------|
| `S` | Scale Sequence | 1,2,3... | Lesson notes | H |
| `Tab` | Note Names | A,C,D... | Entire fretboard | - |
| `H` | Fingers | 1,2,3,4 | Lesson notes | S, U |
| `U` | Upcoming | Arrows | Next 3 notes | H |
| `Space` | Play/Pause | - | Metronome | All |
| `M` | Metro Settings | Panel | Metronome | All |
| `F` | Fret Count | 12/15/22 | Fretboard | All |

## Usage Examples

### Learning a Scale
1. Press `S` → See sequence numbers
2. Press `H` → See which fingers to use
3. Press `Space` → Start playing

### Understanding Note Positions
1. Press `Tab` → See all notes on fretboard
2. Find patterns and relationships
3. Press `Tab` again to turn off

### Practicing with Finger Guide
1. Press `H` → See finger numbers
2. Press `U` → See upcoming notes
3. Press `Space` → Practice

### Adjusting Metronome
1. Press `M` → Open settings
2. Press `+/-` to adjust BPM
3. Press `1-4` to change time signature
4. Press `S` to change sound
5. Press `M` to close

**All controls work perfectly!** 🎸✨
