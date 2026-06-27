package service

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// LineNotifyService: LINE Notify API クライアント
type LineNotifyService struct {
	token string
}

// NewLineNotifyService: LINE Notify サービス初期化
func NewLineNotifyService(token string) *LineNotifyService {
	return &LineNotifyService{token: token}
}

// SendNotification: LINE Notify で通知を送信
func (lns *LineNotifyService) SendNotification(phoneNumber, label string, reportCount int) error {
	if lns.token == "" {
		return fmt.Errorf("LINE_NOTIFY_TOKEN not set")
	}

	message := fmt.Sprintf(
		"📞 知らない番号から着信\n"+
			"番号: %s\n"+
			"時刻: %s\n"+
			"全国記録: %d件\n\n"+
			"お母さんに確認してみてください。",
		phoneNumber,
		"14:32", // TODO: 実際の時刻に置き換え
		reportCount,
	)

	client := &http.Client{}
	endpoint := "https://notify-api.line.me/api/notify"

	data := url.Values{
		"message": {message},
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+lns.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LINE API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}
