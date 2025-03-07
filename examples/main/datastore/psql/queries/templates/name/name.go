package name

import (
	. "github.com/switchupcb/jet/v2/postgres"
)

// SQL returns return an SQL statement.
//
// You can use Jet to write type-safe SQL queries.
//
// Read https://github.com/go-jet/jet#lets-write-some-sql-queries-in-go for more information.
func SQL() (string, error) {
	stmt := SELECT(Accounts.AllColumns).FROM(Accounts)

	query, _ := stmt.Sql()

	return query, nil
}
