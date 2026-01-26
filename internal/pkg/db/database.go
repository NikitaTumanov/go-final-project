package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    title TEXT NOT NULL,
    comment TEXT,
    repeat TEXT CHECK (length(repeat) <= 128)
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`

func createDir(dbFile string) error {
	dir := filepath.Dir(dbFile)
	return os.MkdirAll(dir, 0755)
}

func checkDBIsEmpty(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%'").Scan(&count)

	return count == 0, err
}

func Init(dbFile string) error {
	envFile := os.Getenv("TODO_DBFILE")
	if envFile != "" {
		dbFile = envFile
	}

	err := createDir(dbFile)
	if err != nil {
		return err
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	isEmpty, err := checkDBIsEmpty(DB)
	if err != nil {
		return err
	}

	if isEmpty {
		_, err := DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return DB.Ping()
}
