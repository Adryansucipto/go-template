package tx

import (
	"context"
	"go-template/pkg/resource/model"

	"gorm.io/gorm"
)

type txCtxKey struct{}

// InjectTx embeds a *gorm.DB transaction into the context.
func InjectTx(ctx context.Context, txDB *gorm.DB) context.Context {
	return context.WithValue(ctx, txCtxKey{}, txDB)
}

// ExtractTx retrieves an active *gorm.DB transaction from context if present,
// otherwise falls back to the default *gorm.DB connection.
func ExtractTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if txDB, ok := ctx.Value(txCtxKey{}).(*gorm.DB); ok && txDB != nil {
		return txDB.WithContext(ctx)
	}
	return defaultDB.WithContext(ctx)
}

type Manager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type gormTxManager struct {
	db *gorm.DB
}

func New(db model.Database) Manager {
	return &gormTxManager{
		db: db.PostgreSql,
	}
}

func (m *gormTxManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	if m.db == nil {
		return fn(ctx)
	}
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := InjectTx(ctx, tx)
		return fn(txCtx)
	})
}

