package repository

import (
	"context"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
)

type MessageRepository interface {
	Create(ctx context.Context, message *entity.MessageHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.MessageHistory, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.MessageHistory, error)
	Update(ctx context.Context, message *entity.MessageHistory) error
}

