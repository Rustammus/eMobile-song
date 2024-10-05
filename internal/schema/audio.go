package schema

import (
	"database/sql"
	"eMobile/internal/dto"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"net/url"
	"time"
)

type RequestAudioCreate struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

func (schema *RequestAudioCreate) ToDTO() (*dto.AudioCreate, error) {
	errStr := ""
	if schema.Group == "" {
		errStr += "'group' is required and cannot be empty;"
	}
	if schema.Song == "" {
		errStr += "'song' is required and cannot be empty;"
	}

	if errStr != "" {
		return nil, errors.New(errStr)
	}

	return &dto.AudioCreate{
		Group: schema.Group,
		Song:  schema.Song,
	}, nil
}

type RequestAudioUpdate struct {
	Group       *string `json:"group"`
	Song        *string `json:"song"`
	ReleaseDate *string `json:"release_date"`
	Link        *string `json:"link"`
	Lyrics      *string `json:"lyrics"`
}

func (schema *RequestAudioUpdate) ToDTO() (*dto.AudioUpdate, error) {
	dto := &dto.AudioUpdate{}
	count := 0

	errStr := ""

	if schema.Group != nil {
		if *schema.Group == "" {
			errStr += "group cannot be empty;"
		} else {
			dto.Group.String = *schema.Group
			dto.Group.Valid = true
			count++
		}
	}
	if schema.Song != nil {
		if *schema.Song == "" {
			errStr += "song cannot be empty;"
		} else {
			dto.Song.String = *schema.Song
			dto.Song.Valid = true
			count++
		}
	}
	if schema.ReleaseDate != nil {
		t, err := time.Parse("2006-01-02", *schema.ReleaseDate)
		if err != nil {
			errStr += "invalid date format, example: 2006-09-25;"
		} else {
			dto.ReleaseDate.Time = t
			dto.ReleaseDate.Valid = true
			count++
		}
	}
	if schema.Link != nil {
		if *schema.Link == "" {
			errStr += "link cannot be empty;"
		} else {
			dto.Link.String = *schema.Link
			dto.Link.Valid = true
			count++
		}
	}
	if schema.Lyrics != nil {
		if *schema.Lyrics == "" {
			errStr += "lyrics cannot be empty;"
		} else {
			dto.LyricsRaw.String = *schema.Lyrics
			dto.LyricsRaw.Valid = true
			count++
		}
	}
	if errStr != "" {
		return nil, errors.New(errStr)
	} else if count == 0 {
		return nil, errors.New("at least one argument is required")
	}

	return dto, nil
}

type RequestAudioFilter struct {
	Group             string `json:"group"`
	Song              string `json:"song"`
	ReleaseDateAfter  string `json:"after"`
	ReleaseDateBefore string `json:"before"`
	Link              string `json:"link"`
	Lyric             string `json:"lyric"`
}

func (schema *RequestAudioFilter) ToDTO() (*dto.AudioFilter, error) {
	filterDTO := &dto.AudioFilter{
		Group:             sql.NullString{},
		Song:              sql.NullString{},
		ReleaseDateAfter:  pgtype.Date{},
		ReleaseDateBefore: pgtype.Date{},
		Link:              sql.NullString{},
		Lyric:             sql.NullString{},
	}

	empty := true

	if schema.Group != "" {
		filterDTO.Group = sql.NullString{String: schema.Group, Valid: true}
		empty = false
	}
	if schema.Song != "" {
		filterDTO.Song = sql.NullString{String: schema.Song, Valid: true}
		empty = false
	}
	if schema.ReleaseDateAfter != "" {
		t, err := time.Parse("2006-01-02", schema.ReleaseDateAfter)
		if err != nil {
			return nil, errors.New("invalid after format, example: 2006-09-25")
		}
		filterDTO.ReleaseDateAfter = pgtype.Date{Time: t, Valid: true}
		empty = false
	}
	if schema.ReleaseDateBefore != "" {
		t, err := time.Parse("2006-01-02", schema.ReleaseDateBefore)
		if err != nil {
			return nil, errors.New("invalid before format, example: 2006-09-25")
		}
		filterDTO.ReleaseDateBefore = pgtype.Date{Time: t, Valid: true}
		empty = false
	}
	if schema.Link != "" {
		filterDTO.Link = sql.NullString{String: schema.Link, Valid: true}
		empty = false
	}
	if schema.Lyric != "" {
		filterDTO.Lyric = sql.NullString{String: schema.Lyric, Valid: true}
		empty = false
	}
	if empty {
		return nil, nil
	}
	return filterDTO, nil
}

func (schema *RequestAudioFilter) ScanQuery(u *url.URL) {
	q := u.Query()

	schema.Group = q.Get("group")
	schema.Song = q.Get("song")
	schema.ReleaseDateAfter = q.Get("after")
	schema.ReleaseDateBefore = q.Get("before")
	schema.Link = q.Get("link")
	schema.Lyric = q.Get("lyric")
}

type ResponseAudioRead struct {
	UUID        pgtype.UUID        `json:"uuid"`
	Group       string             `json:"group"`
	Song        string             `json:"song"`
	ReleaseDate pgtype.Date        `json:"release_date" swaggertype:"string"`
	Link        string             `json:"link"`
	CreatedAt   pgtype.Timestamptz `json:"created_at" swaggertype:"string"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at" swaggertype:"string"`
}

func (schema *ResponseAudioRead) FromDTO(dto *dto.AudioRead) {
	schema.UUID = dto.UUID
	schema.Group = dto.Group
	schema.Song = dto.Song
	schema.ReleaseDate = dto.ReleaseDate
	schema.Link = dto.Link
	schema.CreatedAt = dto.CreatedAt
	schema.UpdatedAt = dto.UpdatedAt
}

type ResponseAudioReadFull struct {
	UUID        pgtype.UUID         `json:"uuid"`
	Group       string              `json:"group"`
	Song        string              `json:"song"`
	ReleaseDate pgtype.Date         `json:"release_date" swaggertype:"string"`
	Link        string              `json:"link"`
	Lyrics      []ResponseLyricRead `json:"lyrics,omitempty"`
	CreatedAt   pgtype.Timestamptz  `json:"created_at" swaggertype:"string"`
	UpdatedAt   pgtype.Timestamptz  `json:"updated_at" swaggertype:"string"`
}

func (schema *ResponseAudioReadFull) FromDTO(dto *dto.AudioRead) {
	schema.UUID = dto.UUID
	schema.Group = dto.Group
	schema.Song = dto.Song
	schema.ReleaseDate = dto.ReleaseDate
	schema.Link = dto.Link
	schema.CreatedAt = dto.CreatedAt
	schema.UpdatedAt = dto.UpdatedAt
}

func (schema *ResponseAudioReadFull) FromDTOFull(dto *dto.AudioReadFull) {
	lyrics := make([]ResponseLyricRead, 0, len(dto.Lyrics))
	for i := 0; i < len(dto.Lyrics); i++ {
		lyric := ResponseLyricRead{}
		lyric.FromDTO(&dto.Lyrics[i])
		lyrics = append(lyrics, lyric)
	}

	schema.UUID = dto.UUID
	schema.Group = dto.Group
	schema.Song = dto.Song
	schema.ReleaseDate = dto.ReleaseDate
	schema.Link = dto.Link
	schema.Lyrics = lyrics
	schema.CreatedAt = dto.CreatedAt
	schema.UpdatedAt = dto.UpdatedAt
}
