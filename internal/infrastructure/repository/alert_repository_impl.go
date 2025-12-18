package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/infrastructure/database"
)

type alertRepository struct {
	db *database.PostgresDB
}

func NewAlertRepository(db *database.PostgresDB) *alertRepository {
	return &alertRepository{db: db}
}

func (r *alertRepository) Create(ctx context.Context, alert *entity.AlertLog) error {
	query := `INSERT INTO alert_logs (id, user_id, alert_type, alert_content, scheduled_time,
	          sent_at, is_sent, status, error_message, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		alert.ID, alert.UserID, alert.AlertType, alert.AlertContent, alert.ScheduledTime,
		alert.SentAt, alert.IsSent, alert.Status, alert.ErrorMessage, alert.CreatedAt)
	return err
}

func (r *alertRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AlertLog, error) {
	query := `SELECT id, user_id, alert_type, alert_content, scheduled_time, sent_at,
	          is_sent, status, error_message, created_at
	          FROM alert_logs WHERE id = $1`
	
	alert := &entity.AlertLog{}
	var sentAt sql.NullTime
	
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&alert.ID, &alert.UserID, &alert.AlertType, &alert.AlertContent, &alert.ScheduledTime,
		&sentAt, &alert.IsSent, &alert.Status, &alert.ErrorMessage, &alert.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if sentAt.Valid {
		alert.SentAt = &sentAt.Time
	}
	
	return alert, nil
}

func (r *alertRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.AlertLog, error) {
	query := `SELECT id, user_id, alert_type, alert_content, scheduled_time, sent_at,
	          is_sent, status, error_message, created_at
	          FROM alert_logs WHERE user_id = $1 ORDER BY scheduled_time DESC`
	
	rows, err := r.db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var alerts []*entity.AlertLog
	for rows.Next() {
		alert := &entity.AlertLog{}
		var sentAt sql.NullTime
		
		err := rows.Scan(
			&alert.ID, &alert.UserID, &alert.AlertType, &alert.AlertContent, &alert.ScheduledTime,
			&sentAt, &alert.IsSent, &alert.Status, &alert.ErrorMessage, &alert.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		if sentAt.Valid {
			alert.SentAt = &sentAt.Time
		}
		
		alerts = append(alerts, alert)
	}
	return alerts, rows.Err()
}

func (r *alertRepository) Update(ctx context.Context, alert *entity.AlertLog) error {
	query := `UPDATE alert_logs SET alert_content = $1, sent_at = $2, is_sent = $3, status = $4, error_message = $5
	          WHERE id = $6`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		alert.AlertContent, alert.SentAt, alert.IsSent, alert.Status, alert.ErrorMessage, alert.ID)
	return err
}

func (r *alertRepository) GetPendingAlerts(ctx context.Context, alertType entity.AlertType) ([]*entity.AlertLog, error) {
	query := `SELECT id, user_id, alert_type, alert_content, scheduled_time, sent_at,
	          is_sent, status, error_message, created_at
	          FROM alert_logs WHERE alert_type = $1 AND is_sent = false AND status = $2
	          ORDER BY scheduled_time ASC`
	
	rows, err := r.db.DB.QueryContext(ctx, query, alertType, entity.AlertStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var alerts []*entity.AlertLog
	for rows.Next() {
		alert := &entity.AlertLog{}
		var sentAt sql.NullTime
		
		err := rows.Scan(
			&alert.ID, &alert.UserID, &alert.AlertType, &alert.AlertContent, &alert.ScheduledTime,
			&sentAt, &alert.IsSent, &alert.Status, &alert.ErrorMessage, &alert.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		if sentAt.Valid {
			alert.SentAt = &sentAt.Time
		}
		
		alerts = append(alerts, alert)
	}
	return alerts, rows.Err()
}

func (r *alertRepository) GetByScheduledTime(ctx context.Context, startTime, endTime time.Time) ([]*entity.AlertLog, error) {
	query := `SELECT id, user_id, alert_type, alert_content, scheduled_time, sent_at,
	          is_sent, status, error_message, created_at
	          FROM alert_logs WHERE scheduled_time >= $1 AND scheduled_time <= $2
	          ORDER BY scheduled_time ASC`
	
	rows, err := r.db.DB.QueryContext(ctx, query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var alerts []*entity.AlertLog
	for rows.Next() {
		alert := &entity.AlertLog{}
		var sentAt sql.NullTime
		
		err := rows.Scan(
			&alert.ID, &alert.UserID, &alert.AlertType, &alert.AlertContent, &alert.ScheduledTime,
			&sentAt, &alert.IsSent, &alert.Status, &alert.ErrorMessage, &alert.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		if sentAt.Valid {
			alert.SentAt = &sentAt.Time
		}
		
		alerts = append(alerts, alert)
	}
	return alerts, rows.Err()
}

