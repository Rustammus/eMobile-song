package service

import (
	"eMobile/internal/crud"
	"eMobile/internal/dto"
	"eMobile/internal/repo"
	"eMobile/internal/service/audioService"
	"eMobile/internal/service/lyricService"
	"eMobile/pkg/logging"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type Service struct {
	Audio IAudioService
	Lyric ILyricService
}

type Deps struct {
	Repo    repo.Repository
	Logger  logging.Logger
	InfoURL string
}

func NewService(d *Deps) Service {
	return Service{
		Audio: audioService.NewAudioService(&audioService.Deps{
			Repo:    d.Repo,
			Logger:  d.Logger,
			Http:    http.DefaultClient,
			InfoURL: d.InfoURL,
		}),
		Lyric: lyricService.NewLyricService(&lyricService.Deps{
			Repo:   d.Repo,
			Logger: d.Logger,
		}),
	}
}

//go:generate mockgen -source=service.go -destination=mocks\mock.go

type IAudioService interface {
	Create(audio *dto.AudioCreate) (pgtype.UUID, error)
	Find(uuid pgtype.UUID) (*dto.AudioRead, error)
	FindWithLyric(uuid pgtype.UUID) (*dto.AudioReadFull, error)
	ListByFilter(filter *dto.AudioFilter, pag crud.Pagination) ([]dto.AudioRead, error)
	ListPag(pag crud.Pagination) ([]dto.AudioRead, error)
	Update(uuid pgtype.UUID, audio *dto.AudioUpdate) (*dto.AudioRead, error)
	Delete(uuid pgtype.UUID) error
}

type ILyricService interface {
	ListByAudioPag(uuid pgtype.UUID, pag crud.Pagination) ([]dto.LyricRead, error)
}
