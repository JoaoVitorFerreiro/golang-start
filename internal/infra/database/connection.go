package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(databaseURL string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, err
    }
    
    // Configurações de performance
    config.MaxConns = 30
    config.MinConns = 5
    
    return pgxpool.NewWithConfig(context.Background(), config)
}