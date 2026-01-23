package db

import (
	"database/sql"
	"os"

	"github.com/NikitaTumanov/go-final-project/internal/pkg/model"
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
	_, err := os.Stat(dbFile)

	install := false
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

func AddTask(date, title, comment, repeate string) (sql.Result, error) {
	return DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeate))
}

func parseTasks(rows *sql.Rows) ([]model.Task, error) {
	tasks := make([]model.Task, 0)
	for rows.Next() {
		task := model.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeate)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func Tasks(limit int) ([]model.Task, error) {
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT :limit",
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseTasks(rows)
}

func TasksByTitleOrComment(search string, limit int) ([]model.Task, error) {
	search = "%" + search + "%"
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
		sql.Named("search", search),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseTasks(rows)
}

func TasksByDate(date string, limit int) ([]model.Task, error) {
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit",
		sql.Named("date", date),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return parseTasks(rows)
}

func TasksByID(id string) (model.Task, error) {
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	var task model.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeate)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func ChangeTaskByID(id, date, title, comment, repeate string) (sql.Result, error) {
	return DB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeate = :repeate WHERE id = :id",
		sql.Named("id", id),
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeate", repeate))
}
