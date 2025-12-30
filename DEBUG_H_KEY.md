# Debug H Key Issue

## Expected Behavior
- Press H → `showFingers = true`
- Background: Show finger numbers (1,2,3,4) for all lesson notes
- Active note: Show finger number with bold/underline

## Current Code

### H Key Handler (model.go:164-165)
```go
case "h", "H":
    m.showFingers = !m.showFingers
```
✅ Simple toggle - should work

### Layer 0 Rendering (fretboard.go:128-145)
```go
if props.ShowFingers && !props.ShowScaleShape && !props.ShowAll {
    for key, seqItem := range props.ScaleSequence {
        fingerText := fmt.Sprintf(" %d ", seqItem.Finger)
        // ... render
    }
}
```
⚠️ **PROBLEM**: Condition `!props.ShowAll`

If Tab is ON (showAll=true), this won't render!

### Active Note Rendering (fretboard.go:189-202)
```go
} else if props.ShowFingers {
    displayText = fmt.Sprintf(" %d ", m.Finger)
    // ... render
}
```
✅ Should work for active note

## Hypothesis

**If user previously pressed Tab:**
- `showAll = true`
- Press H → `showFingers = true`
- Layer 0 condition fails: `showFingers && !showAll` = `true && false` = **FALSE**
- Background notes: Show Tab mode (note names) ❌
- Active note: Shows finger number ✅

**Result**: Active note shows finger, but background shows note names!

## Fix

Remove `!props.ShowAll` condition from Layer 0 Finger mode:

```go
// Before
if props.ShowFingers && !props.ShowScaleShape && !props.ShowAll {

// After  
if props.ShowFingers && !props.ShowScaleShape {
```

Or better: Make H key auto-disable Tab:

```go
case "h", "H":
    m.showFingers = !m.showFingers
    if m.showFingers {
        m.showAll = false  // ← Add this
    }
```
