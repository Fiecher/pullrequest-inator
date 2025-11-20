package pg

import (
	"context"
	"pullrequest-manager/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatusPGRepository struct {
	db *pgxpool.Pool
}

func NewStatusPGRepository(db *pgxpool.Pool) *StatusPGRepository {
	return &StatusPGRepository{db: db}
}

func (r *StatusPGRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Status, error) {
	var s models.Status
	if err := r.db.QueryRow(ctx,
		`SELECT id, name FROM pull_request_statuses WHERE id = $1`, id,
	).Scan(&s.ID, &s.Name); err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *StatusPGRepository) List(ctx context.Context) ([]*models.Status, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name FROM pull_request_statuses ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := make([]*models.Status, 0, 8)
	for rows.Next() {
		var s models.Status
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		statuses = append(statuses, &s)
	}
	return statuses, rows.Err()
}
