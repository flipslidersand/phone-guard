package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flipslidersand/phone-guard/backend/internal/model"
	"github.com/flipslidersand/phone-guard/backend/internal/service"
)

// CallHandler: POST /api/calls
type CallHandler struct {
	dbService    *service.DBService
	phoneService *service.PhoneService
}

// NewCallHandler: ハンドラー初期化
func NewCallHandler(dbService *service.DBService, phoneService *service.PhoneService) *CallHandler {
	return &CallHandler{
		dbService:    dbService,
		phoneService: phoneService,
	}
}

// ServeHTTP: POST /api/calls
func (h *CallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req model.CallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, model.CallResponse{Error: "Invalid request"})
		return
	}

	// Validate phone number
	if !h.phoneService.ValidateNumber(req.Number) {
		respondJSON(w, http.StatusBadRequest, model.CallResponse{Error: "Invalid phone number"})
		return
	}

	// Get or create phone number record
	phoneNum, err := h.dbService.GetOrCreatePhoneNumber(r.Context(), req.Number)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, model.CallResponse{
			Error: fmt.Sprintf("Database error: %v", err),
		})
		return
	}

	// Log call
	if err := h.dbService.LogCall(r.Context(), phoneNum.ID, req.DeviceID, false); err != nil {
		respondJSON(w, http.StatusInternalServerError, model.CallResponse{
			Error: fmt.Sprintf("Failed to log call: %v", err),
		})
		return
	}

	// Determine if suspicious
	isSuspicious := h.phoneService.IsSuspicious(phoneNum.ReportCount)

	respondJSON(w, http.StatusOK, model.CallResponse{
		Status:   "recorded",
		Count:    phoneNum.ReportCount,
		IsUnsafe: isSuspicious,
	})
}

// WhitelistHandler: GET/POST /api/whitelist
type WhitelistHandler struct {
	dbService *service.DBService
}

// NewWhitelistHandler: ハンドラー初期化
func NewWhitelistHandler(dbService *service.DBService) *WhitelistHandler {
	return &WhitelistHandler{dbService: dbService}
}

// ServeHTTP: GET/POST /api/whitelist
func (h *WhitelistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	deviceID := r.URL.Query().Get("deviceId")
	if deviceID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing deviceId"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetWhitelist(w, r, deviceID)
	case http.MethodPost:
		h.handlePostWhitelist(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WhitelistHandler) handleGetWhitelist(w http.ResponseWriter, r *http.Request, deviceID string) {
	whitelist, err := h.dbService.GetWhitelist(r.Context(), deviceID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"deviceId":  deviceID,
		"whitelist": whitelist,
		"count":     len(whitelist),
	})
}

func (h *WhitelistHandler) handlePostWhitelist(w http.ResponseWriter, r *http.Request) {
	var req model.WhitelistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if err := h.dbService.AddToWhitelist(r.Context(), req.DeviceID, req.Number, req.Label); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"status": "added"})
}

// NumberLookupHandler: GET /api/numbers?number=...
type NumberLookupHandler struct {
	dbService *service.DBService
}

// NewNumberLookupHandler: ハンドラー初期化
func NewNumberLookupHandler(dbService *service.DBService) *NumberLookupHandler {
	return &NumberLookupHandler{dbService: dbService}
}

// ServeHTTP: GET /api/numbers
func (h *NumberLookupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	number := r.URL.Query().Get("number")
	if number == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing number parameter"})
		return
	}

	phoneNum, err := h.dbService.GetPhoneNumber(r.Context(), number)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if phoneNum == nil {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"number": number,
			"count":  0,
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"number":     phoneNum.Number,
		"count":      phoneNum.ReportCount,
		"firstSeen":  phoneNum.FirstSeenAt,
		"lastSeen":   phoneNum.LastSeenAt,
	})
}

// NotifyHandler: POST /api/notify (for testing)
type NotifyHandler struct {
	lineService *service.LineNotifyService
}

// NewNotifyHandler: ハンドラー初期化
func NewNotifyHandler(lineService *service.LineNotifyService) *NotifyHandler {
	return &NotifyHandler{lineService: lineService}
}

// ServeHTTP: POST /api/notify
func (h *NotifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req model.NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if err := h.lineService.SendNotification(req.Number, "", req.ReportCount); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to send notification: %v", err),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

// Helper: JSON レスポンス
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
