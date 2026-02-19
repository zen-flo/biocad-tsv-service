package models

import (
	"github.com/google/uuid"
	"time"
)

// Message is one line from TSV
type Message struct {
	ID        uuid.UUID `db:"id" json:"id"`
	MQTT      string    `db:"mqtt" json:"mqtt"`           // optional MQTT broker or topic
	UnitGUID  uuid.UUID `db:"unit_guid" json:"unit_guid"` // device id
	MsgId     string    `db:"msg_id" json:"msg_id"`
	Text      string    `db:"text" json:"text"`
	Context   string    `db:"context" json:"context"`
	Class     string    `db:"class" json:"class"`
	Level     int       `db:"level" json:"level"`
	Area      string    `db:"area" json:"area"`
	Addr      string    `db:"addr" json:"addr"`
	Block     *string   `db:"block" json:"block"`
	Type      string    `db:"type" json:"type"`
	Bit       *string   `db:"bit" json:"bit"`
	InvertBit *string   `db:"invert_bit" json:"invert_bit"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
