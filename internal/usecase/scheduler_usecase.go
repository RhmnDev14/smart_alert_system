package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/domain/repository"
	"smart_alert_system/internal/infrastructure/ai"
	"smart_alert_system/internal/infrastructure/whatsapp"
)

type SchedulerUseCase struct {
	userRepo       repository.UserRepository
	activityRepo   repository.ActivityRepository
	healthRepo     repository.HealthRepository
	alertRepo      repository.AlertRepository
	aiService      ai.AIService
	wahaClient     *whatsapp.WahaClient
}

func NewSchedulerUseCase(
	userRepo repository.UserRepository,
	activityRepo repository.ActivityRepository,
	healthRepo repository.HealthRepository,
	alertRepo repository.AlertRepository,
	aiService ai.AIService,
	wahaClient *whatsapp.WahaClient,
) *SchedulerUseCase {
	return &SchedulerUseCase{
		userRepo:     userRepo,
		activityRepo: activityRepo,
		healthRepo:   healthRepo,
		alertRepo:    alertRepo,
		aiService:    aiService,
		wahaClient:   wahaClient,
	}
}

func (uc *SchedulerUseCase) SendMorningAlerts(ctx context.Context) error {
	users, err := uc.userRepo.GetAllActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active users: %w", err)
	}

	for _, user := range users {
		if err := uc.sendMorningAlertForUser(ctx, user.ID, user.WhatsAppNumber); err != nil {
			log.Printf("Error sending morning alert to user %s: %v", user.ID, err)
			continue
		}
	}

	return nil
}

func (uc *SchedulerUseCase) sendMorningAlertForUser(ctx context.Context, userID uuid.UUID, whatsappNumber string) error {
	// Get today's activities
	activities, err := uc.activityRepo.GetTodayActivities(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get activities: %w", err)
	}

	// Get health profile
	healthProfile, _ := uc.healthRepo.GetHealthProfileByUserID(ctx, userID)

	// Generate alert message
	message, err := uc.aiService.GenerateMorningAlert(ctx, activities, healthProfile)
	if err != nil {
		message = uc.generateDefaultMorningAlert(activities)
	}

	// Create alert log
	alert := entity.NewAlertLog(userID, entity.AlertTypeMorning, message, time.Now())
	if err := uc.alertRepo.Create(ctx, alert); err != nil {
		return fmt.Errorf("failed to create alert log: %w", err)
	}

	// Send message
	if err := uc.wahaClient.SendMessage(whatsappNumber, message); err != nil {
		alert.MarkFailed(err)
		uc.alertRepo.Update(ctx, alert)
		return fmt.Errorf("failed to send message: %w", err)
	}

	alert.MarkSent()
	uc.alertRepo.Update(ctx, alert)

	return nil
}

func (uc *SchedulerUseCase) SendEveningSummaries(ctx context.Context) error {
	users, err := uc.userRepo.GetAllActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active users: %w", err)
	}

	for _, user := range users {
		if err := uc.sendEveningSummaryForUser(ctx, user.ID, user.WhatsAppNumber); err != nil {
			log.Printf("Error sending evening summary to user %s: %v", user.ID, err)
			continue
		}
	}

	return nil
}

func (uc *SchedulerUseCase) sendEveningSummaryForUser(ctx context.Context, userID uuid.UUID, whatsappNumber string) error {
	// Get completed activities today
	activities, err := uc.activityRepo.GetCompletedToday(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get completed activities: %w", err)
	}

	// Get health profile
	healthProfile, _ := uc.healthRepo.GetHealthProfileByUserID(ctx, userID)

	// Generate summary message
	message, err := uc.aiService.GenerateEveningSummary(ctx, activities, healthProfile)
	if err != nil {
		message = uc.generateDefaultEveningSummary(activities)
	}

	// Create alert log
	alert := entity.NewAlertLog(userID, entity.AlertTypeEvening, message, time.Now())
	if err := uc.alertRepo.Create(ctx, alert); err != nil {
		return fmt.Errorf("failed to create alert log: %w", err)
	}

	// Send message
	if err := uc.wahaClient.SendMessage(whatsappNumber, message); err != nil {
		alert.MarkFailed(err)
		uc.alertRepo.Update(ctx, alert)
		return fmt.Errorf("failed to send message: %w", err)
	}

	alert.MarkSent()
	uc.alertRepo.Update(ctx, alert)

	return nil
}

func (uc *SchedulerUseCase) generateDefaultMorningAlert(activities []*entity.Activity) string {
	if len(activities) == 0 {
		return "Selamat pagi! 🌅\n\nAnda tidak memiliki kegiatan yang dijadwalkan hari ini. Nikmati hari Anda!"
	}

	msg := "Selamat pagi! 🌅\n\nKegiatan hari ini:\n"
	for i, activity := range activities {
		msg += fmt.Sprintf("%d. %s - %s\n", i+1, activity.Title, activity.ScheduledTime.Format("15:04"))
	}
	msg += "\nSemoga hari Anda menyenangkan!"

	return msg
}

func (uc *SchedulerUseCase) generateDefaultEveningSummary(activities []*entity.Activity) string {
	if len(activities) == 0 {
		return "Selamat malam! 🌙\n\nAnda belum menyelesaikan kegiatan hari ini. Istirahat yang cukup untuk hari esok!"
	}

	msg := "Selamat malam! 🌙\n\nRingkasan hari ini:\n"
	for i, activity := range activities {
		msg += fmt.Sprintf("%d. ✓ %s\n", i+1, activity.Title)
	}
	msg += fmt.Sprintf("\nTotal: %d kegiatan selesai. Istirahat yang cukup!", len(activities))

	return msg
}

