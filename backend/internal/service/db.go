package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/flipslidersand/phone-guard/backend/internal/model"
)

// DBService: PostgreSQL 操作
type DBService struct {
	db *sql.DB
}

// NewDBService: データベース接続初期化
func NewDBService(databaseURL string) (*DBService, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DBService{db: db}, nil
}

// GetOrCreatePhoneNumber: 電話番号を取得または作成
func (ds *DBService) GetOrCreatePhoneNumber(ctx context.Context, number string) (*model.PhoneNumber, error) {
	phoneNum := &model.PhoneNumber{}

	err := ds.db.QueryRowContext(ctx,
		`INSERT INTO phone_numbers (id, number, report_count, first_seen_at, last_seen_at)
		 VALUES ($1, $2, 1, NOW(), NOW())
		 ON CONFLICT (number) DO UPDATE
		 SET report_count = report_count + 1, last_seen_at = EXCLUDED.last_seen_at
		 RETURNING id, number, report_count, first_seen_at, last_seen_at`,
		uuid.New().String(), number).Scan(
		&phoneNum.ID,
		&phoneNum.Number,
		&phoneNum.ReportCount,
		&phoneNum.FirstSeenAt,
		&phoneNum.LastSeenAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get/create phone number: %w", err)
	}

	return phoneNum, nil
}

// LogCall: 着信ログを記録
func (ds *DBService) LogCall(ctx context.Context, numberID, deviceID string, userReported bool) error {
	callID := uuid.New().String()
	_, err := ds.db.ExecContext(ctx,
		`INSERT INTO call_logs (id, number_id, device_id, called_at, user_reported)
		 VALUES ($1, $2, $3, NOW(), $4)`,
		callID, numberID, deviceID, userReported,
	)
	if err != nil {
		return fmt.Errorf("failed to log call: %w", err)
	}
	return nil
}

// GetWhitelist: デバイスのホワイトリスト取得
func (ds *DBService) GetWhitelist(ctx context.Context, deviceID string) ([]model.WhitelistEntry, error) {
	rows, err := ds.db.QueryContext(ctx,
		`SELECT id, device_id, number, label, created_at
		 FROM whitelist WHERE device_id = $1`,
		deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query whitelist: %w", err)
	}
	defer rows.Close()

	var whitelist []model.WhitelistEntry
	for rows.Next() {
		var entry model.WhitelistEntry
		if err := rows.Scan(&entry.ID, &entry.DeviceID, &entry.Number, &entry.Label, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan whitelist entry: %w", err)
		}
		whitelist = append(whitelist, entry)
	}

	return whitelist, nil
}

// AddToWhitelist: ホワイトリストに追加
func (ds *DBService) AddToWhitelist(ctx context.Context, deviceID, number, label string) error {
	entryID := uuid.New().String()
	_, err := ds.db.ExecContext(ctx,
		`INSERT INTO whitelist (id, device_id, number, label, created_at)
		 VALUES ($1, $2, $3, $4, NOW())`,
		entryID, deviceID, number, label,
	)
	if err != nil {
		return fmt.Errorf("failed to add to whitelist: %w", err)
	}
	return nil
}

// GetPhoneNumber: 電話番号の情報取得
func (ds *DBService) GetPhoneNumber(ctx context.Context, number string) (*model.PhoneNumber, error) {
	phoneNum := &model.PhoneNumber{}
	err := ds.db.QueryRowContext(ctx,
		`SELECT id, number, report_count, first_seen_at, last_seen_at
		 FROM phone_numbers WHERE number = $1`,
		number).Scan(
		&phoneNum.ID,
		&phoneNum.Number,
		&phoneNum.ReportCount,
		&phoneNum.FirstSeenAt,
		&phoneNum.LastSeenAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query phone number: %w", err)
	}

	return phoneNum, nil
}

// Close: データベース接続を閉じる
func (ds *DBService) Close() error {
	if ds.db != nil {
		return ds.db.Close()
	}
	return nil
}
