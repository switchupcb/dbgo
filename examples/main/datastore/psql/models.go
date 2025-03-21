// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package psql

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        int32
	FirstName pgtype.Text
	LastName  pgtype.Text
	Email     string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type User struct {
	ID        int32
	Name      pgtype.Text
	Password  pgtype.Text
	Email     string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
