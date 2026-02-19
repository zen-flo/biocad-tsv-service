-- Migration: create parse_errors table
-- Stores errors that occurred during parsing of TSV files

CREATE TABLE "parse_errors" (
                                id uuid PRIMARY KEY DEFAULT gen_random_uuid(),         -- unique identifier
                                filename text,                                         -- file name that caused error
                                raw_line text,                                         -- raw line from the file
                                error_text text,                                       -- description of parsing error
                                created_at timestamp NOT NULL DEFAULT now()            -- timestamp of error creation
);

CREATE INDEX idx_parse_errors_created_at ON "parse_errors"(created_at);
