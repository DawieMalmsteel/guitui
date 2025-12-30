# Test Guitar Pattern - A Minor Pentatonic Position 1

## Standard Guitar A Minor Pentatonic Position 1 (Box 1)

Theo chuẩn guitar, Position 1 của A Minor Pentatonic bắt đầu từ fret 5:

```
e |---5---8---|  (E string - cao nhất)
B |---5---8---|  
G |---5---7---|  
D |---5---7---|  
A |---5---7---|  
E |---5---8---|  (E string - thấp nhất, root note)
   Fret 5-8
```

## Cách chơi đúng (Ascending - từ thấp lên cao)

**Thứ tự chơi từ note thấp nhất → cao nhất:**

1. String 6 (E thấp), Fret 5 → **A** (Root) - Ngón 1
2. String 6 (E thấp), Fret 8 → **C** - Ngón 4
3. String 5 (A), Fret 5 → **D** - Ngón 1  
4. String 5 (A), Fret 7 → **E** - Ngón 3
5. String 4 (D), Fret 5 → **G** - Ngón 1
6. String 4 (D), Fret 7 → **A** - Ngón 3
7. String 3 (G), Fret 5 → **C** - Ngón 1
8. String 3 (G), Fret 7 → **D** - Ngón 3
9. String 2 (B), Fret 5 → **E** - Ngón 1
10. String 2 (B), Fret 8 → **G** - Ngón 4
11. String 1 (E cao), Fret 5 → **A** - Ngón 1
12. String 1 (E cao), Fret 8 → **C** - Ngón 4

## Scale degrees:
A Minor Pentatonic = A C D E G
- A = Root (1)
- C = ♭3
- D = 4
- E = 5
- G = ♭7

## Notes theo pitch (từ thấp → cao):

| Note | Pitch | String | Fret | Finger |
|------|-------|--------|------|--------|
| A    | 110Hz | 6      | 5    | 1      |
| C    | 130Hz | 6      | 8    | 4      |
| D    | 146Hz | 5      | 5    | 1      |
| E    | 164Hz | 5      | 7    | 3      |
| G    | 196Hz | 4      | 5    | 1      |
| A    | 220Hz | 4      | 7    | 3      |
| C    | 261Hz | 3      | 5    | 1      |
| D    | 293Hz | 3      | 7    | 3      |
| E    | 329Hz | 2      | 5    | 1      |
| G    | 392Hz | 2      | 8    | 4      |
| A    | 440Hz | 1      | 5    | 1      |
| C    | 523Hz | 1      | 8    | 4      |

## Kiểm tra code hiện tại

Code đang làm gì:
```go
// Traverse from low string (6th) to high string (1st)
for stringIdx := 5; stringIdx >= 0; stringIdx-- {
    // Get all frets for this string in ascending order
    for i, relFret := range pattern.RelativeFrets {
        actualFret := startFret + relFret
        // Add note
    }
}
```

Có vẻ đúng! Nhưng cần check:
1. **RelativeFrets có đúng không?**
2. **Finger pattern có đúng không?**
3. **Notes có filter đúng scale không?**

