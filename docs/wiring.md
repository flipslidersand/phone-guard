# Wiring Diagram — Fixed Line Phone + ESP32

## Overview

固定電話回線を分岐して、ESP32 でナンバーディスプレイ信号を取得。

## Components

- ESP32 開発ボード（1,000円）
- モジュラー分岐コネクタ（RJ-11）（500円）
- オペアンプ（TL072 等）+ アナログ処理回路（500円）
- ケース（500円）
- **合計: 約 2,500円**

## Wiring Diagram

```
固定電話回線
  ├── 赤線（RX+）
  ├── 緑線（RX-）
  ├── 黄線（TX+）  
  └── 黒線（TX-）
       ↓
  [分岐モジュラージャック]
       ↙          ↘
  既存の           [アナログフロントエンド]
  固定電話機        (オペアンプ + フィルタ)
                       ↓
                  ESP32 ADC Pin 34
```

## Analog Frontend

ナンバーディスプレイ信号（1200bps FSK）は低レベル（-20dBm 程度）なため、
オペアンプで増幅 + フィルタリング後に ADC に入力。

1. **Low-Pass Filter**: 3500 Hz カットオフ（ナイキスト周波数 2100 Hz 以下を通す）
2. **Amplifier**: ゲイン 20-30 dB（TL072 or LM358）
3. **Bias**: 1.65V DC バイアス（ADC input 0.0V - 3.3V）

### Simple Op-Amp Circuit

```
固定電話 RX（AC信号）
  ├── 0.1µF Capacitor (AC coupling)
  ├── 1kΩ Resistor
  └── TL072 Non-inverting amp
       ├── Gain = 1 + R2/R1 (20x推奨)
       ├── +Power: 5V / -Power: GND
       ├── Output → 1.65V DC Bias
       └── → ESP32 ADC Pin 34
```

## Testing

1. 通常の着信でアナログ信号が波形として見えるか確認
2. ロジックアナライザで FSK デコード可能か確認
3. GitHub 実装で 10-20 回の着信テスト

## Safety Notes

- 固定電話回線は高電圧（-48V DC）を含む
- モジュラー分岐は **RX 線（低電圧）のみ** を使用
- TX 線（高電圧）には接触しない
