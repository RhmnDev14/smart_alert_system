package usecase

import (
	"context"
	"fmt"
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
