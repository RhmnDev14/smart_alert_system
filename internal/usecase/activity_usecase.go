package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/domain/repository"

	"github.com/google/uuid"
)

type ActivityUseCase struct {
	activityRepo repository.ActivityRepository
	userRepo     repository.UserRepository
	categoryRepo repository.CategoryRepository
}

func NewActivityUseCase(
	activityRepo repository.ActivityRepository,
	userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository,
) *ActivityUseCase {
	return &ActivityUseCase{
		activityRepo: activityRepo,
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
	}
}

func (uc *ActivityUseCase) CreateActivity(ctx context.Context, userID uuid.UUID, data entity.ActivityIntentData) (*entity.Activity, error) {
	// Validate user exists
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Set default scheduled time if not provided
	scheduledTime := time.Now().Add(1 * time.Hour)
	if data.ScheduledTime != nil {
		scheduledTime = *data.ScheduledTime
	}

	// Set default priority
	priority := data.Priority
	if priority == 0 {
		priority = 3
	}

	activity := entity.NewActivity(userID, data.Title, data.Description, scheduledTime, priority)
	if data.CategoryID != nil {
		activity.CategoryID = data.CategoryID
	}

	if err := uc.activityRepo.Create(ctx, activity); err != nil {
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	return activity, nil
}

// GetActivities retrieves activities for a user based on query filters.
// Data is always scoped to the user's own activities only.
func (uc *ActivityUseCase) GetActivities(ctx context.Context, userID uuid.UUID, query *entity.QueryData) ([]*entity.Activity, error) {
	if query == nil {
		return uc.activityRepo.GetByUserID(ctx, userID)
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	switch query.FilterType {
	case "today":
		return uc.activityRepo.GetByUserIDAndDate(ctx, userID, now)

	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		return uc.activityRepo.GetByUserIDAndDate(ctx, userID, tomorrow)

	case "date":
		if query.Date != "" {
			// Try parsing date in various formats
			var parsedDate time.Time
			var parseErr error
			formats := []string{"2006-01-02", "2006-01-02T15:04:05-07:00", "2006-01-02T15:04:05Z"}
			for _, f := range formats {
				parsedDate, parseErr = time.Parse(f, query.Date)
				if parseErr == nil {
					break
				}
			}
			if parseErr != nil {
				return nil, fmt.Errorf("invalid date format: %s", query.Date)
			}
			return uc.activityRepo.GetByUserIDAndDate(ctx, userID, parsedDate)
		}
		return uc.activityRepo.GetByUserID(ctx, userID)

	case "status":
		if query.Status != "" {
			status := entity.ActivityStatus(strings.ToLower(query.Status))
			return uc.activityRepo.GetByUserIDAndStatus(ctx, userID, status)
		}
		return uc.activityRepo.GetByUserID(ctx, userID)

	case "search":
		if query.Keyword != "" {
			return uc.activityRepo.SearchByKeyword(ctx, userID, query.Keyword)
		}
		return uc.activityRepo.GetByUserID(ctx, userID)

	case "all":
		return uc.activityRepo.GetByUserID(ctx, userID)

	default:
		return uc.activityRepo.GetByUserID(ctx, userID)
	}
}

// UpdateActivity finds and updates an activity belonging to the user.
// It searches by title (fuzzy match) and applies the requested changes.
// SECURITY: Only activities owned by the user can be updated.
func (uc *ActivityUseCase) UpdateActivity(ctx context.Context, userID uuid.UUID, data *entity.UpdateData) (*entity.Activity, error) {
	if data == nil || data.SearchTitle == "" {
		return nil, fmt.Errorf("search_title is required to identify the activity to update")
	}

	// Find activity by fuzzy title match, scoped to user
	activities, err := uc.activityRepo.GetByUserIDAndTitleFuzzy(ctx, userID, data.SearchTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to search activity: %w", err)
	}
	if len(activities) == 0 {
		return nil, fmt.Errorf("tidak ditemukan kegiatan dengan judul '%s'", data.SearchTitle)
	}

	// Use the first (most recent) match
	activity := activities[0]

	// Double-check ownership (defense in depth)
	if activity.UserID != userID {
		return nil, fmt.Errorf("anda tidak memiliki akses ke kegiatan ini")
	}

	// Apply updates
	if data.NewTitle != nil && *data.NewTitle != "" {
		activity.Title = *data.NewTitle
	}
	if data.NewDescription != nil {
		activity.Description = *data.NewDescription
	}
	if data.NewScheduledTime != nil && *data.NewScheduledTime != "" {
		t, err := time.Parse(time.RFC3339, *data.NewScheduledTime)
		if err == nil {
			activity.ScheduledTime = t
		}
	}
	if data.NewStatus != nil && *data.NewStatus != "" {
		newStatus := entity.ActivityStatus(strings.ToLower(*data.NewStatus))
		switch newStatus {
		case entity.ActivityStatusCompleted:
			activity.Complete()
		case entity.ActivityStatusCancelled:
			activity.Cancel()
		case entity.ActivityStatusPending:
			activity.Status = entity.ActivityStatusPending
			activity.UpdatedAt = time.Now()
		}
	}
	if data.NewPriority != nil && *data.NewPriority > 0 {
		activity.Priority = *data.NewPriority
		activity.UpdatedAt = time.Now()
	}

	activity.UpdatedAt = time.Now()

	if err := uc.activityRepo.Update(ctx, activity); err != nil {
		return nil, fmt.Errorf("failed to update activity: %w", err)
	}

	return activity, nil
}
