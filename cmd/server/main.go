package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"smart_alert_system/internal/config"
	"smart_alert_system/internal/handler"
	"smart_alert_system/internal/infrastructure/ai"
	"smart_alert_system/internal/infrastructure/database"
	infraRepo "smart_alert_system/internal/infrastructure/repository"
	"smart_alert_system/internal/infrastructure/scheduler"
	"smart_alert_system/internal/infrastructure/whatsapp"
	"smart_alert_system/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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

	// Initialize infrastructure services
	wahaClient := whatsapp.NewWahaClient(cfg.WahaServerURL, cfg.WahaAPIKey)

	// Initialize AI Service
	var aiService ai.AIService

	if cfg.AIProvider == "ollama" || cfg.AIBaseURL != "" {
		// Using Ollama (free, local)
		if cfg.AIBaseURL == "" {
			cfg.AIBaseURL = "http://localhost:11434/v1"
		}
		if cfg.AIModel == "" {
			cfg.AIModel = "llama3.2" // Default Ollama model
		}
		log.Printf("✓ AI Service: Ollama (Free)")
		log.Printf("  Base URL: %s", cfg.AIBaseURL)
		log.Printf("  Model: %s", cfg.AIModel)
		aiService = ai.NewOpenAIService("", cfg.AIModel, cfg.AIBaseURL)
	} else {
		// Using OpenAI or other provider
		if cfg.AIApiKey == "" {
			log.Fatalf("❌ AI_API_KEY is not set in .env file. Please add your OpenAI API key, or use Ollama by setting AI_PROVIDER=ollama")
		}
		if cfg.AIModel == "" {
			log.Printf("⚠️  AI_MODEL not set, using default: gpt-3.5-turbo")
			cfg.AIModel = "gpt-3.5-turbo"
		}

		// Normalize model name - ensure it has 'gpt-' prefix if it's a version number
		normalizedModel := normalizeModelName(cfg.AIModel)
		if normalizedModel != cfg.AIModel {
			log.Printf("⚠️  Model name normalized: %s -> %s", cfg.AIModel, normalizedModel)
			cfg.AIModel = normalizedModel
		}

		log.Printf("✓ AI Service: OpenAI")
		log.Printf("  Model: %s", cfg.AIModel)
		aiService = ai.NewOpenAIService(cfg.AIApiKey, cfg.AIModel, cfg.AIBaseURL)
	}

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	activityUseCase := usecase.NewActivityUseCase(activityRepo, userRepo, categoryRepo)
	schedulerUseCase := usecase.NewSchedulerUseCase(
		userRepo,
		activityRepo,
		healthRepo,
		alertRepo,
		aiService,
		wahaClient,
	)

	// Initialize handlers
	whatsappHandler := handler.NewWhatsAppHandler(
		userUseCase,
		activityUseCase,
		aiService,
		wahaClient,
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

	// Setup HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/webhook", whatsappHandler.HandleWebhook).Methods("POST")
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
		log.Printf("Server starting on port %s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// normalizeModelName ensures the model name has the correct format
// If user provides "3.5-turbo", it will be normalized to "gpt-3.5-turbo"
func normalizeModelName(model string) string {
	if model == "" {
		return "gpt-3.5-turbo"
	}

	// If model doesn't start with "gpt-" and looks like a version number, add prefix
	if !strings.HasPrefix(model, "gpt-") && !strings.HasPrefix(model, "o1-") {
		// Check if it looks like a model version (e.g., "3.5-turbo", "4", "4-turbo")
		if strings.HasPrefix(model, "3.5") || strings.HasPrefix(model, "4") {
			return "gpt-" + model
		}
	}

	return model
}
