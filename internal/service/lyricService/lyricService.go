package lyricService

import (
	"context"
	"eMobile/internal/crud"
	"eMobile/internal/dto"
	"eMobile/internal/repo"
	"eMobile/pkg/logging"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type LyricService struct {
	r repo.Repository
	l logging.Logger
}

type Deps struct {
	Repo   repo.Repository
	Logger logging.Logger
}

func NewLyricService(d *Deps) *LyricService {
	return &LyricService{
		r: d.Repo,
		l: d.Logger,
	}
}

func (s *LyricService) ListByAudioPag(uuid pgtype.UUID, pag crud.Pagination) ([]dto.LyricRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lyrics, err := s.r.Lyric.ListByAudioPag(ctx, uuid, pag)
	return lyrics, err
}

// TODO create, update, delete??
