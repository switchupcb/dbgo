package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/switchupcb/dbgo/examples/main/datastore/psql"
)

// main.go is a manually developed file which shows you how to use dbgo generated code with
// manually developed database connection and connection pool functions defined in the following files:
//   - ./datastore/psql/psql_connection.go
//   - ./datastore/paql/psql_db.go
func main() {
	// Connect to database for initialization.
	_ = psql.Database()

	// Select the account with an ID = 1 from the database.
	_, err := psql.Database().SelectAccount(
		context.Background(),
		1,
	)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Fatalf("error selecting account in database: %q", err)
	}

	// Insert an account.
	psql.Database().InsertAccount(
		context.Background(),
		psql.InsertAccountParams{
			FirstName: pgtype.Text{
				String: "First",
				Valid:  true,
			},
			LastName: pgtype.Text{
				String: "Last",
				Valid:  true,
			},
			Email: "mail@example.com",
			CreatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
		},
	)

	if err != nil {
		log.Fatalf("error inserting account in database: %q", err)
	}
}
