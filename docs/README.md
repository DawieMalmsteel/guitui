# Guitar Tab Format Documentation

This directory contains documentation for the ASCII Guitar Tab format used in this application.

## 📚 Documentation Files

### [TAB_FORMAT.md](./TAB_FORMAT.md)
**Complete specification** of the guitar tab format including:
- File structure and metadata
- All technique notations (bends, slides, harmonics, etc.)
- Rhythm and timing notation
- Parser implementation details
- Writing guidelines and best practices

### [TAB_QUICK_REFERENCE.md](./TAB_QUICK_REFERENCE.md)
**Quick lookup guide** for common symbols:
- Essential symbols table
- Common patterns
- Quick examples
- Legend template

## 🎸 Example Lessons

The `lessons_tab/` directory contains example lesson files:

1. **01_a_minor_pentatonic_box1.tab** - Basic scale pattern
2. **02_chromatic_1234.tab** - Finger exercise
3. **03_blues_lick_bends.tab** - String bending techniques
4. **04_metal_palm_mute.tab** - Palm muting and power chords
5. **05_tapping_exercise.tab** - Two-hand tapping

## 🚀 Quick Start

### Creating a Lesson

1. Create a `.tab` file
2. Add metadata header:
   ```
   TITLE: Your Lesson Name
   BPM: 120
   KEY: A
   CATEGORY: scale | exercise | song | technique
   DIFFICULTY: beginner | intermediate | advanced
   TUNING: EADGBE
   ```

3. Add the tab:
   ```
   e|--------------------------|
   B|--------------------------|
   G|--------------------------|
   D|--------------------------|
   A|--------------------------|
   E|--------------------------|
   ```

4. Add notes section:
   ```
   NOTES:
   Description and tips here
   ```

### Basic Notation

- **Fret numbers**: `5`, `7`, `12`
- **Finger info**: `5f1` (fret 5, finger 1)
- **Hammer-on**: `5h7`
- **Pull-off**: `7p5`
- **Bend**: `7b9`
- **Slide**: `5/7` or `7\5`
- **Vibrato**: `5~`
- **Trill**: `5l7` (changed from 'tr' for clarity)
- **Tap**: `12t`
- **Harmonic**: `<12>` or `7*`
- **Muted**: `x`

## 📖 Learning Path

1. **Start here**: Read [TAB_QUICK_REFERENCE.md](./TAB_QUICK_REFERENCE.md)
2. **Deep dive**: Study [TAB_FORMAT.md](./TAB_FORMAT.md)
3. **Practice**: Look at examples in `lessons_tab/`
4. **Create**: Write your own lesson files

## 🔧 Technical Details

### File Extensions
- `.tab` (preferred)
- `.txt` (also supported)
- `.guitar` (custom extension)

### String Indexing

**In tab notation:**
```
e| = String 1 (high E) - thinnest
B| = String 2
G| = String 3
D| = String 4
A| = String 5
E| = String 6 (low E) - thickest
```

**In code (0-indexed):**
```
index 0 = String 6 (E)
index 1 = String 5 (A)
index 2 = String 4 (D)
index 3 = String 3 (G)
index 4 = String 2 (B)
index 5 = String 1 (e)
```

### Supported Tunings

Standard tuning (EADGBE) is default. Other tunings can be specified:
- Drop D: DADGBE
- Drop C: CGCFAD
- Half step down: Eb Ab Db Gb Bb Eb
- Custom: Any 6-letter combination

## 🎯 Why This Format?

### Advantages

✅ **Standard**: Used by millions of guitarists worldwide
✅ **Visual**: See the fretboard layout
✅ **Import**: Copy/paste from Ultimate Guitar, Songsterr, etc.
✅ **Edit**: Write in any text editor
✅ **Learn**: Intuitive for guitar players
✅ **Flexible**: Supports all techniques

### Comparison

| Format | Readable | Standard | Import | Techniques |
|--------|----------|----------|--------|------------|
| JSON | ❌ | ✅ | ❌ | Limited |
| YAML | ✅ | ✅ | ❌ | Limited |
| ASCII Tab | ✅✅ | ✅✅ | ✅✅ | ✅✅ |

## 📋 Notation Changes

### From Standard

**Changed notation for clarity:**
- `tr` → `l` (trill) - easier to distinguish from other symbols

**All other notations match standard guitar tab format.**

## 🌐 Resources

### Tab Websites
- [Ultimate Guitar](https://www.ultimate-guitar.com) - Largest tab database
- [Songsterr](https://www.songsterr.com) - Interactive tabs
- [Guitar Pro](https://www.guitar-pro.com) - Professional software

### Learning
- Search "how to read guitar tabs" for tutorials
- Most YouTube guitar lessons use this format
- Free tabs available for thousands of songs

## 🔄 Migration from JSON

If you have existing JSON lessons:
1. Extract fret/finger information
2. Arrange in tab grid format
3. Add metadata header
4. Add technique notations as needed

The parser will handle conversion to internal format.

## 📝 Contributing

When adding new lessons:
- Follow the format specification
- Include finger numbers for educational value
- Add detailed NOTES section
- Use consistent spacing and alignment
- Test tab visually before committing

## 🆘 Support

Questions about notation? Check:
1. [TAB_QUICK_REFERENCE.md](./TAB_QUICK_REFERENCE.md) - Common symbols
2. [TAB_FORMAT.md](./TAB_FORMAT.md) - Complete specification
3. Example files in `lessons_tab/`

---

**Version:** 1.0  
**Last Updated:** 2025-12-30  
**Format:** ASCII Guitar Tab Standard
