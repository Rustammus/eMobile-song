package schema

import "github.com/jackc/pgx/v5/pgtype"

type ResponseUUID struct {
	pgtype.UUID `json:"uuid" example:"da6f6e2c-ef5d-4276-b0a1-5067e77278ca"`
}
