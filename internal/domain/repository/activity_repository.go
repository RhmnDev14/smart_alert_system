package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
)

type ActivityRepository interface {
	Create(ctx context.Context, activity *entity.Activity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Activity, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Activity, error)
	GetByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*entity.Activity, error)
	GetByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status entity.ActivityStatus) ([]*entity.Activity, error)
	Update(ctx context.Context, activity *entity.Activity) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTodayActivities(ctx context.Context, userID uuid.UUID) ([]*entity.Activity, error)
	GetCompletedToday(ctx context.Context, userID uuid.UUID) ([]*entity.Activity, error)
}

