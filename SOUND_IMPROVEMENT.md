# Metronome Sound Improvement

## Vấn đề cũ
- Âm thanh sine wave đơn giản, khó chịu khi nghe lâu
- Không có envelope (attack/decay) tạo âm thanh "flat"
- Chỉ có 1 loại âm thanh

## Cải tiến mới

### 1. **Wood Block Sound** (Mặc định - Tốt nhất)
- Âm thanh giống gõ khối gỗ của metronome thật
- Sử dụng 2 tần số harmonics (base + harmonic)
- Attack cực nhanh (1-2ms) + Decay tự nhiên (30-40ms)
- Accent beat: 1800Hz + 3600Hz
- Normal beat: 1400Hz + 2800Hz

### 2. **Mechanical Click**
- Âm thanh click sắc, ngắn gọn
- Attack 1ms + Decay 12-15ms
- Tần số cao hơn: 2000-2400Hz
- Âm thanh giống metronome cơ học

### 3. **Digital Beep**
- Âm thanh beep điện tử sạch
- Attack chậm hơn (5ms) + Decay (20ms)
- Tần số âm nhạc: A4 (440Hz) và A5 (880Hz)
- Nhẹ nhàng hơn, phù hợp cho practice nhẹ

## Kỹ thuật Audio

### Envelope Streamer
```go
type envelopeStreamer struct {
    Streamer       beep.Streamer
    attackSamples  int  // Attack phase duration
    decaySamples   int  // Decay phase duration
    currentSample  int
}
```

- **Attack**: Volume tăng từ 0 → 1 (tạo "punch")
- **Decay**: Volume giảm từ 1 → 0 (tự nhiên)

### Harmonic Mixing
```go
// Wood block uses two frequencies
baseTone := generators.SineTone(sampleRate, 1400.0)
harmonic := generators.SineTone(sampleRate, 2800.0)
mixed := beep.Mix(baseTone, volumeStreamer{harmonic, 0.3})
```

## Sử dụng

1. **Chọn sound type**: Press `[S]` khi metronome settings mở (phím `[M]`)
2. **Cycle qua các sounds**: Wood Block → Mechanical → Digital → Wood Block...
3. **Sound mặc định**: Wood Block (nghe tự nhiên nhất)

## Tham số Chi tiết

### Wood Block
| Beat   | Base Freq | Harmonic | Attack | Decay |
|--------|-----------|----------|--------|-------|
| Accent | 1800 Hz   | 3600 Hz  | 2 ms   | 40 ms |
| Normal | 1400 Hz   | 2800 Hz  | 1 ms   | 30 ms |

### Mechanical
| Beat   | Freq     | Attack | Decay |
|--------|----------|--------|-------|
| Accent | 2400 Hz  | 1 ms   | 15 ms |
| Normal | 2000 Hz  | 1 ms   | 12 ms |

### Digital
| Beat   | Freq     | Attack | Decay | Duration |
|--------|----------|--------|-------|----------|
| Accent | 880 Hz   | 5 ms   | 20 ms | 80 ms    |
| Normal | 440 Hz   | 5 ms   | 20 ms | 60 ms    |

## Keyboard Shortcuts

Khi metronome settings mở (`[M]`):
- `[S]` - Cycle sound types
- `[+/-]` - Adjust BPM
- `[1-4]` - Change time signature
- `[Space]` - Play/Pause
- `[M]` - Close settings

## Kết quả
✅ Âm thanh tự nhiên, dễ nghe hơn nhiều
✅ 3 lựa chọn sound cho preference khác nhau
✅ Attack/Decay envelope tạo âm percussive tự nhiên
✅ Harmonic mixing cho wood block sound chân thực
