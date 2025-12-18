package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"smart_alert_system/internal/domain/entity"
	"smart_alert_system/internal/infrastructure/database"
)

type activityRepository struct {
	db *database.PostgresDB
}

func NewActivityRepository(db *database.PostgresDB) *activityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(ctx context.Context, activity *entity.Activity) error {
	query := `INSERT INTO activities (id, user_id, category_id, title, description, scheduled_time, 
	          reminder_time, status, priority, created_at, updated_at, completed_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		activity.ID, activity.UserID, activity.CategoryID, activity.Title, activity.Description,
		activity.ScheduledTime, activity.ReminderTime, activity.Status, activity.Priority,
		activity.CreatedAt, activity.UpdatedAt, activity.CompletedAt)
	return err
}

func (r *activityRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Activity, error) {
	query := `SELECT id, user_id, category_id, title, description, scheduled_time, reminder_time,
	          status, priority, created_at, updated_at, completed_at
	          FROM activities WHERE id = $1`
	
	activity := &entity.Activity{}
	var categoryID sql.NullString
	var reminderTime, completedAt sql.NullTime
	
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&activity.ID, &activity.UserID, &categoryID, &activity.Title, &activity.Description,
		&activity.ScheduledTime, &reminderTime, &activity.Status, &activity.Priority,
		&activity.CreatedAt, &activity.UpdatedAt, &completedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	if categoryID.Valid {
		id, _ := uuid.Parse(categoryID.String)
		activity.CategoryID = &id
	}
	if reminderTime.Valid {
		activity.ReminderTime = &reminderTime.Time
	}
	if completedAt.Valid {
		activity.CompletedAt = &completedAt.Time
	}
	
	return activity, nil
}

func (r *activityRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Activity, error) {
	query := `SELECT id, user_id, category_id, title, description, scheduled_time, reminder_time,
	          status, priority, created_at, updated_at, completed_at
	          FROM activities WHERE user_id = $1 ORDER BY scheduled_time ASC`
	
	return r.scanActivities(ctx, query, userID)
}

func (r *activityRepository) GetByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*entity.Activity, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	query := `SELECT id, user_id, category_id, title, description, scheduled_time, reminder_time,
	          status, priority, created_at, updated_at, completed_at
	          FROM activities WHERE user_id = $1 AND scheduled_time >= $2 AND scheduled_time < $3
	          ORDER BY scheduled_time ASC`
	
	rows, err := r.db.DB.QueryContext(ctx, query, userID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanActivityRows(rows)
}

func (r *activityRepository) GetByUserIDAndStatus(ctx context.Context, userID uuid.UUID, status entity.ActivityStatus) ([]*entity.Activity, error) {
	query := `SELECT id, user_id, category_id, title, description, scheduled_time, reminder_time,
	          status, priority, created_at, updated_at, completed_at
	          FROM activities WHERE user_id = $1 AND status = $2 ORDER BY scheduled_time ASC`
	
	return r.scanActivities(ctx, query, userID, status)
}

func (r *activityRepository) Update(ctx context.Context, activity *entity.Activity) error {
	query := `UPDATE activities SET category_id = $1, title = $2, description = $3, scheduled_time = $4,
	          reminder_time = $5, status = $6, priority = $7, updated_at = $8, completed_at = $9
	          WHERE id = $10`
	
	_, err := r.db.DB.ExecContext(ctx, query,
		activity.CategoryID, activity.Title, activity.Description, activity.ScheduledTime,
		activity.ReminderTime, activity.Status, activity.Priority, activity.UpdatedAt,
		activity.CompletedAt, activity.ID)
	return err
}

func (r *activityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM activities WHERE id = $1`
	_, err := r.db.DB.ExecContext(ctx, query, id)
	return err
}

func (r *activityRepository) GetTodayActivities(ctx context.Context, userID uuid.UUID) ([]*entity.Activity, error) {
	now := time.Now()
	return r.GetByUserIDAndDate(ctx, userID, now)
}

func (r *activityRepository) GetCompletedToday(ctx context.Context, userID uuid.UUID) ([]*entity.Activity, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	query := `SELECT id, user_id, category_id, title, description, scheduled_time, reminder_time,
	          status, priority, created_at, updated_at, completed_at
	          FROM activities WHERE user_id = $1 AND status = $2 AND completed_at >= $3 AND completed_at < $4
	          ORDER BY completed_at ASC`
	
	rows, err := r.db.DB.QueryContext(ctx, query, userID, entity.ActivityStatusCompleted, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanActivityRows(rows)
}

func (r *activityRepository) scanActivities(ctx context.Context, query string, args ...interface{}) ([]*entity.Activity, error) {
	rows, err := r.db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanActivityRows(rows)
}

func (r *activityRepository) scanActivityRows(rows *sql.Rows) ([]*entity.Activity, error) {
	var activities []*entity.Activity
	for rows.Next() {
		activity := &entity.Activity{}
		var categoryID sql.NullString
		var reminderTime, completedAt sql.NullTime
		
		err := rows.Scan(
			&activity.ID, &activity.UserID, &categoryID, &activity.Title, &activity.Description,
			&activity.ScheduledTime, &reminderTime, &activity.Status, &activity.Priority,
			&activity.CreatedAt, &activity.UpdatedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		
		if categoryID.Valid {
			id, _ := uuid.Parse(categoryID.String)
			activity.CategoryID = &id
		}
		if reminderTime.Valid {
			activity.ReminderTime = &reminderTime.Time
		}
		if completedAt.Valid {
			activity.CompletedAt = &completedAt.Time
		}
		
		activities = append(activities, activity)
	}
	return activities, rows.Err()
}

