package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/infrastructure/database"
)

type categoryRepository struct {
	db *database.PostgresDB
}

func NewCategoryRepository(db *database.PostgresDB) *categoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]*entity.ActivityCategory, error) {
	query := `SELECT id, name, description, icon, color FROM activity_categories ORDER BY name`
	
	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var categories []*entity.ActivityCategory
	for rows.Next() {
		category := &entity.ActivityCategory{}
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Icon, &category.Color)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ActivityCategory, error) {
	query := `SELECT id, name, description, icon, color FROM activity_categories WHERE id = $1`
	
	category := &entity.ActivityCategory{}
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&category.ID, &category.Name, &category.Description, &category.Icon, &category.Color)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return category, err
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (*entity.ActivityCategory, error) {
	query := `SELECT id, name, description, icon, color FROM activity_categories WHERE name = $1`
	
	category := &entity.ActivityCategory{}
	err := r.db.DB.QueryRowContext(ctx, query, name).Scan(
		&category.ID, &category.Name, &category.Description, &category.Icon, &category.Color)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return category, err
}

