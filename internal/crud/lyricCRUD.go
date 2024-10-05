package crud

import (
	"context"
	"eMobile/internal/dto"
	"eMobile/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type LyricCRUD struct {
	c Client
	l logging.Logger
}

func NewLyricCRUD(c Client, l logging.Logger) *LyricCRUD {
	return &LyricCRUD{c: c, l: l}
}

func (c *LyricCRUD) ListByAudioPag(ctx context.Context, audioUUID pgtype.UUID, pag Pagination) ([]dto.LyricRead, error) {
	q := `SELECT uuid, audio_uuid, "order", text, created_at, updated_at
		  FROM public.lyrics
		  WHERE audio_uuid = $1
		  ORDER BY "order"
		  LIMIT $2 OFFSET $3`

	lyrics := make([]dto.LyricRead, 0, pag.Limit)
	rows, err := c.c.Query(ctx, q, audioUUID, pag.Limit, pag.Offset)
	defer rows.Close()
	if err != nil {
		return lyrics, err
	}
	for rows.Next() {
		lyric := dto.LyricRead{}
		err = rows.Scan(&lyric.UUID, &lyric.AudioUUID, &lyric.Order, &lyric.Text, &lyric.CreatedAt, &lyric.UpdatedAt)
		if err != nil {
			return nil, err
		}
		lyrics = append(lyrics, lyric)
	}
	if len(lyrics) == 0 {
		return nil, pgx.ErrNoRows
	}

	return lyrics, nil
}

func (c *LyricCRUD) FindByUUID(ctx context.Context, uuid pgtype.UUID) (*dto.LyricRead, error) {
	q := `SELECT uuid, audio_uuid, "order", text, created_at, updated_at
		  FROM public.lyrics
		  WHERE uuid = $1`
	lyric := &dto.LyricRead{}
	err := c.c.QueryRow(ctx, q, uuid).Scan(lyric.UUID, lyric.AudioUUID, lyric.Order, lyric.Text, lyric.CreatedAt, lyric.UpdatedAt)

	return lyric, err
}

func (c *LyricCRUD) DeleteAllByAudio(ctx context.Context, uuid pgtype.UUID) error {
	q := `DELETE FROM public.lyrics WHERE audio_uuid=$1`
	tag, err := c.c.Exec(ctx, q, uuid)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
