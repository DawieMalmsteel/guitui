# Guitar Engine - Keyboard Controls

## TUI App - Hoàn toàn điều khiển bằng phím tắt

Đây là Terminal UI (TUI) application, không sử dụng chuột. Tất cả chức năng được điều khiển bằng bàn phím.

## 🎵 Main Controls

| Phím | Chức năng |
|------|-----------|
| `Space` | Play / Pause metronome và progression |
| `Enter` | Chọn bài lesson trong danh sách |
| `↑` `↓` | Di chuyển trong danh sách lessons |
| `Ctrl+C` / `Q` | Thoát app |

## 🎛️ Metronome Settings

| Phím | Chức năng |
|------|-----------|
| `M` | Mở/đóng **Metronome Settings Panel** |

### Khi Settings Panel mở (nhấn M):

| Phím | Chức năng |
|------|-----------|
| `+` hoặc `=` | Tăng BPM (+5) |
| `-` hoặc `_` | Giảm BPM (-5) |
| `1` | Set time signature 4/4 (Common) |
| `2` | Set time signature 3/4 (Waltz) |
| `3` | Set time signature 6/8 (Compound) |
| `4` | Set time signature 2/4 (March) |
| `S` | Cycle sound types: Wood → Mechanical → Digital |
| `M` | Đóng settings panel |

**Range BPM**: 40 - 240

## 🎸 Display Modes

| Phím | Chức năng |
|------|-----------|
| `H` | Toggle **Finger Helper** - Hiển thị số ngón tay |
| `S` | Toggle **Scale Shape/Sequence** - Hiển thị thứ tự note |
| `Tab` | Toggle **Note Names** - Hiển thị tên note |
| `U` | Toggle **Upcoming markers** - Hiển thị note sắp tới |

> **Note**: Phím `S` có 2 chức năng:
> - Khi **Settings Panel đóng**: Toggle Scale Shape
> - Khi **Settings Panel mở**: Cycle sound types

## 🎼 Position & Layout

| Phím | Chức năng |
|------|-----------|
| `N` | Switch Position (CAGED/3NPS positions) |
| `P` | Toggle Position Type (CAGED ↔ 3NPS) |
| `F` | Toggle Fret count (12 ↔ 24 frets) |

## 🎨 Metronome Settings Panel

Khi nhấn phím `M`, panel settings sẽ xuất hiện ở giữa màn hình với layout:

```
╔═══════════════════════════════════════════╗
║      ♪ METRONOME SETTINGS ♪               ║
║            ▶ PLAYING                      ║
║                                           ║
║ ──────────────────────────────────────    ║
║                                           ║
║ TEMPO (BPM)  ▸ 120 ◂                     ║
║   Press [+] to increase, [-] to decrease  ║
║                                           ║
║ TIME SIGNATURE  4/4 (Common Time)         ║
║   [1] 4/4   [2] 3/4   [3] 6/8   [4] 2/4  ║
║                                           ║
║ SOUND TYPE  🪵 Wood Block                ║
║   Natural wood percussion                 ║
║   Press [S] to cycle sounds               ║
║                                           ║
║ ──────────────────────────────────────    ║
║                                           ║
║ KEYBOARD CONTROLS                         ║
║                                           ║
║   [Space]  Play / Pause metronome         ║
║   [M]      Close this menu                ║
║   [+/-]    Adjust tempo                   ║
║   [1-4]    Change time signature          ║
║   [S]      Cycle sound types              ║
╚═══════════════════════════════════════════╝
```

## 🔊 Sound Types

| Type | Mô tả | Đặc điểm |
|------|-------|----------|
| 🪵 **Wood Block** | Natural wood percussion | Âm thanh tự nhiên, harmonics phong phú |
| ⚙️ **Mechanical** | Sharp mechanical click | Click sắc, ngắn, crisper |
| 🔔 **Digital Beep** | Clean electronic tone | Beep điện tử, nhẹ nhàng |

## 💡 Tips

1. **Metronome first**: Nhấn `M` để mở settings, adjust BPM và sound type trước khi practice
2. **Practice flow**: 
   - Chọn lesson (`Enter`)
   - Mở metronome settings (`M`)
   - Điều chỉnh BPM phù hợp
   - Đóng settings (`M`)
   - Play (`Space`)
3. **Display modes**: Dùng `H`, `S`, `Tab`, `U` để toggle các mode hiển thị theo nhu cầu
4. **Position switching**: Dùng `N` để practice các position khác nhau của cùng 1 scale

## 🎹 Quick Start Example

```
1. Run app: ./guitui
2. Select lesson: ↑↓ + Enter
3. Open metronome: M
4. Set BPM 80: Press - nhiều lần
5. Choose Wood sound: S (cycle đến Wood)
6. Close settings: M
7. Start practice: Space
8. Toggle helpers: H (fingers), S (sequence)
```

## ⚠️ Important Notes

- **NO MOUSE**: Đây là TUI app, không sử dụng chuột
- **Settings Panel**: Phải mở panel (`M`) mới điều chỉnh được metronome
- **Sound changes**: Chỉ có hiệu lực khi panel settings đang mở
- **BPM limits**: Minimum 40, Maximum 240
