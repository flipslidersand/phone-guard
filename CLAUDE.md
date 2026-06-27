# phone-guard CLAUDE.md

## Project Overview

ESP32 + FSK デコード + LINE 通知で、知らない番号からの着信を家族に通知する IοT デバイス。

**設計思想**: 録音しない。判定しない。番号だけ検知して通知。

## Tech Stack

| Layer | Tech | Purpose |
|-------|------|---------|
| **Embedded** | ESP32 + C++ | FSK デコード、LINE 通知 |
| **Backend** | Go (Cloud Run) | 番号DB 照合、統計 |
| **DB** | PostgreSQL | 番号スコアテーブル |
| **Notification** | LINE Notify API | 家族への通知 |

## Phase 1: Firmware MVP

### Key Implementation

**FSK デコード**:
- ナンバーディスプレイは着信1〜2回鳴動後の無音期間に 1200bps FSK 信号で番号が乗ってくる
- ESP32 ADC でサンプリング → ソフトウェア FSK デコード → 発信者番号取得
- GitHub 先行実装を参考に実装

**ローカル電話帳キャッシュ**:
- EEPROM に登録済み番号を保存
- 未登録番号のみ通知

**LINE Notify Integration**:
- デバイス→LINE 一方向通知
- 自動取得トークンで認証

### Flow

```
着信（1回目鳴動）
  ↓
無音期間（500ms）でFSK信号をADCでキャプチャ
  ↓
ソフトウェアFSKデコード → 発信者番号取得
  ↓
ローカル電話帳照合
  ↓
未登録 → LINE Notify で家族に通知
登録済み → スルー
```

## Phase 2: Cloud Integration

### Database Schema

```sql
CREATE TABLE phone_numbers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  number VARCHAR(20) NOT NULL UNIQUE,
  report_count INTEGER DEFAULT 1,
  first_seen_at TIMESTAMP DEFAULT NOW(),
  last_seen_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE call_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  number_id UUID REFERENCES phone_numbers(id),
  device_id VARCHAR(64),
  called_at TIMESTAMP DEFAULT NOW(),
  user_reported BOOLEAN DEFAULT FALSE
);

CREATE TABLE whitelist (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_id VARCHAR(64),
  number VARCHAR(20),
  label VARCHAR(100),  -- 「長男」「かかりつけ医」等
  created_at TIMESTAMP DEFAULT NOW()
);
```

### API Integration

- デバイスから番号を POST → DB 照合
- 「この番号は全国で〇〇件記録あり」を LINE 通知に追加

## Development Notes

- **FSK参考実装**: GitHub を検索（Arduino/ESP32 向け）
- **LINE Notify トークン**: 環境変数 `LINE_NOTIFY_TOKEN` で管理
- **Device ID**: デバイス固有の MAC アドレスなど（個人特定なし）

## Links

- **Wiring Diagram**: docs/wiring.md
- **FSK Specification**: docs/fsk-spec.md
- **GitHub**: github.com/flipslidersand/phone-guard (TBD)
