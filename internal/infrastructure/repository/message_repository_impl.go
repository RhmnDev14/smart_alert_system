package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/infrastructure/database"
)

type messageRepository struct {
	db *database.PostgresDB
}

func NewMessageRepository(db *database.PostgresDB) *messageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *entity.MessageHistory) error {
	query := `INSERT INTO message_history (id, user_id, message_content, message_type, intent_detected,
	          ai_response, received_at, sent_at, is_processed, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		message.ID, message.UserID, message.MessageContent, message.MessageType,
		message.IntentDetected, message.AIResponse, message.ReceivedAt, message.SentAt,
		message.IsProcessed, message.CreatedAt)
	return err
}

func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MessageHistory, error) {
	query := `SELECT id, user_id, message_content, message_type, intent_detected, ai_response,
	          received_at, sent_at, is_processed, created_at
	          FROM message_history WHERE id = $1`
	
	message := &entity.MessageHistory{}
	var receivedAt, sentAt sql.NullTime
	
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&message.ID, &message.UserID, &message.MessageContent, &message.MessageType,
		&message.IntentDetected, &message.AIResponse, &receivedAt, &sentAt,
		&message.IsProcessed, &message.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if receivedAt.Valid {
		message.ReceivedAt = &receivedAt.Time
	}
	if sentAt.Valid {
		message.SentAt = &sentAt.Time
	}
	
	return message, nil
}

func (r *messageRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.MessageHistory, error) {
	query := `SELECT id, user_id, message_content, message_type, intent_detected, ai_response,
	          received_at, sent_at, is_processed, created_at
	          FROM message_history WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`
	
	rows, err := r.db.DB.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []*entity.MessageHistory
	for rows.Next() {
		message := &entity.MessageHistory{}
		var receivedAt, sentAt sql.NullTime
		
		err := rows.Scan(
			&message.ID, &message.UserID, &message.MessageContent, &message.MessageType,
			&message.IntentDetected, &message.AIResponse, &receivedAt, &sentAt,
			&message.IsProcessed, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		if receivedAt.Valid {
			message.ReceivedAt = &receivedAt.Time
		}
		if sentAt.Valid {
			message.SentAt = &sentAt.Time
		}
		
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (r *messageRepository) Update(ctx context.Context, message *entity.MessageHistory) error {
	query := `UPDATE message_history SET message_content = $1, message_type = $2, intent_detected = $3,
	          ai_response = $4, received_at = $5, sent_at = $6, is_processed = $7
	          WHERE id = $8`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		message.MessageContent, message.MessageType, message.IntentDetected,
		message.AIResponse, message.ReceivedAt, message.SentAt, message.IsProcessed, message.ID)
	return err
}

