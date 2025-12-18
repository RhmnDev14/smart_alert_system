package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/infrastructure/database"
)

type healthRepository struct {
	db *database.PostgresDB
}

func NewHealthRepository(db *database.PostgresDB) *healthRepository {
	return &healthRepository{db: db}
}

func (r *healthRepository) GetHealthProfileByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserHealthProfile, error) {
	query := `SELECT id, user_id, age, gender, medical_conditions, allergies, medications,
	          activity_preferences, health_goals, created_at, updated_at
	          FROM user_health_profiles WHERE user_id = $1`
	
	profile := &entity.UserHealthProfile{}
	var age sql.NullInt64
	
	err := r.db.DB.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID, &profile.UserID, &age, &profile.Gender, &profile.MedicalConditions,
		&profile.Allergies, &profile.Medications, &profile.ActivityPreferences,
		&profile.HealthGoals, &profile.CreatedAt, &profile.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if age.Valid {
		ageInt := int(age.Int64)
		profile.Age = &ageInt
	}
	
	return profile, nil
}

func (r *healthRepository) CreateOrUpdateHealthProfile(ctx context.Context, profile *entity.UserHealthProfile) error {
	query := `INSERT INTO user_health_profiles (id, user_id, age, gender, medical_conditions, allergies,
	          medications, activity_preferences, health_goals, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	          ON CONFLICT (user_id) DO UPDATE SET
	          age = EXCLUDED.age, gender = EXCLUDED.gender, medical_conditions = EXCLUDED.medical_conditions,
	          allergies = EXCLUDED.allergies, medications = EXCLUDED.medications,
	          activity_preferences = EXCLUDED.activity_preferences, health_goals = EXCLUDED.health_goals,
	          updated_at = EXCLUDED.updated_at`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		profile.ID, profile.UserID, profile.Age, profile.Gender, profile.MedicalConditions,
		profile.Allergies, profile.Medications, profile.ActivityPreferences, profile.HealthGoals,
		profile.CreatedAt, profile.UpdatedAt)
	return err
}

func (r *healthRepository) GetRecommendationTypes(ctx context.Context) ([]*entity.RecommendationType, error) {
	query := `SELECT id, name, description, trigger_condition FROM recommendation_types ORDER BY name`
	
	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var types []*entity.RecommendationType
	for rows.Next() {
		recType := &entity.RecommendationType{}
		err := rows.Scan(&recType.ID, &recType.Name, &recType.Description, &recType.TriggerCondition)
		if err != nil {
			return nil, err
		}
		types = append(types, recType)
	}
	return types, rows.Err()
}

func (r *healthRepository) CreateRecommendation(ctx context.Context, recommendation *entity.HealthRecommendation) error {
	query := `INSERT INTO health_recommendations (id, user_id, recommendation_type_id, recommendation_text,
	          activity_id, generated_at, sent_at, is_read, priority)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		recommendation.ID, recommendation.UserID, recommendation.RecommendationTypeID, recommendation.RecommendationText,
		recommendation.ActivityID, recommendation.GeneratedAt, recommendation.SentAt,
		recommendation.IsRead, recommendation.Priority)
	return err
}

func (r *healthRepository) GetRecommendationsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.HealthRecommendation, error) {
	query := `SELECT id, user_id, recommendation_type_id, recommendation_text, activity_id,
	          generated_at, sent_at, is_read, priority
	          FROM health_recommendations WHERE user_id = $1 ORDER BY generated_at DESC`
	
	rows, err := r.db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var recommendations []*entity.HealthRecommendation
	for rows.Next() {
		rec := &entity.HealthRecommendation{}
		var recTypeID, activityID sql.NullString
		var sentAt sql.NullTime
		
		err := rows.Scan(
			&rec.ID, &rec.UserID, &recTypeID, &rec.RecommendationText, &activityID,
			&rec.GeneratedAt, &sentAt, &rec.IsRead, &rec.Priority)
		if err != nil {
			return nil, err
		}
		
		if recTypeID.Valid {
			id, _ := uuid.Parse(recTypeID.String)
			rec.RecommendationTypeID = &id
		}
		if activityID.Valid {
			id, _ := uuid.Parse(activityID.String)
			rec.ActivityID = &id
		}
		if sentAt.Valid {
			rec.SentAt = &sentAt.Time
		}
		
		recommendations = append(recommendations, rec)
	}
	return recommendations, rows.Err()
}

