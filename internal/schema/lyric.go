package schema

import (
	"eMobile/internal/dto"
	"github.com/jackc/pgx/v5/pgtype"
)

type RequestLyricCreate struct {
	AudioUUID string `json:"audio_uuid" example:"da6f6e2c-ef5d-4276-b0a1-5067e77278ca"`
	Order     int    `json:"order" example:"0"`
	Text      string `json:"text" example:"Never gonna give you up"`
}

type RequestLyricUpdate struct {
	AudioUUID string `json:"audio_uuid" example:"da6f6e2c-ef5d-4276-b0a1-5067e77278ca"`
	Order     int    `json:"order"`
	Text      string `json:"text"`
}

type ResponseLyricRead struct {
	UUID      pgtype.UUID        `json:"uuid" example:"da6f6e2c-ef5d-4276-b0a1-5067e77278ca"`
	AudioUUID pgtype.UUID        `json:"audio_uuid" example:"da6f6e2c-ef5d-4276-b0a1-5067e77278ca"`
	Order     int                `json:"order" example:"0"`
	Text      string             `json:"text" example:"Never gonna give you up"`
	CreatedAt pgtype.Timestamptz `json:"created_at" swaggertype:"string" example:"2024-10-05T12:57:19.752+05:00"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at" swaggertype:"string" example:"2024-10-05T12:57:19.752+05:00"`
}

func (schema *ResponseLyricRead) FromDTO(dto *dto.LyricRead) {
	schema.UUID = dto.UUID
	schema.AudioUUID = dto.AudioUUID
	schema.Order = dto.Order
	schema.Text = dto.Text
	schema.CreatedAt = dto.CreatedAt
	schema.UpdatedAt = dto.UpdatedAt
}
