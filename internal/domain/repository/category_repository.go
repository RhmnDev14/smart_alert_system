package repository

import (
	"context"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]*entity.ActivityCategory, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ActivityCategory, error)
	GetByName(ctx context.Context, name string) (*entity.ActivityCategory, error)
}

