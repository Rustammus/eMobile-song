package schema

import (
	"eMobile/internal/dto"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type AudioInfoRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type ResponseAudioInfo struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (schema *ResponseAudioInfo) ToDTO() (*dto.AudioInfo, error) {
	date, err := time.Parse("02.01.2006", schema.ReleaseDate)
	if err != nil {
		return nil, err
	}

	if schema.Text == "" {
		return nil, errors.New("got empty 'text' in response")
	}
	if schema.Link == "" {
		return nil, errors.New("got empty 'link' in response")
	}

	return &dto.AudioInfo{
		ReleaseDate: pgtype.Date{
			Time:  date,
			Valid: true,
		},
		Text: schema.Text,
		Link: schema.Link,
	}, nil
}
