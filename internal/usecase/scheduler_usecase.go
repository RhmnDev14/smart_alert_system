package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/domain/repository"
	"smart_alert_system/internal/infrastructure/ai"
	"smart_alert_system/internal/infrastructure/queue"
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
	producer       queue.TaskProducer
	txManager      repository.TransactionManager
}

func NewSchedulerUseCase(
	userRepo repository.UserRepository,
	activityRepo repository.ActivityRepository,
	healthRepo repository.HealthRepository,
	alertRepo repository.AlertRepository,
	aiService ai.AIService,
	telegramClient *telegram.TelegramClient,
	producer queue.TaskProducer,
	txManager repository.TransactionManager,
) *SchedulerUseCase {
	return &SchedulerUseCase{
		userRepo:       userRepo,
		activityRepo:   activityRepo,
		healthRepo:     healthRepo,
		alertRepo:      alertRepo,
		aiService:      aiService,
		telegramClient: telegramClient,
		producer:       producer,
		txManager:      txManager,
	}
}

func (uc *SchedulerUseCase) SendMorningAlerts(ctx context.Context) error {
	users, err := uc.userRepo.GetAllActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active users: %w", err)
	}

	for _, user := range users {
		payload := []byte(user.ID.String())
		if err := uc.producer.Publish("scheduler:process_morning_alert", payload, 1); err != nil {
			log.Printf("Failed to publish task to process morning alert for user %s: %v", user.ID, err)
		} else {
			log.Printf("Successfully published morning alert task for %s", user.ID)
		}
	}

	return nil
}

func (uc *SchedulerUseCase) ProcessSingleMorningAlert(ctx context.Context, userIDStr string) error {
	return uc.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return fmt.Errorf("invalid user id format: %w", err)
		}

		user, err := uc.userRepo.GetByID(txCtx, userID)
		if err != nil || user == nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Get today's activities
		allActivities, err := uc.activityRepo.GetTodayActivities(txCtx, userID)
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

		// Generate alert message (health profile no longer needed for morning alert)
		message, err := uc.aiService.GenerateMorningAlert(txCtx, activities, nil)
		if err != nil {
			message = uc.generateDefaultMorningAlert(activities)
		}

		// Create alert log
		alert := entity.NewAlertLog(userID, entity.AlertTypeMorning, message, time.Now())
		if err := uc.alertRepo.Create(txCtx, alert); err != nil {
			return fmt.Errorf("failed to create alert log: %w", err)
		}

		// Send message via Telegram
		if err := uc.telegramClient.SendMessageByStringID(user.WhatsAppNumber, message); err != nil {
			alert.MarkFailed(err)
			uc.alertRepo.Update(txCtx, alert)
			return fmt.Errorf("failed to send message: %w", err)
		}

		alert.MarkSent()
		uc.alertRepo.Update(txCtx, alert)

		return nil
	})
}

func (uc *SchedulerUseCase) SendEveningSummaries(ctx context.Context) error {
	users, err := uc.userRepo.GetAllActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active users: %w", err)
	}

	for _, user := range users {
		payload := []byte(user.ID.String())
		if err := uc.producer.Publish("scheduler:process_evening_summary", payload, 1); err != nil {
			log.Printf("Failed to publish task to process evening summary for user %s: %v", user.ID, err)
		} else {
			log.Printf("Successfully published evening summary task for %s", user.ID)
		}
	}

	return nil
}

func (uc *SchedulerUseCase) ProcessSingleEveningSummary(ctx context.Context, userIDStr string) error {
	return uc.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return fmt.Errorf("invalid user id format: %w", err)
		}

		user, err := uc.userRepo.GetByID(txCtx, userID)
		if err != nil || user == nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Get completed activities today
		activities, err := uc.activityRepo.GetCompletedToday(txCtx, userID)
		if err != nil {
			return fmt.Errorf("failed to get completed activities: %w", err)
		}

		// Generate summary message (health profile no longer needed)
		message, err := uc.aiService.GenerateEveningSummary(txCtx, activities, nil)
		if err != nil {
			message = uc.generateDefaultEveningSummary(activities)
		}

		// Create alert log
		alert := entity.NewAlertLog(userID, entity.AlertTypeEvening, message, time.Now())
		if err := uc.alertRepo.Create(txCtx, alert); err != nil {
			return fmt.Errorf("failed to create alert log: %w", err)
		}

		// Send message via Telegram
		if err := uc.telegramClient.SendMessageByStringID(user.WhatsAppNumber, message); err != nil {
			alert.MarkFailed(err)
			uc.alertRepo.Update(txCtx, alert)
			return fmt.Errorf("failed to send message: %w", err)
		}

		alert.MarkSent()
		uc.alertRepo.Update(txCtx, alert)

		return nil
	})
}

