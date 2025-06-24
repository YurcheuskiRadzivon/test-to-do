package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

func (txM *TxManager) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := txM.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
