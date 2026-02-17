package models

import (
	"github.com/google/uuid"
	"time"
)

// ProcessedFile is a file that has already been processed
type ProcessedFile struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Filename    string    `db:"filename" json:"filename"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at"`
	Status      string    `db:"status" json:"status"` // success / failed
}
