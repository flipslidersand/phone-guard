package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/flipslidersand/phone-guard/backend/internal/handler"
	"github.com/flipslidersand/phone-guard/backend/internal/service"
)

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	lineNotifyToken := os.Getenv("LINE_NOTIFY_TOKEN")
	if lineNotifyToken == "" {
		log.Println("Warning: LINE_NOTIFY_TOKEN not set")
	}

	log.Printf("Starting Phone Guard API on port %s...\n", port)

	// 1. PostgreSQL 接続
	dbService, err := service.NewDBService(databaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbService.Close()

	// 2. LINE Notify サービス初期化
	lineService := service.NewLineNotifyService(lineNotifyToken)

	// 3. ビジネスロジックサービス
	phoneService := service.NewPhoneService()

	// 4. HTTP ハンドラー登録
	mux := http.NewServeMux()

	callHandler := handler.NewCallHandler(dbService, phoneService)
	whitelistHandler := handler.NewWhitelistHandler(dbService)
	notifyHandler := handler.NewNotifyHandler(lineService)
	numberLookupHandler := handler.NewNumberLookupHandler(dbService)

	mux.Handle("/api/calls", callHandler)
	mux.Handle("/api/whitelist", whitelistHandler)
	mux.Handle("/api/notify", notifyHandler)
	mux.Handle("/api/numbers", numberLookupHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	// 5. サーバー起動
	server := &http.Server{
		Addr:         net.JoinHostPort("", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("Server listening on %s", server.Addr)

	// Graceful shutdown 待機
	<-sigChan
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
