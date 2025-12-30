# Changelog

## [2025-12-30] - Beat Parsing & Color Fixes

### Fixed
- **Beat parsing logic**: Fixed parser to correctly group all notes within the same beat into a single Step
  - Previously: Each note was creating a separate Step, causing notes to play sequentially even when in same beat
  - Now: All notes between `|` delimiters are grouped into one beat/step and play together
  
- **Beat alignment**: Fixed cell filtering logic that was skipping leading empty cells
  - Previously: Empty cells at start were removed, causing strings to become misaligned
  - Now: All cells between `|` are preserved, including empty ones (rest beats)
  
- **Finger colors**: Fixed active note color display
  - Previously: Active notes showed fixed peach color regardless of finger
  - Now: Active notes display with background color matching the finger number:
    - Finger 0 (open): Red
    - Finger 1 (index): Teal
    - Finger 2 (middle): Yellow
    - Finger 3 (ring): Peach
    - Finger 4 (pinky): Mauve (Purple) - was incorrectly Red before

### Changed
- **Tab file format requirements**: All strings must now have equal number of `|` delimiters
  - Parser now properly synchronizes beats across all strings
  - Empty cells `|-----|` or `||` represent rest beats
  - Example tab file updated: `01_a_minor_pentatonic_box1.tab`

### Documentation
- Updated `TAB_FORMAT.md` with beat alignment rules and examples
- Updated `TAB_QUICK_REFERENCE.md` with correct format examples
- Updated `README.md` with finger color information and format rules
- Added visual examples showing correct vs incorrect tab formatting

### Technical Details
- Modified `internal/lesson/tabparser.go`:
  - `parseSteps()`: Fixed to group markers by beat instead of creating separate steps
  - Cell filtering logic: Keeps all cells between delimiters for proper alignment
- Modified `internal/ui/components/fretboard.go`:
  - `buildActiveLayer()`: Changed to use `getFingerStyle()` for active notes
  - Finger color mapping: Updated finger 4 from Red to Mauve

## Previous Versions

Previous changes were not documented. This is the first changelog entry.
