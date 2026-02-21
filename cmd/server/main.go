package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smart_alert_system/internal/config"
	"smart_alert_system/internal/handler"
	"smart_alert_system/internal/infrastructure/ai"
	"smart_alert_system/internal/infrastructure/database"
	infraRepo "smart_alert_system/internal/infrastructure/repository"
	"smart_alert_system/internal/infrastructure/scheduler"
	"smart_alert_system/internal/infrastructure/telegram"
	"smart_alert_system/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate Telegram Bot Token
	if cfg.TelegramBotToken == "" {
		log.Fatalf("❌ TELEGRAM_BOT_TOKEN is not set in .env file. Please add your Telegram Bot Token (get it from @BotFather)")
	}

	// Connect to database
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✓ Connected to database")

	// Initialize repositories
	userRepo := infraRepo.NewUserRepository(db)
	activityRepo := infraRepo.NewActivityRepository(db)
	messageRepo := infraRepo.NewMessageRepository(db)
	alertRepo := infraRepo.NewAlertRepository(db)
	healthRepo := infraRepo.NewHealthRepository(db)
	categoryRepo := infraRepo.NewCategoryRepository(db)

	// Initialize Telegram Bot client
	telegramClient := telegram.NewTelegramClient(cfg.TelegramBotToken)

	// Verify bot token
	botInfo, err := telegramClient.GetMe()
	if err != nil {
		log.Fatalf("❌ Failed to verify Telegram Bot Token: %v", err)
	}
	log.Printf("✓ Telegram Bot verified: @%s", botInfo.Username)

	// Initialize AI Service
	var aiService ai.AIService

	if cfg.AIModel == "" {
		cfg.AIModel = "gpt-3.5-turbo"
	}

	log.Printf("✓ AI Service: %s", cfg.AIProvider)
	log.Printf("  Model: %s", cfg.AIModel)
	if cfg.AIBaseURL != "" {
		log.Printf("  Base URL: %s", cfg.AIBaseURL)
	}
	aiService = ai.NewOpenAIService(cfg.AIApiKey, cfg.AIModel, cfg.AIBaseURL)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	activityUseCase := usecase.NewActivityUseCase(activityRepo, userRepo, categoryRepo)
	schedulerUseCase := usecase.NewSchedulerUseCase(
		userRepo,
		activityRepo,
		healthRepo,
		alertRepo,
		aiService,
		telegramClient,
	)

	// Initialize Telegram handler
	telegramHandler := handler.NewTelegramHandler(
		userUseCase,
		activityUseCase,
		aiService,
		telegramClient,
		messageRepo,
		alertRepo,
	)

	// Setup scheduler
	location, err := cfg.GetLocation()
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	sched := scheduler.NewScheduler(schedulerUseCase, cfg.MorningAlertTime, cfg.EveningSummaryTime, location)
	if err := sched.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer sched.Stop()

	// Start Telegram long-polling in background
	pollingCtx, pollingCancel := context.WithCancel(context.Background())
	defer pollingCancel()
	go telegramHandler.StartLongPolling(pollingCtx)

	// Setup HTTP router (keep webhook endpoint as fallback + health check)
	router := mux.NewRouter()
	router.HandleFunc("/webhook", telegramHandler.HandleWebhook).Methods("POST")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Start HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.AppPort)
		log.Printf("🤖 Telegram Bot is running with long-polling mode")
		log.Printf("📡 Webhook endpoint available at POST /webhook")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	pollingCancel()

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
