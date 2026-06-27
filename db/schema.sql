-- Phone numbers: 全国から報告された番号
CREATE TABLE IF NOT EXISTS phone_numbers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  number VARCHAR(20) NOT NULL UNIQUE,
  report_count INTEGER DEFAULT 1,
  first_seen_at TIMESTAMP DEFAULT NOW(),
  last_seen_at TIMESTAMP DEFAULT NOW(),

  INDEX idx_number (number),
  INDEX idx_report_count (report_count)
);

-- Call logs: デバイスごとの着信履歴
CREATE TABLE IF NOT EXISTS call_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  number_id UUID REFERENCES phone_numbers(id) ON DELETE CASCADE,
  device_id VARCHAR(64),  -- デバイス識別（個人特定なし）
  called_at TIMESTAMP DEFAULT NOW(),
  user_reported BOOLEAN DEFAULT FALSE,

  INDEX idx_device_id (device_id),
  INDEX idx_called_at (called_at)
);

-- Whitelist: 家族・かかりつけ医等の登録済み番号
CREATE TABLE IF NOT EXISTS whitelist (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_id VARCHAR(64),
  number VARCHAR(20),
  label VARCHAR(100),  -- 「長男」「かかりつけ医」等
  created_at TIMESTAMP DEFAULT NOW(),

  INDEX idx_device_id (device_id),
  UNIQUE KEY unique_device_number (device_id, number)
);
