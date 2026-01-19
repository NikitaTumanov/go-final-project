package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    title TEXT NOT NULL,
    comment TEXT,
    repeat TEXT CHECK (length(repeat) <= 128)
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`

// type Database struct {
// 	db *sql.DB
// }

// func NewDatabase(dbFile string) *Database {
// 	db, err := Init(dbFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return &Database{
// 		db: db,
// 	}
// }

func Init(dbFile string) error {
	envFile := os.Getenv("TODO_DBFILE")
	if envFile != "" {
		dbFile = envFile
	}

	notExists := checkDB(dbFile)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if notExists {
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkDB(dbFile string) bool {
	_, err := os.Stat(dbFile)

	install := false
	if err != nil {
		install = true
	}
	return install
}
