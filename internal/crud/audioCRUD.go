package crud

import (
	"context"
	"eMobile/internal/dto"
	"eMobile/pkg/logging"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
	"strings"
)

type AudioCRUD struct {
	db     Client
	logger logging.Logger
}

func NewAudioCRUD(c Client, l logging.Logger) *AudioCRUD {
	return &AudioCRUD{db: c, logger: l}
}

func (c *AudioCRUD) CreateWithLyrics(ctx context.Context, audio *dto.AudioCreateFull) (pgtype.UUID, error) {
	qAudio := `INSERT INTO public.audios 
    	  ("group", song, release_date, link, created_at, updated_at)
    	  VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP(3), CURRENT_TIMESTAMP(3))
    	  RETURNING uuid;`

	// insert audio
	uuid := pgtype.UUID{}
	trx, err := c.db.Begin(ctx)
	if err != nil {
		return uuid, err
	}
	defer trx.Rollback(ctx)

	err = trx.QueryRow(ctx, qAudio, audio.Group, audio.Song, audio.ReleaseDate, audio.Link).Scan(&uuid)
	if err != nil {
		return uuid, err
	}

	err = c.insertLyrics(ctx, trx, uuid, audio.Lyrics)
	if err != nil {
		return uuid, err
	}

	return uuid, trx.Commit(ctx)
}

func (c *AudioCRUD) insertLyrics(ctx context.Context, trx pgx.Tx, audioUUID pgtype.UUID, lyrics []dto.LyricCreate) error {
	qLyrics := `INSERT INTO public.lyrics
				(audio_uuid, "order", text, created_at, updated_at)
				VALUES `

	// prepare lyrics query
	counter := 1
	qValues := make([]string, 0, len(lyrics))
	values := make([]any, 0, len(lyrics))
	for i := 0; i < len(lyrics); i++ {
		qValue := fmt.Sprintf("($%d, $%d, $%d, CURRENT_TIMESTAMP(3), CURRENT_TIMESTAMP(3))",
			counter, counter+1, counter+2)
		counter += 3
		qValues = append(qValues, qValue)

		lyric := &lyrics[i]
		values = append(values, audioUUID, lyric.Order, lyric.Text)
	}

	qLyrics += strings.Join(qValues, ",") + ";"

	// insert lyrics
	_, err := trx.Exec(ctx, qLyrics, values...)
	return err
}

