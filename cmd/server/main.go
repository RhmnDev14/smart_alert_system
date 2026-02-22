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
	"smart_alert_system/internal/infrastructure/jobqueue"
	"smart_alert_system/internal/infrastructure/queue"
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
		log.Fatalf("❌ TELEGRAM_BOT_TOKEN is not set in .env file. Please add your Telegram Bot Token")
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

	// Initialize Global Task Producer & Consumer
	producer := queue.NewTaskProducer(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
	defer producer.Close()

	consumer := queue.NewTaskConsumer(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB, 5)

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
		producer,
		db,
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

	// Setup task processor & register routes
	taskProcessor := jobqueue.NewProcessor(schedulerUseCase)
	consumer.RegisterHandler(jobqueue.TaskMorningAlert, taskProcessor.HandleMorningAlert)
	consumer.RegisterHandler(jobqueue.TaskProcessMorningAlert, taskProcessor.HandleProcessMorningAlert)

	consumer.RegisterHandler(jobqueue.TaskEveningSummary, taskProcessor.HandleEveningSummary)
	consumer.RegisterHandler(jobqueue.TaskProcessEveningSummary, taskProcessor.HandleProcessEveningSummary)

	consumer.RegisterHandler(jobqueue.TaskActivityReminder, taskProcessor.HandleActivityReminder)
	consumer.RegisterHandler(jobqueue.TaskProcessActivityReminder, taskProcessor.HandleProcessActivityReminder)

	go func() {
		if err := consumer.Start(); err != nil {
			log.Fatalf("Task Consumer Server failed to start: %v", err)
		}
	}()

	// Setup scheduler
	location, err := cfg.GetLocation()
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	sched := scheduler.NewScheduler(producer, cfg.MorningAlertTime, cfg.EveningSummaryTime, location)
	if err := sched.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer sched.Stop()

	// Start Telegram long-polling in background
	pollingCtx, pollingCancel := context.WithCancel(context.Background())
	defer pollingCancel()
	go telegramHandler.StartLongPolling(pollingCtx)

	// Setup HTTP router
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

	go func() {
		log.Printf("🚀 HTTP Server starting on port %s", cfg.AppPort)
		log.Printf("🤖 Telegram Bot is running with long-polling mode")
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
	consumer.Stop()

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
