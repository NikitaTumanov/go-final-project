package db

import (
	"database/sql"
	"os"

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

func checkDBFile(dbFile string) bool {
	install := false
	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}

	return install
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

	isNotExists := checkDBFile(dbFile)

	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	isEmpty, err := checkDBIsEmpty(DB)
	if err != nil {
		return err
	}

	if isNotExists || isEmpty {
		_, err := DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return DB.Ping()
}
