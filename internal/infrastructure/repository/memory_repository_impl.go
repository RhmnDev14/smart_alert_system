package repository

import (
	"context"

	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/infrastructure/database"

	"github.com/google/uuid"
)

type memoryRepository struct {
	db *database.PostgresDB
}

func NewMemoryRepository(db *database.PostgresDB) *memoryRepository {
	return &memoryRepository{db: db}
}

func (r *memoryRepository) Create(ctx context.Context, memory *entity.UserMemory) error {
	query := `INSERT INTO user_memories (id, user_id, memory_type, content, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Ext(ctx).ExecContext(ctx, query,
		memory.ID, memory.UserID, memory.MemoryType, memory.Content,
		memory.CreatedAt, memory.UpdatedAt)
	return err
}

func (r *memoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.UserMemory, error) {
	query := `SELECT id, user_id, memory_type, content, created_at, updated_at
	          FROM user_memories WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Ext(ctx).QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []*entity.UserMemory
	for rows.Next() {
		memory := &entity.UserMemory{}
		err := rows.Scan(
			&memory.ID, &memory.UserID, &memory.MemoryType, &memory.Content,
			&memory.CreatedAt, &memory.UpdatedAt)
		if err != nil {
			return nil, err
		}
		memories = append(memories, memory)
	}
	return memories, rows.Err()
}

func (r *memoryRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_memories WHERE id = $1`
	_, err := r.db.Ext(ctx).ExecContext(ctx, query, id)
	return err
}

func (r *memoryRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_memories WHERE user_id = $1`
	_, err := r.db.Ext(ctx).ExecContext(ctx, query, userID)
	return err
}
