package repository

import "context"

type TransactionManager interface {
	RunInTx(ctx context.Context, fn func(txCtx context.Context) error) error
}
