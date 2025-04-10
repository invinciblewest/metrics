package pgstorage

import "database/sql"

func InstallSchema(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS metrics (
		id TEXT NOT NULL,
		type TEXT NOT NULL,
		value DOUBLE PRECISION
	);
	CREATE UNIQUE INDEX IF NOT EXISTS unique_id_type ON metrics (id, type);
	`)
	if err != nil {
		return err
	}
	return nil
}
