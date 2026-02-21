package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/domain/repository"
	"smart_alert_system/internal/infrastructure/ai"
	"smart_alert_system/internal/infrastructure/telegram"

	"github.com/google/uuid"
)

type SchedulerUseCase struct {
	userRepo       repository.UserRepository
	activityRepo   repository.ActivityRepository
	healthRepo     repository.HealthRepository
	alertRepo      repository.AlertRepository
	aiService      ai.AIService
	telegramClient *telegram.TelegramClient
}

func NewSchedulerUseCase(
	userRepo repository.UserRepository,
	activityRepo repository.ActivityRepository,
	healthRepo repository.HealthRepository,
	alertRepo repository.AlertRepository,
	aiService ai.AIService,
	telegramClient *telegram.TelegramClient,
) *SchedulerUseCase {
	return &SchedulerUseCase{
		userRepo:       userRepo,
		activityRepo:   activityRepo,
		healthRepo:     healthRepo,
		alertRepo:      alertRepo,
		aiService:      aiService,
		telegramClient: telegramClient,
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

func (uc *SchedulerUseCase) sendMorningAlertForUser(ctx context.Context, userID uuid.UUID, chatID string) error {
	// Get today's activities
	allActivities, err := uc.activityRepo.GetTodayActivities(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get activities: %w", err)
	}

	// Filter only pending activities for morning alert
	var activities []*entity.Activity
	for _, act := range allActivities {
		if act.Status == entity.ActivityStatusPending {
			activities = append(activities, act)
		}
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

	// Send message via Telegram
	if err := uc.telegramClient.SendMessageByStringID(chatID, message); err != nil {
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

func (uc *SchedulerUseCase) sendEveningSummaryForUser(ctx context.Context, userID uuid.UUID, chatID string) error {
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

	// Send message via Telegram
	if err := uc.telegramClient.SendMessageByStringID(chatID, message); err != nil {
		alert.MarkFailed(err)
		uc.alertRepo.Update(ctx, alert)
		return fmt.Errorf("failed to send message: %w", err)
	}

	alert.MarkSent()
	uc.alertRepo.Update(ctx, alert)

	return nil
}

func (uc *SchedulerUseCase) SendActivityReminders(ctx context.Context) error {
	// We want to remind activities that are scheduled up to the current minute
	now := time.Now()

	activities, err := uc.activityRepo.GetPendingActivitiesToRemind(ctx, now)
	if err != nil {
		return fmt.Errorf("failed to get activities to remind: %w", err)
	}

	for _, activity := range activities {
		// Get user to send message
		user, err := uc.userRepo.GetByID(ctx, activity.UserID)
		if err != nil || user == nil {
			log.Printf("Error getting user for activity %s: %v", activity.ID, err)
			continue
		}

		// Adjust time to user's timezone if available
		loc, _ := time.LoadLocation("Asia/Jakarta")
		if user.Timezone != "" {
			if l, err := time.LoadLocation(user.Timezone); err == nil {
				loc = l
			}
		}
		localScheduledTime := activity.ScheduledTime.In(loc)

		// Generate reminder message dynamically with AI
		message, err := uc.aiService.GenerateActivityReminder(ctx, activity.Title, activity.Description, localScheduledTime.Format("15:04"))
		if err != nil {
			log.Printf("⚠️ AI API failed to generate reminder for %s, using fallback", activity.Title)
			// Fallback if AI fails
			message = fmt.Sprintf("⏰ *PENGINGAT KEGIATAN*\n\nSudah waktunya untuk kegiatan:\n*❖ %s*\n\nWaktu: %s\n",
				activity.Title, localScheduledTime.Format("15:04"))

			if activity.Description != "" {
				message += fmt.Sprintf("Catatan: %s\n", activity.Description)
			}
		}

		// Send message via Telegram
		if err := uc.telegramClient.SendMessageByStringID(user.WhatsAppNumber, message); err != nil {
			log.Printf("Failed to send reminder for activity %s (User %s): %v", activity.ID, user.WhatsAppNumber, err)
			// Proceed to update reminder_time anyway so it doesn't try again and spam the logs/API
		}

		// Update activity to mark it as reminded AND completed
		activity.ReminderTime = &now
		activity.Status = entity.ActivityStatusCompleted
		activity.CompletedAt = &now
		if err := uc.activityRepo.Update(ctx, activity); err != nil {
			log.Printf("Failed to update status & reminder_time for activity %s: %v", activity.ID, err)
		} else {
			log.Printf("✓ Reminder processed and marked completed for activity %s", activity.Title)
		}
	}

	return nil
}

func (uc *SchedulerUseCase) generateDefaultMorningAlert(activities []*entity.Activity) string {
	if len(activities) == 0 {
		return "Selamat pagi! 🌅\n\nSemoga harimu dipenuhi dengan kebahagiaan dan produktivitas. \"Hari ini adalah kesempatan baru untuk menjadi lebih baik.\"\n\nSaat ini Anda belum memiliki kegiatan yang dijadwalkan hari ini. Adakah rencana atau aktivitas yang ingin Anda lakukan hari ini? Beritahu saya, agar saya bisa mengingatkan Anda!\n\nSmart Alert System\nDevelop by Rahman Umardi"
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	msg := "Selamat pagi! 🌅\n\nSemoga harimu penuh semangat. \"Setiap langkah kecil yang kamu ambil hari ini, akan membawamu lebih dekat ke tujuanmu.\"\n\nBerikut adalah jadwal kegiatan Anda hari ini:\n"
	for i, activity := range activities {
		localTime := activity.ScheduledTime.In(loc)
		msg += fmt.Sprintf("%d. %s - %s\n", i+1, activity.Title, localTime.Format("15:04"))
	}
	msg += "\nApakah ada kegiatan tambahan lain yang ingin Anda rencanakan hari ini? Beritahu saya 😊\n\nSmart Alert System\nDevelop by Rahman Umardi"

	return msg
}

func (uc *SchedulerUseCase) generateDefaultEveningSummary(activities []*entity.Activity) string {
	if len(activities) == 0 {
		return "Selamat malam! 🌙\n\nHari ini Anda tidak merekam kegiatan yang diselesaikan. Mari bersiap dan rencanakan aktivitas yang baik untuk esok hari.\n\nRekomendasi: Cobalah tidur sedikit lebih awal malam ini untuk menjaga kebugaran tubuh esok pagi. Istirahat yang cukup sangat penting!\n\nSmart Alert System\nDevelop by Rahman Umardi"
	}

	msg := "Selamat malam! 🌙\n\nWah, Anda telah melewati hari ini dengan luar biasa! Berikut adalah ringkasan kegiatan yang telah Anda selesaikan hari ini:\n\n"
	for i, activity := range activities {
		msg += fmt.Sprintf("%d. ✓ %s\n", i+1, activity.Title)
	}

	msg += fmt.Sprintf("\nTotal: %d kegiatan diselesaikan. Bagus sekali!\n\nAnalisis & Rekomendasi:\nAnda memiliki tingkat produktivitas yang terpantau aktif hari ini. Pastikan Anda merilekskan otot-otot dan menjauhi layar sebelum tidur demi kualitas istirahat yang optimal.\n\nSelamat beristirahat dengan damai!\n\nSmart Alert System\nDevelop by Rahman Umardi", len(activities))

	return msg
}
