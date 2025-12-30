# Guitar Tab Notation - Quick Reference

## 🎸 Essential Symbols

### Basic
| Symbol | Meaning | Example |
|--------|---------|---------|
| `5` | Fret number | Play 5th fret |
| `0` | Open string | No finger |
| `5f1` | Fret + finger | 5th fret, index finger |
| `x` | Muted | Dead note |
| `-` | No note | Spacing |

### Techniques
| Symbol | Name | Example | Description |
|--------|------|---------|-------------|
| `h` | Hammer-on | `5h7` | Pick 5, hammer to 7 |
| `p` | Pull-off | `7p5` | Pick 7, pull to 5 |
| `b` | Bend | `7b9` | Bend 7 to pitch of 9 |
| `r` | Release | `7b9r7` | Bend then release |
| `/` | Slide up | `5/7` | Slide from 5 to 7 |
| `\` | Slide down | `7\5` | Slide from 7 to 5 |
| `~` | Vibrato | `5~` | Shake the string |
| `l` | Trill | `5l7` | Rapid hammer/pull |
| `t` | Tap | `12t` | Right-hand tap |
| `*` | Pinch harmonic | `7*` | Artificial harmonic |
| `<>` | Natural harmonic | `<12>` | Touch lightly at 12 |
| `v` | Whammy down | `5v` | Dip whammy bar |
| `^` | Whammy up | `5^` | Lift whammy bar |
| `s` | Legato slide | `5s7` | Smooth slide |
| `PM` | Palm mute | Below tab | Muted picking |

### Picking
| Symbol | Meaning |
|--------|---------|
| `d` | Down stroke ∏ |
| `u` | Up stroke V |
| `>` | Accent (loud) |
| `.` | Staccato (short) |

## 📋 File Format

```
TITLE: Lesson Name
BPM: 120
KEY: A
CATEGORY: scale | exercise | song
DIFFICULTY: beginner | intermediate | advanced
TUNING: EADGBE

e|5(f1)|-----|7(f3)|-----|
B|-----|5(f1)|-----|7(f3)|
G|-----|-----|5(f2)|-----|
D|-----|-----|-----|5(f2)|
A|-----|-----|-----|-----|
E|-----|-----|-----|-----|

NOTES:
Description here
```

**CRITICAL RULES:**
- `|` separates beats - each cell is ONE beat
- **ALL strings MUST have the same number of `|` delimiters**
- Empty cells `|-----|` or `||` = rest beats
- Dashes `-` are visual only (use any amount)
- Notes in same vertical column play together (chords)

## 🎯 Common Patterns

### Hammer-on/Pull-off Legato
```
e|--5h7p5h7p5--|
```

### Triplet Run
```
e|--5-7-8--7-5-3--|
     3       3
```

### Bend and Vibrato
```
e|--7b9~--|
```

### Slide and Hammer
```
e|--5/7h8--|
```

### Pinch Harmonic Dive
```
e|--7*v---|
```

### Trill Pattern
```
e|--5l7ll5l7--|
```

## 🔤 String Names

```
e = 1st string (high E) - thinnest
B = 2nd string
G = 3rd string  
D = 4th string
A = 5th string
E = 6th string (low E) - thickest
```

## 👆 Finger Numbers

```
0 = Open string
1 = Index finger
2 = Middle finger
3 = Ring finger
4 = Pinky
```

## ⚡ Quick Examples

### Simple Scale (Sequential Notes)
```
e|5(f1)|-----|7(f3)|-----|8(f4)|
B|-----|5(f1)|-----|7(f3)|-----|
G|-----|-----|5(f2)|-----|-----|
```
Each beat plays one note. All strings have 5 beats.

### Chord Progression (Notes Together)
```
e|5(f1)|7(f3)|8(f4)|
B|5(f1)|7(f3)|8(f4)|
G|5(f2)|7(f4)|9(f4)|
```
Each beat plays multiple strings together (chord).

### Metal Riff
```
E|0(f0)|0(f0)|0(f0)|3(f3)|5(f1)|
A|-----|-----|-----|-----|-----|
   PM-----------------------
```

### Blues Lick
```
e|8(f4)|-----|5(f1)|-----|
B|-----|-----|-----|8(f4)|
```

### Sweep Arpeggio
```
e|--12----------------
B|-----10-------------
G|--------9-----------
D|----------11--------
   →→→→→
```

## 🎼 Common Combinations

| Notation | Name | Description |
|----------|------|-------------|
| `5h7p5` | Trill fragment | Quick hammer/pull |
| `7b9r7` | Bend/release | Full step up and down |
| `5/7~` | Slide vibrato | Slide then shake |
| `<12>*` | Harm + pinch | Combined harmonics |
| `5h7t12p7` | Tap lick | Hammer, tap, pull |
| `3PM` | Palm mute | Muted power chord |

## 📖 Legend Template

```
LEGEND:
h   = hammer-on          l   = trill
p   = pull-off           ~   = vibrato  
b   = bend               v   = whammy down
r   = release            ^   = whammy up
/   = slide up           *   = pinch harm
\   = slide down         <>  = natural harm
s   = legato slide       t   = tap
x   = muted              PM  = palm mute
d   = down stroke        u   = up stroke
```

---

**Tip:** Most online guitar tabs use these exact symbols. Copy/paste works!

**File Extension:** Save as `.tab` or `.txt`

**Version:** 1.0
