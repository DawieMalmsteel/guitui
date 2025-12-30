# Key Behavior - S và N

## Phím S - Scale Shape / Sequence Mode

### Khi metronome settings ĐÓNG:
- `s` hoặc `S` → Toggle **Scale Shape mode**
- Hiển thị số thứ tự (1, 2, 3...) của notes trong scale pattern
- Nếu bật Scale Shape → tự động TẮT Tab mode (Note Names)

### Khi metronome settings MỞ (nhấn M):
- `s` → Cycle sound types (Wood → Mechanical → Digital)
- `S` → Vẫn toggle Scale Shape mode

## Phím N - Switch Position

### Chức năng:
- Chuyển sang position tiếp theo trong cùng scale
- Ví dụ: Position 1 → Position 2 → Position 3 → ... → Position 1

### Khi nhấn N:
1. ✅ `currentPosition` tăng lên
2. ✅ Regenerate `steps` từ position mới
3. ✅ Reset `currentStep = 0`
4. ✅ **ScaleSequence tự động rebuild** từ steps mới

### Display modes (S, H, Tab, U) vẫn GIỮ NGUYÊN khi nhấn N

**Ví dụ:**
```
1. Nhấn S → Bật Scale Shape mode (hiện số thứ tự)
2. Nhấn N → Switch sang Position 2
3. Scale Shape mode VẪN BẬT → Hiện số thứ tự của Position 2
```

## Behavior hiện tại (ĐÚNG)

| Action | Result |
|--------|--------|
| Nhấn `S` | Toggle Scale Shape mode |
| Nhấn `N` | Switch position, **giữ nguyên** display modes |
| Nhấn `S` → `N` | Position mới với Scale Shape **VẪN BẬT** |
| Nhấn `N` → `S` | Bật Scale Shape cho position hiện tại |

## Có thể bạn muốn?

### Option 1: Auto-refresh Scale Shape khi switch position
Khi nhấn N, tự động bật lại Scale Shape mode?

### Option 2: Show pattern type info
Hiển thị pattern type (box/sequence/exercise) trong info bar?

### Option 3: Reset display modes khi switch position
Khi nhấn N, reset tất cả display modes về default?

Bạn muốn behavior nào?
