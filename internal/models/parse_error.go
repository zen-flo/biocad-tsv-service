package models

import (
	"github.com/google/uuid"
	"time"
)

// ParseError is an error when parsing a file
type ParseError struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Filename  string    `db:"filename" json:"filename"`
	RawLine   string    `db:"raw_line" json:"raw_line"`
	ErrorText string    `db:"error_text" json:"error_text"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
