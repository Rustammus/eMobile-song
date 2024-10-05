package schema

import (
	"eMobile/internal/dto"
	"github.com/jackc/pgx/v5/pgtype"
)

type RequestLyricCreate struct {
	AudioUUID string `json:"audio_uuid"`
	Order     int    `json:"order"`
	Text      string `json:"text"`
}

type RequestLyricUpdate struct {
	AudioUUID string `json:"audio_uuid"`
	Order     int    `json:"order"`
	Text      string `json:"text"`
}

type ResponseLyricRead struct {
	UUID      pgtype.UUID        `json:"uuid"`
	AudioUUID pgtype.UUID        `json:"audio_uuid"`
	Order     int                `json:"order"`
	Text      string             `json:"text"`
	CreatedAt pgtype.Timestamptz `json:"created_at" swaggertype:"string"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at" swaggertype:"string"`
}

func (schema *ResponseLyricRead) FromDTO(dto *dto.LyricRead) {
	schema.UUID = dto.UUID
	schema.AudioUUID = dto.AudioUUID
	schema.Order = dto.Order
	schema.Text = dto.Text
	schema.CreatedAt = dto.CreatedAt
	schema.UpdatedAt = dto.UpdatedAt
}
