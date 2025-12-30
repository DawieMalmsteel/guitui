# Guitar Tab Format Specification

## 📋 Overview

This application uses **ASCII Guitar Tab** format for lesson files - the same format used by guitarists worldwide. This makes it easy to:
- ✅ Write lessons by hand
- ✅ Import tabs from websites (Ultimate Guitar, Songsterr, etc.)
- ✅ Copy/paste existing tabs
- ✅ Visual and intuitive

---

## 📄 File Structure

### Basic Template

```
TITLE: Song/Exercise Name
BPM: 120
KEY: A
CATEGORY: scale | exercise | song | technique
DIFFICULTY: beginner | intermediate | advanced
TUNING: EADGBE

e|-----|-----|
B|-----|-----|
G|-----|-----|
D|-----|-----|
A|-----|-----|
E|-----|-----|

NOTES:
Optional notes about the lesson
```

### Example Lesson File

```
TITLE: A Minor Pentatonic - Box 1
BPM: 80
KEY: A
CATEGORY: scale
DIFFICULTY: beginner
TUNING: EADGBE

e|5(f1)|8(f4)|
B|5(f1)|7(f3)|
G|5(f1)|7(f3)|
D|5(f1)|7(f3)|
A|5(f1)|8(f4)|
E|5(f1)|8(f4)|

NOTES:
- Start with index finger at 5th fret
- Box pattern - all fingers stay in position
- Practice ascending and descending
```

---

## 🎸 String Notation

### String Names (Standard Tuning)

```
e = 1st string (high E) - thinnest
B = 2nd string (B)
G = 3rd string (G)
D = 4th string (D)
A = 5th string (A)
E = 6th string (low E) - thickest
```

### String Order in Tab

```
e|-----|  (String 1 - highest pitch, thinnest)
B|-----|  (String 2)
G|-----|  (String 3)
D|-----|  (String 4)
A|-----|  (String 5)
E|-----|  (String 6 - lowest pitch, thickest)
```

**Note:** In code, strings are indexed 0-5:
- String 6 (E) = index 0
- String 1 (e) = index 5

---

## 🎯 Basic Notation

### Fret Numbers

```
e|--0--1--2--3--5--7--8--12--|
   (play these frets on string)

0 = open string (no finger)
1 = 1st fret
5 = 5th fret
etc.
```

### Beat Timing with Pipe Delimiters

**IMPORTANT:** Beats are separated by pipe `|` delimiters. Each cell between pipes is one beat.

```
e|5(f1)|7(f3)|-----|8(f4)|
B|-----|5(f1)|7(f3)|-----|
G|-----|-----|-----|-----|

Beat 1: Note on e string (fret 5, finger 1)
Beat 2: Notes on e and B strings together (chord)
Beat 3: Rest (empty or dashes only)
Beat 4: Note on e string (fret 8, finger 4)
```

**Rules:**
- `|` separates beats
- Notes in the same cell play together (chords)
- Empty cells or cells with only `---` are rests
- Dashes `-` are visual only (ignored by parser)
- Use dashes to align tab for readability

### Finger Numbers (Optional)

Add finger information using parentheses notation:

```
e|5(f1)|7(f3)|8(f4)|

Format: {fret}(f{finger})
  5(f1) = fret 5, finger 1 (index)
  7(f3) = fret 7, finger 3 (ring)
  8(f4) = fret 8, finger 4 (pinky)
```

**Finger Numbers:**
- `0` = open string
- `1` = index finger
- `2` = middle finger
- `3` = ring finger
- `4` = pinky

### Alternative Notations

All these formats are supported for finger specification:

```
5(f1)    = fret 5, finger 1 (recommended)
5        = fret 5, no finger info
```

---

## 🎵 Technique Notation

### 1. Muted/Dead Notes

```
MUTED STRING:
e|--x--x--x--|  (muted, no pitch)

PALM MUTE:
e|--5--5--5--|
   PM-------

GHOST NOTE:
e|--(5)--|  (barely audible)
```

### 2. Bends

```
FULL BEND (whole step):
e|--7b9--|
   (bend from fret 7 to pitch of fret 9)

HALF BEND (1/2 step):
e|--7b8--|

QUARTER BEND:
e|--7b1/4--|

BEND AND RELEASE:
e|--7b9r7--|
   (bend up, then release down)

PRE-BEND:
e|--b7--|
   (bend string before picking)

PRE-BEND AND RELEASE:
e|--b7r5--|
```

### 3. Vibrato

```
VIBRATO:
e|--5~--|  or  e|--5~~~--|
   (shake the string)

WIDE VIBRATO:
e|--5~~--|  (wider shakes)

WHAMMY BAR VIBRATO:
e|--5~v~--|
```

