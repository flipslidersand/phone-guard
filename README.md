# Phone Guard — 固定電話詐欺対策デバイス

固定電話の回線に物理割り込みして、知らない番号から着信があったら家族の LINE に自動通知するIoTデバイス。

- **録音しない** — 音声はサーバに送らない
- **判定しない** — AIで詐欺か判断しない。家族が判断する
- **番号だけ共有** — ユーザーが任意でDB に提供

## 技術スタック

- **Embedded**: ESP32 + FSK デコード（C++, PlatformIO）
- **Backend**: Go (Cloud Run)
- **DB**: PostgreSQL (Cloud SQL)
- **Notification**: LINE Notify API
- **Deployment**: GCP

## ハードウェア構成

```
固定電話回線
    ↓
[分岐モジュラージャック]
    ↓          ↓
固定電話機   ESP32
    ├── FSKデコード
    ├── Wi-Fi
    └── LINE通知
```

部品: ESP32(1k) + モジュラー分岐(500) + アナログ回路(500) + ケース(500) = 約3,000円

## フェーズ設計

### Phase 1 (Firmware MVP)
- ESP32 単体で FSK デコード
- ローカル電話帳照合
- LINE Notify 通知

### Phase 2 (Cloud Integration)
- 番号DB と照合（全国記録件数表示）

### Phase 3 (Smart Features)
- 高スコア番号の自動切断
- 統計分析

### Phase 4 (Production)
- 基板量産、ECで販売

## Project Structure

```
phone-guard/
├── esp32/            # Firmware (C++ / PlatformIO)
│   ├── src/
│   │   └── main.cpp
│   ├── lib/
│   └── platformio.ini
├── backend/          # Go API
│   ├── cmd/api/
│   ├── internal/
│   └── go.mod
├── db/
│   └── schema.sql
└── docs/
    ├── fsk-spec.md
    └── wiring.md
```

## Getting Started

### Firmware (Phase 1)

```bash
cd esp32
pio run -t upload  # Upload to ESP32
pio device monitor  # Serial monitor
```

### Backend (Phase 2)

```bash
cd backend
go mod download
go run cmd/api/main.go
```

### Database (PostgreSQL)

```bash
psql -U postgres -f db/schema.sql
```

## API Endpoints

- `GET /api/numbers/:number` — 番号情報照合
- `POST /api/calls` — 着信ログ記録
- `GET /api/whitelist` — ホワイトリスト取得
- `POST /api/whitelist` — ホワイトリスト追加

## License

MIT