func (c *AudioCRUD) ListByPag(ctx context.Context, pag Pagination) ([]dto.AudioRead, error) {
	q := `SELECT uuid, "group", song, release_date, link, created_at, updated_at
		  FROM public.audios
		  LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(ctx, q, pag.Limit, pag.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	audios := make([]dto.AudioRead, 0, pag.Limit)
	for rows.Next() {
		a := dto.AudioRead{}
		err = rows.Scan(&a.UUID, &a.Group, &a.Song, &a.ReleaseDate, &a.Link, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		audios = append(audios, a)
	}
	return audios, rows.Err()
}

func (c *AudioCRUD) FindByUUID(ctx context.Context, uuid pgtype.UUID) (*dto.AudioRead, error) {
	q := `SELECT uuid, "group", song, release_date, link, created_at, updated_at
		  FROM public.audios 
		  WHERE uuid=$1`
	a := dto.AudioRead{}
	err := c.db.QueryRow(ctx, q, uuid).Scan(&a.UUID, &a.Group, &a.Song, &a.ReleaseDate, &a.Link, &a.CreatedAt, &a.UpdatedAt)

	return &a, err
}

func (c *AudioCRUD) FindByUUIDWithLyrics(ctx context.Context, uuid pgtype.UUID) (*dto.AudioReadFull, error) {
	// select audio
	qAudio := `SELECT uuid, "group", song, release_date, link, created_at, updated_at
		  FROM public.audios 
		  WHERE uuid=$1`

	a := dto.AudioReadFull{}
	err := c.db.QueryRow(ctx, qAudio, uuid).Scan(&a.UUID, &a.Group, &a.Song, &a.ReleaseDate, &a.Link, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// select lyrics rows
	qLyrics := `SELECT uuid, audio_uuid, "order", text, created_at, updated_at
				FROM public.lyrics
				WHERE audio_uuid = $1
		  		ORDER BY "order"`

	lyrics := make([]dto.LyricRead, 0)
	rows, err := c.db.Query(ctx, qLyrics, uuid)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		lyric := dto.LyricRead{}
		err = rows.Scan(&lyric.UUID, &lyric.AudioUUID, &lyric.Order, &lyric.Text, &lyric.CreatedAt, &lyric.UpdatedAt)
		if err != nil {
			return nil, err
		}
		lyrics = append(lyrics, lyric)
	}
	a.Lyrics = lyrics

	return &a, rows.Err()
}

func (c *AudioCRUD) ListByFilter(ctx context.Context, filter *dto.AudioFilter, pag Pagination) ([]dto.AudioRead, error) {
	var baseQuery string
	if filter.Lyric.Valid {
		baseQuery = `SELECT DISTINCT a.uuid, a."group", a.song, a.release_date, a.link, a.created_at, a.updated_at
					  FROM public.audios a
					  JOIN public.lyrics l ON a.uuid = l.audio_uuid
					  WHERE `
	} else {
		baseQuery = `SELECT uuid, "group", song, release_date, link, created_at, updated_at 
				 	  FROM public.audios a
				 	  WHERE `
	}

	endQuery := ` LIMIT $1 OFFSET $2;`

	values := make([]any, 0, 8)
	values = append(values, pag.Limit, pag.Offset)

	q, values := c.buildWhereQuery(baseQuery, filter, values)
	q += endQuery

	rows, err := c.db.Query(ctx, q, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	audios := make([]dto.AudioRead, 0, pag.Limit)
	for rows.Next() {
		a := dto.AudioRead{}
		err = rows.Scan(&a.UUID, &a.Group, &a.Song, &a.ReleaseDate, &a.Link, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		audios = append(audios, a)
	}
	return audios, rows.Err()
}

func (c *AudioCRUD) buildWhereQuery(base string, filter *dto.AudioFilter, values []any) (string, []any) {
	conditions := make([]string, 0, 6)
	counter := 3

	if filter.Group.Valid {
		conditions = append(conditions, "a.group = $"+strconv.Itoa(counter))
		values = append(values, filter.Group.String)
		counter++
	}
	if filter.ReleaseDateAfter.Valid {
		conditions = append(conditions, "a.release_date >= $"+strconv.Itoa(counter))
		values = append(values, filter.ReleaseDateAfter)
		counter++
	}
	if filter.ReleaseDateBefore.Valid {
		conditions = append(conditions, "a.release_date <= $"+strconv.Itoa(counter))
		values = append(values, filter.ReleaseDateBefore)
		counter++
	}
	if filter.Link.Valid {
		conditions = append(conditions, "a.link = $"+strconv.Itoa(counter))
		values = append(values, filter.Link.String)
		counter++
	}
	if filter.Song.Valid {
		conditions = append(conditions, "to_tsvector('english', a.song) @@ phraseto_tsquery('english', $"+
			strconv.Itoa(counter)+")")
		values = append(values, filter.Song.String)
		counter++
	}
	if filter.Lyric.Valid {
		conditions = append(conditions, "to_tsvector('english', l.text) @@ phraseto_tsquery('english', $"+
			strconv.Itoa(counter)+")")
		values = append(values, filter.Lyric.String)
		counter++
	}

	base += strings.Join(conditions, " AND ")
	return base, values
}

func (c *AudioCRUD) Update(ctx context.Context, uuid pgtype.UUID, audio *dto.AudioUpdate) (*dto.AudioRead, error) {
	baseQuery := `UPDATE public.audios 
		  SET (updated_at, %s) = ROW(CURRENT_TIMESTAMP(3), %s) 
		  WHERE uuid=$1 
		  RETURNING uuid, "group", song, release_date, link, created_at, updated_at;`

	var values []any
	values = append(values, uuid)
	q, values := c.buildUpdateQuery(baseQuery, values, audio)

	trx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer trx.Rollback(ctx)

	// Update audio
	rAudio := dto.AudioRead{}
	err = trx.QueryRow(ctx, q, values...).
		Scan(&rAudio.UUID, &rAudio.Group, &rAudio.Song, &rAudio.ReleaseDate, &rAudio.Link, &rAudio.CreatedAt, &rAudio.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Delete all lyrics and insert new
	if audio.Lyrics != nil && len(audio.Lyrics) > 0 {
		lyricsDelQuery := `DELETE FROM public.lyrics WHERE audio_uuid=$1;`
		_, err = trx.Exec(ctx, lyricsDelQuery, rAudio.UUID)
		if err != nil {
			return nil, err
		}

		err = c.insertLyrics(ctx, trx, rAudio.UUID, audio.Lyrics)
		if err != nil {
			return nil, err
		}
	}

	return &rAudio, trx.Commit(ctx)
}

func (c *AudioCRUD) buildUpdateQuery(base string, values []any, audio *dto.AudioUpdate) (string, []any) {
	var names []string
	var ids []string
	count := 2

	if audio.Group.Valid {
		names = append(names, "\"group\"")
		ids = append(ids, "$"+strconv.Itoa(count))
		values = append(values, audio.Group.String)
		count++
	}
	if audio.Song.Valid {
		names = append(names, "song")
		ids = append(ids, "$"+strconv.Itoa(count))
		values = append(values, audio.Song.String)
		count++
	}
	if audio.ReleaseDate.Valid {
		names = append(names, "release_date")
		ids = append(ids, "$"+strconv.Itoa(count))
		values = append(values, audio.ReleaseDate)
		count++
	}
	if audio.Link.Valid {
		names = append(names, "link")
		ids = append(ids, "$"+strconv.Itoa(count))
		values = append(values, audio.Link.String)
		count++
	}

	q := fmt.Sprintf(base, strings.Join(names, ","), strings.Join(ids, ","))
	return q, values
}

func (c *AudioCRUD) Delete(ctx context.Context, uuid pgtype.UUID) error {
	q := "DELETE FROM public.audios WHERE uuid=$1"

	tag, err := c.db.Exec(ctx, q, uuid)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
