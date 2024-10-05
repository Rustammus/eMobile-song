package schema

import "github.com/jackc/pgx/v5/pgtype"

type ResponseUUID struct {
	pgtype.UUID `json:"uuid"`
}
