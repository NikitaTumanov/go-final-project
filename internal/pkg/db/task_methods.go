package db

import (
	"database/sql"

	"github.com/NikitaTumanov/go-final-project/internal/pkg/model"
)

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

func TaskByID(id string) (model.Task, error) {
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	var task model.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeate)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func AddTask(date, title, comment, repeate string) (sql.Result, error) {
	return DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeate))
}

func ChangeTaskByID(id, date, title, comment, repeat string) (sql.Result, error) {
	return DB.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", id),
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeat))
}

func DeleteTaskByID(id string) (sql.Result, error) {
	return DB.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
}