func (uc *SchedulerUseCase) SendActivityReminders(ctx context.Context) error {
	// We want to remind activities that are scheduled up to the current minute
	now := time.Now()

	activities, err := uc.activityRepo.GetPendingActivitiesToRemind(ctx, now)
	if err != nil {
		return fmt.Errorf("failed to get activities to remind: %w", err)
	}

	for _, activity := range activities {
		// Just enqueue the payload
		payload := []byte(activity.ID.String())
		// Assuming we will publish using "scheduler:process_activity_reminder"
		// Retries param is set to 2.
		if err := uc.producer.Publish("scheduler:process_activity_reminder", payload, 1); err != nil {
			log.Printf("Failed to publish task to process activity reminder %s: %v", activity.ID, err)
		} else {
			log.Printf("Successfully published activity reminder task for %s", activity.ID)
		}
	}

	return nil
}

func (uc *SchedulerUseCase) ProcessSingleActivityReminder(ctx context.Context, activityIDStr string) error {
	return uc.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		now := time.Now()
		activityID, err := uuid.Parse(activityIDStr)
		if err != nil {
			return fmt.Errorf("invalid activity id format: %w", err)
		}

		activity, err := uc.activityRepo.GetByID(txCtx, activityID)
		if err != nil || activity == nil {
			return fmt.Errorf("failed to get activity %s: %w", activityIDStr, err)
		}

		// Double check to ensure it hasn't been completed or reminded concurrently
		if activity.Status == entity.ActivityStatusCompleted || activity.ReminderTime != nil {
			log.Printf("Activity %s already processed or reminded. Skipping.", activityIDStr)
			return nil
		}

		user, err := uc.userRepo.GetByID(txCtx, activity.UserID)
		if err != nil || user == nil {
			return fmt.Errorf("failed to get user for activity %s: %w", activityID, err)
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
		message, err := uc.aiService.GenerateActivityReminder(txCtx, activity.Title, activity.Description, localScheduledTime.Format("15:04"))
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
			// Return error to trigger queue retry (up to 2 times, as defined during Publish)
			return err
		}

		// Update activity to mark it as reminded AND completed
		activity.ReminderTime = &now
		activity.Status = entity.ActivityStatusCompleted
		activity.CompletedAt = &now
		if err := uc.activityRepo.Update(txCtx, activity); err != nil {
			log.Printf("Failed to update status & reminder_time for activity %s: %v", activity.ID, err)
			return err
		} else {
			log.Printf("✓ Reminder sent and marked completed for activity %s", activity.Title)
		}

		return nil
	})
}

func (uc *SchedulerUseCase) generateDefaultMorningAlert(activities []*entity.Activity) string {
	if len(activities) == 0 {
		return "Selamat pagi! ☀️\n\nBelum ada jadwal kegiatan untuk hari ini.\nApakah ada kegiatan yang ingin kamu rencanakan hari ini? Beritahu saya ya! 😊\n\nSmart Alert System\nDevelop by Rahman Umardi"
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	msg := "Selamat pagi! ☀️\n\nBerikut jadwal kegiatan kamu hari ini:\n\n"
	for i, activity := range activities {
		localTime := activity.ScheduledTime.In(loc)
		msg += fmt.Sprintf("%d. %s — %s\n", i+1, activity.Title, localTime.Format("15:04"))
	}
	msg += "\nApakah ada kegiatan lain yang ingin ditambahkan untuk hari ini? 😊\n\nSmart Alert System\nDevelop by Rahman Umardi"

	return msg
}

func (uc *SchedulerUseCase) generateDefaultEveningSummary(activities []*entity.Activity) string {
	if len(activities) == 0 {
		return "Selamat malam! 🌙\n\nTidak ada kegiatan yang tercatat hari ini.\nIstirahat yang cukup ya! 😊\n\nSmart Alert System\nDevelop by Rahman Umardi"
	}

	msg := "Selamat malam! 🌙\n\nBerikut rangkuman kegiatan yang telah kamu selesaikan hari ini:\n\n"
	for i, activity := range activities {
		msg += fmt.Sprintf("%d. ✅ %s\n", i+1, activity.Title)
	}
	msg += fmt.Sprintf("\nTotal: %d kegiatan diselesaikan. \n\nSmart Alert System\nDevelop by Rahman Umardi", len(activities))

	return msg
}
