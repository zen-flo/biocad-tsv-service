package repository

import (
	"biocad-tsv-service/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ProcessedFileRepo struct {
	db *pgxpool.Pool
}

func NewProcessedFileRepo(db *pgxpool.Pool) *ProcessedFileRepo {
	return &ProcessedFileRepo{db: db}
}

func (r *ProcessedFileRepo) Insert(ctx context.Context, file *models.ProcessedFile) error {
	if file.ID == uuid.Nil {
		file.ID = uuid.New()
	}
	if file.ProcessedAt.IsZero() {
		file.ProcessedAt = time.Now()
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO "processed_files" (id, filename, processed_at, status)
		VALUES ($1,$2,$3,$4)
	`,
		file.ID, file.Filename, file.ProcessedAt, file.Status,
	)
	return err
}

// List returns processed files with pagination
func (r *ProcessedFileRepo) List(ctx context.Context, limit, offset int) ([]models.ProcessedFile, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, filename, processed_at, status
		FROM "processed_files"
		ORDER BY processed_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list processed_files failed: %w", err)
	}
	defer rows.Close()

	var files []models.ProcessedFile
	for rows.Next() {
		var f models.ProcessedFile
		if err := rows.Scan(&f.ID, &f.Filename, &f.ProcessedAt, &f.Status); err != nil {
			return nil, fmt.Errorf("scan processed_file failed: %w", err)
		}
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return files, nil
}

// IsProcessed checks if a file with the given filename has already been processed successfully
func (r *ProcessedFileRepo) IsProcessed(ctx context.Context, filename string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 
            FROM processed_files 
            WHERE filename=$1 AND status='success'
        )
    `, filename).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if file is processed: %w", err)
	}
	return exists, nil
}
