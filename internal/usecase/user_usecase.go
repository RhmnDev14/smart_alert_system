package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/domain/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) GetOrCreateUser(ctx context.Context, whatsappNumber, name, timezone string) (*entity.User, error) {
	user, err := uc.userRepo.GetByWhatsAppNumber(ctx, whatsappNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		// Create new user
		if timezone == "" {
			timezone = "Asia/Jakarta"
		}
		user = entity.NewUser(whatsappNumber, name, timezone)
		if err := uc.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Update last interaction
		now := time.Now()
		if err := uc.userRepo.UpdateLastInteraction(ctx, user.ID, now); err != nil {
			return nil, fmt.Errorf("failed to update last interaction: %w", err)
		}
		user.LastInteractionAt = &now
	}

	return user, nil
}

func (uc *UserUseCase) MarkAsNotFirstTime(ctx context.Context, userID uuid.UUID) error {
	return uc.userRepo.MarkAsNotFirstTime(ctx, userID)
}

func (uc *UserUseCase) GetAllActiveUsers(ctx context.Context) ([]*entity.User, error) {
	return uc.userRepo.GetAllActive(ctx)
}

