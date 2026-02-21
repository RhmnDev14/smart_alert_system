package repository

import (
	"context"
	"smart_alert_system/internal/domain/entity"
)

type WorkOrderRepository interface {
	// Create creates a new work order
	Create(ctx context.Context, workOrder *entity.WorkOrder) error

	// GetByID retrieves work order by ID
	GetByID(ctx context.Context, idOrder string) (*entity.WorkOrder, error)

	// GetByNoTelp retrieves work orders by phone number
	GetByNoTelp(ctx context.Context, noTelp string) ([]*entity.WorkOrder, error)

	// GetByTglKunjungan retrieves work orders scheduled for a specific date
	GetByTglKunjungan(ctx context.Context, date string) ([]*entity.WorkOrder, error)

	// GetPending retrieves pending work orders
	GetPending(ctx context.Context) ([]*entity.WorkOrder, error)

	// GetPendingByNoTelp retrieves pending work orders for a phone number
	GetPendingByNoTelp(ctx context.Context, noTelp string) ([]*entity.WorkOrder, error)

	// Update updates work order
	Update(ctx context.Context, workOrder *entity.WorkOrder) error

	// UpdateSendWhatsapp updates send_whatsapp flag
	UpdateSendWhatsapp(ctx context.Context, idOrder string, sent bool) error

	// Delete deletes work order
	Delete(ctx context.Context, idOrder string) error

	// BulkCreate creates multiple work orders
	BulkCreate(ctx context.Context, workOrders []*entity.WorkOrder) error
}

