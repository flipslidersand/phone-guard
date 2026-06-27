# FSK Signal Specification (ナンバーディスプレイ)

## Overview

ナンバーディスプレイは **着信1〜2回鳴動後の無音期間** に、**1200bps FSK** で発信者番号を乗せる。

## Signal Sequence

```
着信開始（1回目の鳴動）
  ↓ ~1秒
無音期間（~500ms）  ← FSK信号はここ
  ↓
着信鳴動（2回目）
```

## FSK Parameters

| Parameter | Value |
|-----------|-------|
| **Baud Rate** | 1200 bps |
| **Frequency Mark** | 1300 Hz |
| **Frequency Space** | 2100 Hz |
| **Modulation** | AFSK (Audio FSK) |
| **Data Format** | DTMF-like sequence (digits + control codes) |

## Decoding Steps

1. **Capture**: ESP32 ADC でアナログ信号をサンプリング（8kHz以上推奨）
2. **Detect**: マーク周波数（1300Hz）と スペース周波数（2100Hz）を Goertzel アルゴリズムで検出
3. **Demodulate**: FSK 信号から 1/0 ビット列を復調
4. **Decode**: 発信者番号（10-11桁）を抽出

## Reference

- Arduino FSK Decoder Library (GitHub): `arduino-fsk-decoder`
- ESP32 ADC Documentation: https://docs.espressif.com/projects/esp-idf/
