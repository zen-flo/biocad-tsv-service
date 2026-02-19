package parser

import (
	"biocad-tsv-service/internal/models"
	"biocad-tsv-service/internal/repository"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// column indexes in TSV file
const (
	colMQTT      = 0
	colN         = 1
	colUnitGUID  = 2
	colMsgID     = 3
	colText      = 4
	colContext   = 5
	colClass     = 6
	colLevel     = 7
	colArea      = 8
	colAddr      = 9
	colBlock     = 10
	colType      = 11
	colBit       = 12
	colInvertBit = 13
	expectedCols = 14
)

// ParseTSVFile reads a TSV file and stores messages into the database
func ParseTSVFile(
	ctx context.Context,
	filePath string,
	msgRepo *repository.MessageRepo,
	pfRepo *repository.ProcessedFileRepo,
	errRepo *repository.ParseErrorRepo,
) ([]*models.Message, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file %s", filePath)
		}
	}()

	reader := csv.NewReader(f)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1 // allow variable number of columns

	var hadErrors bool
	var processedMessages []*models.Message

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return processedMessages, fmt.Errorf("failed to read record from TSV file %s: %w", filePath, err)
		}

		if len(record) < expectedCols {
			// if the line is too short, save the error
			hadErrors = true
			_ = errRepo.Insert(ctx, &models.ParseError{
				ID:        uuid.New(),
				Filename:  filePath,
				RawLine:   strings.Join(record, "\t"),
				ErrorText: fmt.Sprintf("not enough columns, expected %d got %d", expectedCols, len(record)),
				CreatedAt: time.Now(),
			})
			continue
		}

		// parse unit GUID
		unitGUID, err := uuid.Parse(record[colUnitGUID])
		if err != nil {
			hadErrors = true
			_ = errRepo.Insert(ctx, &models.ParseError{
				ID:        uuid.New(),
				Filename:  filePath,
				RawLine:   strings.Join(record, "\t"),
				ErrorText: fmt.Sprintf("invalid unit_guid value: %v", err),
				CreatedAt: time.Now(),
			})
			continue
		}

		// parse level safety
		level, err := strconv.Atoi(record[colLevel])
		if err != nil {
			hadErrors = true
			_ = errRepo.Insert(ctx, &models.ParseError{
				ID:        uuid.New(),
				Filename:  filePath,
				RawLine:   strings.Join(record, "\t"),
				ErrorText: fmt.Sprintf("invalid level value: %v", err),
				CreatedAt: time.Now(),
			})
			continue
		}

		msg := &models.Message{
			ID:        uuid.New(),
			MQTT:      record[colMQTT],
			UnitGUID:  unitGUID,
			MsgId:     record[colMsgID],
			Text:      record[colText],
			Context:   record[colContext],
			Class:     record[colClass],
			Level:     level,
			Area:      record[colArea],
			Addr:      record[colAddr],
			Block:     emptyToNil(record[colBlock]),
			Type:      record[colType],
			Bit:       emptyToNil(record[colBit]),
			InvertBit: emptyToNil(record[colInvertBit]),
			CreatedAt: time.Now(),
		}

		if err := msgRepo.Insert(ctx, msg); err != nil {
			hadErrors = true
			_ = errRepo.Insert(ctx, &models.ParseError{
				ID:        uuid.New(),
				Filename:  filePath,
				RawLine:   strings.Join(record, "\t"),
				ErrorText: fmt.Sprintf("database insert failed: %v", err),
				CreatedAt: time.Now(),
			})
		} else {
			processedMessages = append(processedMessages, msg)
		}
	}

	status := "success"
	if hadErrors {
		status = "failed"
	}

	if err := pfRepo.Insert(ctx, &models.ProcessedFile{
		ID:          uuid.New(),
		Filename:    filePath,
		ProcessedAt: time.Now(),
		Status:      status,
	}); err != nil {
		log.Printf("failed to mark file %s as processed: %v", filePath, err)
	}

	return processedMessages, nil
}

func emptyToNil(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}
