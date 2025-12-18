package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/domain/repository"
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

func (uc *ActivityUseCase) GetUserActivities(ctx context.Context, userID uuid.UUID) ([]*entity.Activity, error) {
	return uc.activityRepo.GetByUserID(ctx, userID)
}

func (uc *ActivityUseCase) GetTodayActivities(ctx context.Context, userID uuid.UUID) ([]*entity.Activity, error) {
	return uc.activityRepo.GetTodayActivities(ctx, userID)
}

func (uc *ActivityUseCase) GetCompletedToday(ctx context.Context, userID uuid.UUID) ([]*entity.Activity, error) {
	return uc.activityRepo.GetCompletedToday(ctx, userID)
}

func (uc *ActivityUseCase) UpdateActivity(ctx context.Context, activityID uuid.UUID, data entity.UpdateActivityIntentData) error {
	activity, err := uc.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return fmt.Errorf("failed to get activity: %w", err)
	}
	if activity == nil {
		return fmt.Errorf("activity not found")
	}

	if data.Title != nil {
		activity.Title = *data.Title
	}
	if data.Description != nil {
		activity.Description = *data.Description
	}
	if data.ScheduledTime != nil {
		activity.ScheduledTime = *data.ScheduledTime
	}
	if data.Status != nil {
		activity.Status = entity.ActivityStatus(*data.Status)
	}
	if data.Priority != nil {
		activity.Priority = *data.Priority
	}

	activity.UpdatedAt = time.Now()

	return uc.activityRepo.Update(ctx, activity)
}

func (uc *ActivityUseCase) DeleteActivity(ctx context.Context, activityID uuid.UUID) error {
	return uc.activityRepo.Delete(ctx, activityID)
}

func (uc *ActivityUseCase) CompleteActivity(ctx context.Context, activityID uuid.UUID) error {
	activity, err := uc.activityRepo.GetByID(ctx, activityID)
	if err != nil {
		return fmt.Errorf("failed to get activity: %w", err)
	}
	if activity == nil {
		return fmt.Errorf("activity not found")
	}

	activity.Complete()
	return uc.activityRepo.Update(ctx, activity)
}

