package main

import (
	"biocad-tsv-service/internal/api"
	"biocad-tsv-service/internal/config"
	"biocad-tsv-service/internal/database"
	"biocad-tsv-service/internal/parser"
	"biocad-tsv-service/internal/pdf"
	"biocad-tsv-service/internal/queue"
	"biocad-tsv-service/internal/repository"
	"biocad-tsv-service/internal/util"
	"context"
	"github.com/google/uuid"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const numWorkers = 4

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("[main] failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("[main] failed to validate config: %v", err)
	}

	dbPool, err := database.NewPool(cfg)
	if err != nil {
		log.Fatalf("[main] failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	util.EnsureDirs(cfg.Dirs.Input, cfg.Dirs.Output)

	log.Printf("[main] Loaded config: %s", cfg)
	log.Println("Service started successfully")

	// create repositories
	msgRepo := repository.NewMessageRepo(dbPool)
	pfRepo := repository.NewProcessedFileRepo(dbPool)
	errRepo := repository.NewParseErrorRepo(dbPool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start API server
	apiServer := api.NewServer(msgRepo)
	apiServer.Start(ctx, cfg.Server.Port)

	// channel for files queue
	fileQueue := make(chan string, 100)
	var wg sync.WaitGroup

	queueManager := queue.New()

	// start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, fileQueue, msgRepo, pfRepo, errRepo, queueManager, &wg, cfg.Dirs.Output)
	}

	// start scanner
	scanner := queue.NewScanner(cfg.Dirs.Input, pfRepo, fileQueue, queueManager, 30*time.Second)
	scanner.Start(ctx)

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("[main] Shutdown signal received, stopping scanner and workers...")

	cancel()         // cancel context for any ongoing operations
	close(fileQueue) // signal workers to finish
	wg.Wait()        // wait for all workers
	log.Println("[main] Service stopped gracefully")
}

// worker processes files from the queue
func worker(
	ctx context.Context,
	id int,
	queue <-chan string,
	msgRepo *repository.MessageRepo,
	pfRepo *repository.ProcessedFileRepo,
	errRepo *repository.ParseErrorRepo,
	qm *queue.Manager,
	wg *sync.WaitGroup,
	outDir string,
) {
	defer wg.Done()
	for file := range queue {
		select {
		case <-ctx.Done():
			log.Printf("[worker %d] context canceled, exiting", id)
			return
		default:
		}

		log.Printf("[worker %d] processing file: %s", id, file)
		messages, err := parser.ParseTSVFile(ctx, file, msgRepo, pfRepo, errRepo)
		if err != nil {
			log.Printf("[worker %d] failed to parse file %s: %v", id, file, err)
			qm.Remove(file)
			continue
		} else {
			log.Printf("[worker %d] successfully parsed file %s", id, file)
		}

		// generating a PDF for each unique unitGUID
		unitGUIDMap := make(map[uuid.UUID]struct{})
		for _, msg := range messages {
			unitGUIDMap[msg.UnitGUID] = struct{}{}
		}

		for unitGUID := range unitGUIDMap {
			if err := pdf.GenerateUnitPDF(ctx, outDir, unitGUID, msgRepo); err != nil {
				log.Printf("[worker %d] failed to generate PDF for %s: %v", id, unitGUID, err)
			} else {
				log.Printf("[worker %d] PDF generated for %s", id, unitGUID)
			}
		}

		qm.Remove(file)
	}
}
