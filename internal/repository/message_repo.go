package repository

import (
	"biocad-tsv-service/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type MessageRepo struct {
	db *pgxpool.Pool
}

func NewMessageRepo(db *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) Insert(ctx context.Context, msg *models.Message) error {
	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO "messages" 
		    (id, mqtt, unit_guid, msg_id, text, context, class, level, area, addr, block, 
		     type, bit, invert_bit, created_at) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
	`,
		msg.ID, msg.MQTT, msg.UnitGUID, msg.MsgId, msg.Text, msg.Context, msg.Class,
		msg.Level, msg.Area, msg.Addr, msg.Block, msg.Type, msg.Bit, msg.InvertBit, msg.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert message failed: %w", err)
	}
	return nil
}

// GetByUnitGUID returns all messages for a given device
func (r *MessageRepo) GetByUnitGUID(ctx context.Context, unitGUID uuid.UUID) ([]models.Message, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, mqtt, unit_guid, msg_id, text, context, class, level, area, addr, block, type, bit, invert_bit, created_at
		FROM "messages"
		WHERE unit_guid=$1
		ORDER BY created_at DESC
	`, unitGUID)
	if err != nil {
		return nil, fmt.Errorf("query messages failed: %w", err)
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID, &m.MQTT, &m.UnitGUID, &m.MsgId, &m.Text, &m.Context, &m.Class, &m.Level,
			&m.Area, &m.Addr, &m.Block, &m.Type, &m.Bit, &m.InvertBit, &m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan message failed: %w", err)
		}
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return msgs, nil
}
