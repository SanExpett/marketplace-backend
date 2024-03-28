package repository

import (
	"context"
	"fmt"
	myerrors "github.com/SanExpett/marketplace-backend/pkg/my_errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(ctx context.Context, urlDataBase string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, urlDataBase)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return pool, nil
}
