package repository

import (
	"biocad-tsv-service/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ParseErrorRepo struct {
	db *pgxpool.Pool
}

func NewParseErrorRepo(db *pgxpool.Pool) *ParseErrorRepo {
	return &ParseErrorRepo{db: db}
}

func (r *ParseErrorRepo) Insert(ctx context.Context, e *models.ParseError) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO "parse_errors" (id, filename, raw_line, error_text, created_at)
		VALUES ($1,$2,$3,$4,$5)
	`,
		e.ID, e.Filename, e.RawLine, e.ErrorText, e.CreatedAt,
	)
	return err
}

// List returns parse errors with pagination
func (r *ParseErrorRepo) List(ctx context.Context, limit, offset int) ([]models.ParseError, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, filename, raw_line, error_text, created_at
		FROM "parse_errors"
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list parse_errors failed: %w", err)
	}
	defer rows.Close()

	var errs []models.ParseError
	for rows.Next() {
		var e models.ParseError
		if err := rows.Scan(&e.ID, &e.Filename, &e.RawLine, &e.ErrorText, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan parse_error failed: %w", err)
		}
		errs = append(errs, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return errs, nil
}
