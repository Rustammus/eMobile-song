package audioService

import (
	"context"
	"eMobile/internal/crud"
	"eMobile/internal/dto"
	"eMobile/internal/repo"
	"eMobile/internal/schema"
	"eMobile/pkg/logging"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AudioService struct {
	r       repo.Repository
	http    *http.Client
	l       logging.Logger
	infoURL string
}

type Deps struct {
	Repo    repo.Repository
	Logger  logging.Logger
	Http    *http.Client
	InfoURL string
}

func NewAudioService(d *Deps) *AudioService {
	return &AudioService{
		r:       d.Repo,
		l:       d.Logger,
		http:    d.Http,
		infoURL: d.InfoURL,
	}
}

func (s *AudioService) Create(audio *dto.AudioCreate) (pgtype.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	audioInfo, err := s.getAudioInfo(ctx, audio.Group, audio.Song)
	if err != nil {
		s.l.Logger.Error("Error getting audio info: ", err)
		return pgtype.UUID{}, err
	}

	lyrics := s.splitAudioText(audioInfo.Text)

	audioFull := &dto.AudioCreateFull{
		Group:       audio.Group,
		Song:        audio.Song,
		ReleaseDate: audioInfo.ReleaseDate,
		Link:        audioInfo.Link,
		Lyrics:      lyrics,
	}

	uuid, err := s.r.Audio.CreateWithLyrics(ctx, audioFull)
	if err != nil {
		s.l.Logger.Error("Error on creating audio: ", err)
	}
	return uuid, err
}

func (s *AudioService) getAudioInfo(ctx context.Context, group, song string) (*dto.AudioInfo, error) {
	queryGroup := url.QueryEscape(group)
	querySong := url.QueryEscape(song)

	infoURL := fmt.Sprintf("%s?group=%s&song=%s", s.infoURL, queryGroup, querySong)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, infoURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("response is nil")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("got non 200 status: " + resp.Status) // todo 400 500 handle
	}

	audioInfo := schema.ResponseAudioInfo{}
	err = json.NewDecoder(resp.Body).Decode(&audioInfo)
	if err != nil {
		return nil, err
	}

	audioInfoDTO, err := audioInfo.ToDTO()
	if err != nil {
		return nil, err
	}

	return audioInfoDTO, nil
}

func (s *AudioService) splitAudioText(text string) []dto.LyricCreate {
	texts := strings.Split(text, "\n\n")
	var lyrics []dto.LyricCreate
	for i := 0; i < len(texts); i++ {
		lyrics = append(lyrics, dto.LyricCreate{
			Order: i,
			Text:  texts[i],
		})
	}
	return lyrics
}

func (s *AudioService) Find(uuid pgtype.UUID) (*dto.AudioRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	audio, err := s.r.Audio.FindByUUID(ctx, uuid)
	if err != nil {
		s.l.Error("Error on finding audio by uuid: ", err)
	}
	return audio, err
}

func (s *AudioService) FindWithLyric(uuid pgtype.UUID) (*dto.AudioReadFull, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	audio, err := s.r.Audio.FindByUUIDWithLyrics(ctx, uuid)
	if err != nil {
		s.l.Error("Error on finding audio with lyrics by uuid: ", err)
	}
	return audio, err
}

func (s *AudioService) ListByFilter(filter *dto.AudioFilter, pag crud.Pagination) ([]dto.AudioRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	audios, err := s.r.Audio.ListByFilter(ctx, filter, pag)
	if err != nil {
		s.l.Error("Error on list audios by filter: ", err)
	}
	return audios, err
}

func (s *AudioService) ListPag(pag crud.Pagination) ([]dto.AudioRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	audios, err := s.r.Audio.ListByPag(ctx, pag)
	if err != nil {
		s.l.Error("Error on list audios: ", err)
	}
	return audios, err
}

func (s *AudioService) Update(uuid pgtype.UUID, audio *dto.AudioUpdate) (*dto.AudioRead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if audio.LyricsRaw.Valid {
		lyrics := make([]dto.LyricCreate, 0)
		lyricsTexts := strings.Split(audio.LyricsRaw.String, "\n\n")
		for i := 0; i < len(lyricsTexts); i++ {
			lyric := dto.LyricCreate{
				AudioUUID: uuid,
				Order:     i,
				Text:      lyricsTexts[i],
			}
			lyrics = append(lyrics, lyric)
		}

		audio.Lyrics = lyrics
	}

	readAudio, err := s.r.Audio.Update(ctx, uuid, audio)
	if err != nil {
		s.l.Error("Error on update audio: ", err)
	}
	return readAudio, err
}

func (s *AudioService) Delete(uuid pgtype.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.r.Audio.Delete(ctx, uuid)
	if err != nil {
		s.l.Error("Error on delete audio: ", err)
	}
	return err
}
