package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
)

type AlertRepository interface {
	Create(ctx context.Context, alert *entity.AlertLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AlertLog, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.AlertLog, error)
	Update(ctx context.Context, alert *entity.AlertLog) error
	GetPendingAlerts(ctx context.Context, alertType entity.AlertType) ([]*entity.AlertLog, error)
	GetByScheduledTime(ctx context.Context, startTime, endTime time.Time) ([]*entity.AlertLog, error)
}

