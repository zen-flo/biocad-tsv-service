package queue

import (
	"biocad-tsv-service/internal/repository"
	"context"
	"log"
	"path/filepath"
	"time"
)

// Scanner periodically scans a directory for new TSV files
type Scanner struct {
	InputDir string
	PFRepo   *repository.ProcessedFileRepo
	Queue    chan<- string
	QM       *Manager
	Interval time.Duration
	done     chan struct{}
}

// NewScanner creates a new Scanner
func NewScanner(inputDir string, pfRepo *repository.ProcessedFileRepo, queue chan<- string, qm *Manager, interval time.Duration) *Scanner {
	return &Scanner{
		InputDir: inputDir,
		PFRepo:   pfRepo,
		Queue:    queue,
		QM:       qm,
		Interval: interval,
		done:     make(chan struct{}),
	}
}

// Start launches the scanner goroutine
func (s *Scanner) Start(ctx context.Context) {
	go func() {
		log.Println("[scanner] started")
		ticker := time.NewTicker(s.Interval)
		defer func() {
			ticker.Stop()
			log.Println("[scanner] stopped")
		}()

		// initial scan
		s.scan(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-s.done:
				return
			case <-ticker.C:
				s.scan(ctx)
			}
		}
	}()
}

// Stop stops the scanner
func (s *Scanner) Stop() {
	close(s.done)
}

// scan performs a single scan of the input directory
func (s *Scanner) scan(ctx context.Context) {
	files, err := filepath.Glob(filepath.Join(s.InputDir, "*.tsv"))
	if err != nil {
		log.Printf("[scanner] failed to list TSV files in %s: %v", s.InputDir, err)
		return
	}

	for _, file := range files {
		select {
		case <-ctx.Done():
			return
		default:
		}

		processed, err := s.PFRepo.IsProcessed(ctx, file)
		if err != nil || processed {
			continue
		}

		if s.QM.Add(file) {
			log.Printf("[scanner] queueing new file %s", file)
			s.Queue <- file
		}
	}
}
