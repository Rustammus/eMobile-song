package dto

import (
	"database/sql"
	"github.com/jackc/pgx/v5/pgtype"
)

type Audio struct {
	UUID        pgtype.UUID        `json:"uuid"`
	Group       string             `json:"group"`
	Song        string             `json:"song"`
	ReleaseDate pgtype.Date        `json:"release_date"`
	Link        string             `json:"link"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type AudioRead struct {
	UUID        pgtype.UUID        `json:"uuid"`
	Group       string             `json:"group"`
	Song        string             `json:"song"`
	ReleaseDate pgtype.Date        `json:"release_date"`
	Link        string             `json:"link"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type AudioReadFull struct {
	UUID        pgtype.UUID        `json:"uuid"`
	Group       string             `json:"group"`
	Song        string             `json:"song"`
	ReleaseDate pgtype.Date        `json:"release_date"`
	Link        string             `json:"link"`
	Lyrics      []LyricRead        `json:"lyrics"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type AudioCreate struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type AudioCreateFull struct {
	Group       string      `json:"group"`
	Song        string      `json:"song"`
	ReleaseDate pgtype.Date `json:"release_date"`
	Link        string      `json:"link"`
	Lyrics      []LyricCreate
}

type AudioInfo struct {
	ReleaseDate pgtype.Date `json:"releaseDate"`
	Text        string      `json:"text"`
	Link        string      `json:"link"`
}

type AudioUpdate struct {
	Group       sql.NullString `json:"group"`
	Song        sql.NullString `json:"song"`
	ReleaseDate pgtype.Date    `json:"release_date"`
	Link        sql.NullString `json:"link"`
	LyricsRaw   sql.NullString `json:"lyric_raw"`
	Lyrics      []LyricCreate
}

type AudioFilter struct {
	Group             sql.NullString `json:"group"`
	Song              sql.NullString `json:"song"`
	ReleaseDateAfter  pgtype.Date    `json:"release_date_after"`
	ReleaseDateBefore pgtype.Date    `json:"release_date_before"`
	Link              sql.NullString `json:"link"`
	Lyric             sql.NullString `json:"lyric"`
}
