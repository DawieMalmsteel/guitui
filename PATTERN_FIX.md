# Guitar Pattern Fix - String Order

## 🐛 Vấn đề

Pattern không đúng theo cách chơi guitar tự nhiên.

**Output SAI (trước khi fix):**
```
1. String 1, Fret 5 → A (cao nhất)
2. String 1, Fret 8 → C
...
11. String 6, Fret 5 → A (thấp nhất)
12. String 6, Fret 8 → C
```

Chơi từ **cao xuống thấp** → Nghe không tự nhiên!

## ✅ Output ĐÚNG (sau khi fix)

```
1. String 6, Fret 5 → A (thấp nhất - 110Hz)
2. String 6, Fret 8 → C
3. String 5, Fret 5 → D
4. String 5, Fret 7 → E
...
11. String 1, Fret 5 → A (cao nhất - 440Hz)
12. String 1, Fret 8 → C
```

Chơi từ **thấp lên cao** → Đúng guitar pattern! ✅

## 🔍 Nguyên nhân

### StandardTuning Array Indexing

```go
// Index 0 là dây to nhất (E2), Index 5 là dây nhỏ nhất (E4)
var StandardTuning = []Note{E, A, D, G, B, E}
```

- `StandardTuning[0]` = E = **String 6** (thấp nhất)
- `StandardTuning[1]` = A = **String 5**
- `StandardTuning[2]` = D = **String 4**
- `StandardTuning[3]` = G = **String 3**
- `StandardTuning[4]` = B = **String 2**
- `StandardTuning[5]` = E = **String 1** (cao nhất)

### Code Trước (SAI)

```go
// Loop từ 5 → 0 = String 1 → 6 = Ngược!
for stringIdx := 5; stringIdx >= 0; stringIdx-- {
    pattern := pos.NotePatterns[stringIdx]
    ...
}
```

Loop từ index 5→0 = chơi từ String 1 cao xuống String 6 thấp = **SAI!**

### Code Sau (ĐÚNG)

```go
// Loop từ 0 → 5 = String 6 → 1 = Đúng!
for stringIdx := 0; stringIdx < 6; stringIdx++ {
    pattern := pos.NotePatterns[stringIdx]
    ...
}
```

Loop từ index 0→5 = chơi từ String 6 thấp lên String 1 cao = **ĐÚNG!**

## 🎸 Ví dụ: A Minor Pentatonic Position 1

### Sau khi fix:

```
Pitch   Note  String  Fret  Finger
-----   ----  ------  ----  ------
110Hz   A     6       5     1      ← Bắt đầu (thấp nhất)
130Hz   C     6       8     4
146Hz   D     5       5     1
164Hz   E     5       7     3
196Hz   G     4       5     1
220Hz   A     4       7     3
261Hz   C     3       5     1
293Hz   D     3       7     3
329Hz   E     2       5     1
392Hz   G     2       8     4
440Hz   A     1       5     1
523Hz   C     1       8     4      ← Kết thúc (cao nhất)
```

**Pattern tự nhiên:** Từ note thấp → cao, theo cách guitar player thực tế chơi scales!

## 📊 Impact

### Trước Fix
- ❌ Chơi ngược (cao → thấp)
- ❌ Không tự nhiên
- ❌ Khó học scale patterns
- ❌ Confusing cho người mới

### Sau Fix
- ✅ Chơi đúng (thấp → cao)
- ✅ Tự nhiên theo guitar
- ✅ Dễ học scale box patterns
- ✅ Match với guitar lessons thông thường

## 🎯 Áp dụng cho

Fix này áp dụng cho TẤT CẢ 3 pattern types:

1. **Box Pattern** ✅
2. **Sequence Pattern** ✅ (sử dụng box pattern làm base)
3. **Exercise Pattern** ✅ (chromatic, string skip, hammer-pull)

## 🔧 File Changed

- `internal/lesson/generator.go`
  - Function: `buildScaleBoxMarkers()`
  - Change: Loop direction `5→0` thành `0→5`

## ✨ Kết quả

Bây giờ tất cả patterns đều follow thứ tự guitar chuẩn:
- **Ascending**: String 6→1 (thấp→cao) ✅
- **Descending**: String 1→6 (cao→thấp) ✅

**Guitar patterns giờ 100% chuẩn!** 🎸
