package dto

import "github.com/jackc/pgx/v5/pgtype"

type Lyric struct {
	UUID      pgtype.UUID        `json:"uuid"`
	AudioUUID pgtype.UUID        `json:"audio_uuid"`
	Order     int                `json:"order"`
	Text      string             `json:"text"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type LyricRead struct {
	UUID      pgtype.UUID        `json:"uuid"`
	AudioUUID pgtype.UUID        `json:"audio_uuid"`
	Order     int                `json:"order"`
	Text      string             `json:"text"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type LyricCreate struct {
	AudioUUID pgtype.UUID `json:"audio_uuid"`
	Order     int         `json:"order"`
	Text      string      `json:"text"`
}
