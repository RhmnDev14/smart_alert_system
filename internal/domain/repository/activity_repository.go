package repository

import (
	"context"
	"time"

	"smart_alert_system/internal/domain/entity"

	"github.com/google/uuid"
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
	GetPendingActivitiesToRemind(ctx context.Context, until time.Time) ([]*entity.Activity, error)
	// SearchByKeyword searches activities by keyword in title or description, scoped to a user
	SearchByKeyword(ctx context.Context, userID uuid.UUID, keyword string) ([]*entity.Activity, error)
	// GetByUserIDAndTitleFuzzy finds activities matching a title (case-insensitive LIKE), scoped to a user
	GetByUserIDAndTitleFuzzy(ctx context.Context, userID uuid.UUID, title string) ([]*entity.Activity, error)
}
