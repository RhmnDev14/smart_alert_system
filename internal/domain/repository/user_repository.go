package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByWhatsAppNumber(ctx context.Context, whatsappNumber string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	UpdateLastInteraction(ctx context.Context, userID uuid.UUID, timestamp time.Time) error
	GetAllActive(ctx context.Context) ([]*entity.User, error)
	MarkAsNotFirstTime(ctx context.Context, userID uuid.UUID) error
}

