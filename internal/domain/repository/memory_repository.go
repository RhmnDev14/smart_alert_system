package repository

import (
	"context"

	"smart_alert_system/internal/domain/entity"

	"github.com/google/uuid"
)

type MemoryRepository interface {
	Create(ctx context.Context, memory *entity.UserMemory) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserMemory, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