### 4. Slides

```
SLIDE UP:
e|--5/7--|  (slide from 5 to 7)

SLIDE DOWN:
e|--7\5--|  (slide from 7 to 5)

SLIDE IN FROM BELOW:
e|--/5--|  (slide into 5 from below)

SLIDE IN FROM ABOVE:
e|--\5--|  (slide into 5 from above)

SLIDE OUT:
e|--5/--|  or  e|--5\--|

LEGATO SLIDE:
e|--5s7--|  (smooth, don't pick 2nd note)
```

### 5. Hammer-ons & Pull-offs

```
HAMMER-ON:
e|--5h7--|
   (pick 5, hammer finger onto 7)

PULL-OFF:
e|--7p5--|
   (pick 7, pull finger off to 5)

TRILL (changed from 'tr' to 'l'):
e|--5l7--|  or  e|--5l7l5l7--|
   (rapid hammer-on/pull-off)

ASCENDING:
e|--5h7h8h10--|

DESCENDING:
e|--10p8p7p5--|

COMBO:
e|--5h7p5h7p5--|
```

### 6. Tapping

```
TWO-HAND TAP:
e|--5h12t5h12p5--|
   (t = right-hand tap)

TAP NOTATION:
e|--12t--|  or  e|--T12--|

TAP HARMONIC:
e|--<12>t--|
```

### 7. Harmonics

```
NATURAL HARMONIC:
e|--<12>--|  or  e|--NH12--|
   (lightly touch at fret 12)

PINCH HARMONIC:
e|--7*--|  or  e|--PH7--|
   (artificial harmonic with pick)

TAPPED HARMONIC:
e|--TH12--|

Common positions: <5>, <7>, <12>, <19>, <24>
```

### 8. Tremolo Bar (Whammy)

```
DIP:
e|--5v--|  (dip down)

LIFT:
e|--5^--|  (lift up)

DIVE BOMB:
e|--12\---|  (dive to slack)

SCOOP:
e|--/5--|  (scoop up into note)
```

### 9. Picking

```
DOWN STROKE:
e|--5--5--5--|
   d  d  d

UP STROKE:
e|--5--5--5--|
   u  u  u

ALTERNATE PICKING:
e|--5--5--5--5--|
   d  u  d  u

TREMOLO PICKING:
e|--5≈--|  (very fast repeated picking)

SWEEP PICKING:
e|--12----------
B|-----10-------
G|--------9-----
   →→→ (arrow indicates sweep direction)

RAKE:
e|--x--x--5--|
   (rake across muted strings)
```

### 10. Special Techniques

```
PICK SCRAPE:
e|--PS----|  or  e|--///--|

LET RING:
e|--5-------|
   let ring--

STACCATO:
e|--5.--5.--5.--|  (short, choppy)

ACCENT:
e|-->5--|  or  e|--5>--|  (hit harder)

FEEDBACK:
e|--5---|
   fdbk
```

---

## 🎼 Rhythm & Timing

### Measures

```
MEASURE BARS:
e|--5--7--8--|--5--7--8--|--3--5--7--|
   measure 1   measure 2   measure 3
```

### Repeats

```
SIMPLE REPEAT:
e||:--5--7--|--8--10--:||
   (play this section twice)

1ST/2ND ENDING:
e|--5--7--|--8--10--|
  1.            2.
  (play 1st, repeat, then play 2nd)
```

### Note Durations

```
WHOLE NOTE:
e|--5-------|  (hold for 4 beats)

HALF NOTE:
e|--5----|  (hold for 2 beats)

QUARTER NOTE:
e|--5--|  (1 beat)

EIGHTH NOTE:
e|-5-5-5-5-|  (0.5 beat each)

TRIPLET:
e|--5-5-5--|
     3      (3 notes in time of 2)
```

---

## 📊 Complete Example

### Example 1: Advanced Techniques

```
TITLE: Metal Lick - Harmonic Minor
BPM: 140
KEY: Am
CATEGORY: technique
DIFFICULTY: advanced
TUNING: EADGBE

e|--------------------12h15p12----<12>*---------------------|
B|-----------12h15p13---------15b17r15~-------5h7p5---------|
G|--12h14p12----------------------------7b9r7-------7-------|
D|--------------------------------------------------9b11r9--|
A|----------------------------------------------------------|
E|----------------------------------------------------------|
   PM--------|                  v~~~~~    let ring---------|
   d  u  d  u                                 d   u   d

NOTES:
- Start with palm muting on the low strings
- Pinch harmonic at 12th fret on high e
- Use vibrato on bends for expression
- Let the final notes ring out

LEGEND:
h   = hammer-on          l   = trill (rapid h/p)
p   = pull-off           ~   = vibrato
b   = bend               v   = whammy bar
r   = release            <5> = natural harmonic
/   = slide up           *   = pinch harmonic
\   = slide down         PM  = palm mute
s   = legato slide       t   = tap
x   = muted              d   = down stroke
^   = lift/bend          u   = up stroke
```

