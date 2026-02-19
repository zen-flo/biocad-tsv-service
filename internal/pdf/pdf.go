package pdf

import (
	"biocad-tsv-service/internal/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"

	"github.com/phpdave11/gofpdf"
)

// GenerateUnitPDF creates a PDF file with unitGUID data
func GenerateUnitPDF(ctx context.Context, outDir string, unitGUID uuid.UUID, msgRepo *repository.MessageRepo) error {
	messages, err := msgRepo.GetByUnitGUID(ctx, unitGUID)
	if err != nil {
		return fmt.Errorf("failed to get messages for unit %s: %w", unitGUID, err)
	}
	if len(messages) == 0 {
		return fmt.Errorf("no messages found for unit %s", unitGUID)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(fmt.Sprintf("Unit %s Report", unitGUID), false)
	pdf.AddPage()

	// title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("Unit Report: %s", unitGUID))
	pdf.Ln(12)

	// setting up the table
	pdf.SetFont("Arial", "B", 10)
	header := []string{
		"MsgId",
		"Text",
		"Class",
		"Level",
		"Area",
		"Addr",
		"Block",
		"Type",
		"Bit",
		"InvertBit",
		"CreatedAt",
	}
	colWidths := []float64{20, 40, 15, 10, 20, 20, 15, 15, 15, 20, 25}

	for i, h := range header {
		pdf.CellFormat(colWidths[i], 7, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// table contents
	pdf.SetFont("Arial", "", 10)
	for _, m := range messages {
		values := []string{
			m.MsgId,
			m.Text,
			m.Class,
			fmt.Sprintf("%d", m.Level),
			m.Area,
			m.Addr,
			nilOrString(m.Block),
			m.Type,
			nilOrString(m.Bit),
			nilOrString(m.InvertBit),
			m.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for i, v := range values {
			pdf.CellFormat(colWidths[i], 6, v, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// check the directory
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("failed to create output dir: %w", err)
		}
	}

	// saving PDF
	filePath := filepath.Join(outDir, fmt.Sprintf("%s.pdf", unitGUID))
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return fmt.Errorf("failed to save PDF: %w", err)
	}

	return nil
}

// nilOrString returns the string value or "-"
func nilOrString(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}
