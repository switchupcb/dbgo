package example

// SQL returns return an SQL statement.
//
// You can use Jet to write type-safe SQL queries.
//
// Read https://github.com/go-jet/jet#lets-write-some-sql-queries-in-go for more information.
func SQL() (string, error) {
	return "SELECT * FROM accounts", nil
}
