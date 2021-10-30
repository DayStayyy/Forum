package database

import (
	"database/sql"
)

// ConnectDatabase : Connect / Open to database (database name)
func ConnectDatabase(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// DisconnectDatabase : Disconnect / Close database (*sql.DB)
func DisconnectDatabase(db *sql.DB) error {
	err := db.Close()
	return err
}
