package repository

import (
	"context"

	"connectHome/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HomeRepository struct {
	pool *pgxpool.Pool
}

func NewHomeRepository(pool *pgxpool.Pool) *HomeRepository {
	return &HomeRepository{pool: pool}
}

func (s *HomeRepository) Create(ctx context.Context, name, address string) (model.Home, error) {
	var home model.Home
	err := s.pool.QueryRow(ctx,
		`INSERT INTO homes (name, address) VALUES ($1, $2) RETURNING id, name, address`,
		name, address,
	).Scan(&home.ID, &home.Name, &home.Address)
	return home, err
}

func (s *HomeRepository) List(ctx context.Context) ([]model.Home, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name, address FROM homes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	homes := make([]model.Home, 0)
	for rows.Next() {
		var home model.Home
		if err := rows.Scan(&home.ID, &home.Name, &home.Address); err != nil {
			return nil, err
		}
		homes = append(homes, home)
	}
	return homes, rows.Err()
}
