package model

import "time"

// CallRequest: ESP32 から送信される着信ログ
type CallRequest struct {
	Number    string `json:"number"`
	DeviceID  string `json:"deviceId"`
	Timestamp int64  `json:"timestamp"`
}

// CallResponse: 着信ログ記録のレスポンス
type CallResponse struct {
	Status   string `json:"status"`
	Count    int    `json:"count"`
	IsUnsafe bool   `json:"isUnsafe"`
	Error    string `json:"error,omitempty"`
}

// WhitelistEntry: ホワイトリスト登録エントリ
type WhitelistEntry struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"deviceId"`
	Number    string    `json:"number"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"createdAt"`
}

// WhitelistRequest: ホワイトリスト追加リクエスト
type WhitelistRequest struct {
	Number   string `json:"number"`
	Label    string `json:"label"`
	DeviceID string `json:"deviceId"`
}

// NotifyRequest: LINE 通知リクエスト
type NotifyRequest struct {
	Number    string `json:"number"`
	DeviceID  string `json:"deviceId"`
	ReportCount int `json:"reportCount"`
}

// PhoneNumber: 電話番号テーブルのモデル
type PhoneNumber struct {
	ID         string
	Number     string
	ReportCount int
	FirstSeenAt time.Time
	LastSeenAt  time.Time
}

// CallLog: 着信ログのモデル
type CallLog struct {
	ID         string
	NumberID   string
	DeviceID   string
	CalledAt   time.Time
	UserReported bool
}
