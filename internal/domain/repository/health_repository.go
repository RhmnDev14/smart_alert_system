package repository

import (
	"context"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
)

type HealthRepository interface {
	GetHealthProfileByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserHealthProfile, error)
	CreateOrUpdateHealthProfile(ctx context.Context, profile *entity.UserHealthProfile) error
	GetRecommendationTypes(ctx context.Context) ([]*entity.RecommendationType, error)
	CreateRecommendation(ctx context.Context, recommendation *entity.HealthRecommendation) error
	GetRecommendationsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.HealthRecommendation, error)
}

