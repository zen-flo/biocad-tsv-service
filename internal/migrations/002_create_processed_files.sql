-- Migration: create processed_files table
-- Stores metadata about files that have been processed

CREATE TABLE "processed_files" (
                                   id uuid PRIMARY KEY DEFAULT gen_random_uuid(),          -- unique identifier
                                   filename text UNIQUE NOT NULL,                          -- file name
                                   processed_at timestamp NOT NULL DEFAULT now(),          -- time when file was processed
                                   status text                                             -- processing status: success / failed
);

CREATE INDEX idx_processed_files_processed_at ON "processed_files"(processed_at);