### Example 2: Manual Beat Control

```
TITLE: Manual Beat Demo - Same String Multiple Beats
BPM: 100
KEY: C
CATEGORY: exercise
DIFFICULTY: beginner
TUNING: EADGBE

e|--5(b1f1)--7(b2f3)--8(b3f4)--10(b4f4)--|
B|----------------------------------------|
G|----------------------------------------|
D|----------------------------------------|
A|----------------------------------------|
E|----------------------------------------|

NOTES:
This demonstrates manual beat control.
Same string (high e) playing 4 different notes in sequence.

Without manual beats, this would be impossible to specify 
because the parser would see them all on the same string.

Beat notation:
- 5(b1f1) = fret 5, beat 1, finger 1
- 7(b2f3) = fret 7, beat 2, finger 3
- 8(b3f4) = fret 8, beat 3, finger 4
- 10(b4f4) = fret 10, beat 4, finger 4

LEGEND:
(b1f1) = beat 1, finger 1
(b2f3) = beat 2, finger 3
```

---

## 🔧 Parser Implementation

### Supported Patterns

The parser recognizes these patterns:

```regex
\d+           = fret number (0-24)
\d+f\d+       = fret + finger (5f1, 12f3)
\d+\(\d+\)    = fret + finger alt (5(1), 12(3))
x|X           = muted string
\d+b\d+       = bend (7b9, 5b6)
\d+b\d+r\d+   = bend and release (7b9r7)
\d+h\d+       = hammer-on (5h7)
\d+p\d+       = pull-off (7p5)
\d+l\d+      = trill (5l7)
\d+/\d+       = slide up (5/7)
\d+\\\d+      = slide down (7\5)
/\d+          = slide in from below
\\\d+         = slide in from above
\d+/          = slide out up
\d+\\         = slide out down
\d+s\d+       = legato slide
\d+~+         = vibrato (5~, 7~~~)
<\d+>         = natural harmonic (<12>)
\d+\*         = pinch harmonic (7*)
\d+t          = tap (12t)
\d+v          = whammy dip
\d+\^         = whammy lift
```

### String Line Detection

```
Line starts with:
  e| = string 1 (index 5 in code)
  B| = string 2 (index 4)
  G| = string 3 (index 3)
  D| = string 4 (index 2)
  A| = string 5 (index 1)
  E| = string 6 (index 0)
```

### Metadata Parsing

```
TITLE: {text}
BPM: {number}
KEY: {note}
CATEGORY: {text}
DIFFICULTY: {text}
TUNING: {6 letters}
NOTES: {multiline text}
```

---

## 📝 Writing Guidelines

### Best Practices

1. **Spacing**: Use dashes (`-`) for consistent spacing
2. **Alignment**: Keep notes vertically aligned across strings
3. **Clarity**: Add spaces between distinct phrases
4. **Comments**: Use NOTES section for instructions
5. **Finger hints**: Include finger numbers for learners

### Good Example

```
e|-----5f1-----7f3-----8f4-----|
B|-----5f1-----6f2-----8f4-----|
G|-----5f1-----7f3-------------|
   let ring-----------
```

### Poor Example (avoid)

```
e|--5--7--8|  (inconsistent spacing)
B|--5-6-8--|  (no finger info)
G|5-7------|  (misaligned)
```

---

## 🎯 File Extensions

Supported file extensions:
- `.tab` - Primary format
- `.txt` - Plain text tabs
- `.guitar` - Custom extension

---

## 📚 Resources

### Learning Resources
- [Ultimate Guitar](https://www.ultimate-guitar.com) - Huge tab database
- [Songsterr](https://www.songsterr.com) - Interactive tabs
- [Guitar Pro](https://www.guitar-pro.com) - Professional tab software

### ASCII Tab Tutorials
- Search for "how to read guitar tabs"
- Most online tabs use this exact format

---

## ⚙️ Future Enhancements

Planned additions:
- Multi-voice tabs (bass + lead)
- Drum notation
- Chord diagrams
- Time signature changes mid-song
- Tempo changes
- Effects pedal notations

---

**Last Updated:** 2025-12-30  
**Version:** 1.0  
**Format Compatibility:** ASCII Guitar Tab Standard
