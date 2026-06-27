package service

// PhoneService: 電話番号のバリデーション・フォーマット
type PhoneService struct{}

// NewPhoneService: サービス初期化
func NewPhoneService() *PhoneService {
	return &PhoneService{}
}

// ValidateNumber: 電話番号の検証
// 日本の固定電話・携帯電話形式をチェック
func (ps *PhoneService) ValidateNumber(number string) bool {
	if len(number) < 10 || len(number) > 11 {
		return false
	}

	// Simple check: starts with 0 and all digits
	if number[0] != '0' {
		return false
	}

	for _, c := range number {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

// FormatNumber: 電話番号をフォーマット
// "09012345678" → "090-1234-5678"
func (ps *PhoneService) FormatNumber(number string) string {
	if len(number) == 10 {
		// Fixed line: 0X-XXXX-XXXX
		return number[0:2] + "-" + number[2:6] + "-" + number[6:]
	} else if len(number) == 11 {
		// Mobile: 0XX-XXXX-XXXX
		return number[0:3] + "-" + number[3:7] + "-" + number[7:]
	}
	return number
}

// IsSuspicious: 危険度判定（簡易版）
// Phase 2 で LLM ベースに拡張
func (ps *PhoneService) IsSuspicious(reportCount int) bool {
	return reportCount >= 5
}
