package state

import "github.com/spacemeshos/go-spacemesh/sql"

func StartDB(filePath string, connections int) (db *sql.Database, err error) {
	db, err = sql.Open(filePath, sql.WithConnections(connections))
	return
}
