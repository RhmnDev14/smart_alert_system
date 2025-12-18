package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/infrastructure/database"
)

type userRepository struct {
	db *database.PostgresDB
}

func NewUserRepository(db *database.PostgresDB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, whatsapp_number, name, timezone, is_active, is_first_time, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		user.ID, user.WhatsAppNumber, user.Name, user.Timezone,
		user.IsActive, user.IsFirstTime, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `SELECT id, whatsapp_number, name, timezone, is_active, is_first_time, 
	          created_at, updated_at, last_interaction_at 
	          FROM users WHERE id = $1`
	
	user := &entity.User{}
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.WhatsAppNumber, &user.Name, &user.Timezone,
		&user.IsActive, &user.IsFirstTime, &user.CreatedAt, &user.UpdatedAt, &user.LastInteractionAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) GetByWhatsAppNumber(ctx context.Context, whatsappNumber string) (*entity.User, error) {
	query := `SELECT id, whatsapp_number, name, timezone, is_active, is_first_time, 
	          created_at, updated_at, last_interaction_at 
	          FROM users WHERE whatsapp_number = $1`
	
	user := &entity.User{}
	err := r.db.DB.QueryRowContext(ctx, query, whatsappNumber).Scan(
		&user.ID, &user.WhatsAppNumber, &user.Name, &user.Timezone,
		&user.IsActive, &user.IsFirstTime, &user.CreatedAt, &user.UpdatedAt, &user.LastInteractionAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET name = $1, timezone = $2, is_active = $3, 
	          is_first_time = $4, updated_at = $5, last_interaction_at = $6 
	          WHERE id = $7`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		user.Name, user.Timezone, user.IsActive, user.IsFirstTime,
		user.UpdatedAt, user.LastInteractionAt, user.ID)
	return err
}

func (r *userRepository) UpdateLastInteraction(ctx context.Context, userID uuid.UUID, timestamp time.Time) error {
	query := `UPDATE users SET last_interaction_at = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.DB.ExecContext(ctx, query, timestamp, time.Now(), userID)
	return err
}

func (r *userRepository) GetAllActive(ctx context.Context) ([]*entity.User, error) {
	query := `SELECT id, whatsapp_number, name, timezone, is_active, is_first_time, 
	          created_at, updated_at, last_interaction_at 
	          FROM users WHERE is_active = true`
	
	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID, &user.WhatsAppNumber, &user.Name, &user.Timezone,
			&user.IsActive, &user.IsFirstTime, &user.CreatedAt, &user.UpdatedAt, &user.LastInteractionAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *userRepository) MarkAsNotFirstTime(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET is_first_time = false, updated_at = $1 WHERE id = $2`
	_, err := r.db.DB.ExecContext(ctx, query, time.Now(), userID)
	return err
}

