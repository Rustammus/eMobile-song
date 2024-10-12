package repo

import (
	"context"
	"eMobile/internal/crud"
	"eMobile/internal/dto"
	"eMobile/pkg/logging"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository struct {
	Audio AudioRepository
	Lyric LyricRepository
}

// NewRepository
// return all-in-one repository
func NewRepository(c crud.Client, l logging.Logger) Repository {
	return Repository{
		Audio: crud.NewAudioCRUD(c, l),
		Lyric: crud.NewLyricCRUD(c, l),
	}
}

type AudioRepository interface {
	CreateWithLyrics(ctx context.Context, audio *dto.AudioCreateFull) (pgtype.UUID, error)
	ListByPag(ctx context.Context, pag crud.Pagination) ([]dto.AudioRead, error)
	FindByUUID(ctx context.Context, uuid pgtype.UUID) (*dto.AudioRead, error)
	FindByUUIDWithLyrics(ctx context.Context, uuid pgtype.UUID) (*dto.AudioReadFull, error)
	ListByFilter(ctx context.Context, filter *dto.AudioFilter, pag crud.Pagination) ([]dto.AudioRead, error)
	Update(ctx context.Context, uuid pgtype.UUID, audio *dto.AudioUpdate) (*dto.AudioRead, error)
	Delete(ctx context.Context, uuid pgtype.UUID) error
}

type LyricRepository interface {
	ListByAudioPag(ctx context.Context, audioUUID pgtype.UUID, pag crud.Pagination) ([]dto.LyricRead, error)
	FindByUUID(ctx context.Context, uuid pgtype.UUID) (*dto.LyricRead, error)
	DeleteAllByAudio(ctx context.Context, uuid pgtype.UUID) error
}
