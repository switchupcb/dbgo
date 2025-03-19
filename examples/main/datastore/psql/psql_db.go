package psql

import "log"

// db represents a single instance (singleton) of a database.
var db *Queries

// Database returns a single instance of the database at any point.
func Database() *Queries {
	if db == nil {
		pool, err := Pool()
		if err != nil {
			log.Fatal(err.Error())
		}

		db = New(pool)
	}

	return db
}
