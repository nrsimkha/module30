package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// описаниe модели данных сущностей,
// хранящихся в БД
type user struct {
	id   int
	name string
}

type label struct {
	id   int
	name string
}

type task struct {
	id          int
	opened      int
	closed      int
	author_id   int
	assigned_id int
	title       string
	content     string
}

var ctx context.Context = context.Background()

func main() {
	// Подключение к БД.
	pwd := os.Getenv("dbpass")
	dbpool, err := pgxpool.New(ctx, "postgres://postgres:"+pwd+"@localhost:5432/taskTracker")
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()
	// Подключение подключения
	err = dbpool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//Создать новые задачи,
	taskData := []task{
		{author_id: 1, assigned_id: 2, title: "Complete task", content: "complete task fast"},
	}
	err = addTasks(dbpool, taskData)
	if err != nil {
		log.Fatal(err)
	}

	//Получить список всех задач,
	tasks, err := getTasks(dbpool)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)

	//Получить задачу по id
	taskById, err := getTaskById(dbpool, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(taskById)

	//Получить список задач по автору,
	tasksByAuthor, err := getTasksByAuthor(dbpool, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasksByAuthor)

	//Получить список задач по метке,
	tasksByLabel, err := getTasksByLabel(dbpool, "Bug")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasksByLabel)

	//Обновить задачу по id,
	changedTask := task{author_id: 1, assigned_id: 1, title: "Complete task", content: "complete task fast"}
	err = updateTaskById(dbpool, 3, changedTask)

	//Удалить задачу по id.
	err = deleteTaskById(dbpool, 7)
	if err != nil {
		log.Fatal(err)
	}

}

func addTasks(db *pgxpool.Pool, tasks []task) error {
	_, err := db.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
			closed BIGINT DEFAULT 0,
			author_id INTEGER REFERENCES users(id) DEFAULT 0,
			assigned_id INTEGER REFERENCES users(id) DEFAULT 0,
			title TEXT NOT NULL,
			content TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		_, err := db.Exec(ctx, `INSERT INTO tasks (author_id, assigned_id, title, content) VALUES ($1, $2, $3, $4)`,
			t.author_id, t.assigned_id, t.title, t.content)
		if err != nil {
			return err
		}
	}
	return nil
}

func getTasks(db *pgxpool.Pool) ([]task, error) {
	rows, err := db.Query(ctx, `
		SELECT * FROM tasks ORDER BY id;
	`)
	if err != nil {
		return nil, err
	}
	var tasks []task

	for rows.Next() {
		var t task
		err = rows.Scan(
			&t.id,
			&t.opened,
			&t.closed,
			&t.author_id,
			&t.assigned_id,
			&t.title,
			&t.content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)

	}
	return tasks, rows.Err()
}

func getTaskById(db *pgxpool.Pool, id int) (*task, error) {
	row := db.QueryRow(ctx, `SELECT * FROM tasks WHERE id=($1);`, id)
	var t task
	err := row.Scan(&t.id,
		&t.opened,
		&t.closed,
		&t.author_id,
		&t.assigned_id,
		&t.title,
		&t.content)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func updateTaskById(db *pgxpool.Pool, id int, t task) error {
	_, err := db.Exec(ctx, `UPDATE tasks SET author_id=$1, assigned_id=$2, closed=$3, title=$4, content=$5 WHERE id = $6;`,
		t.author_id, t.assigned_id, t.closed, t.title, t.content, id)
	if err != nil {
		return err
	}
	return nil
}

func deleteTaskById(db *pgxpool.Pool, id int) error {
	_, err := db.Exec(ctx, `DELETE FROM tasks WHERE id = $1;`, id)
	if err != nil {
		return err
	}
	return nil
}

func getTasksByAuthor(db *pgxpool.Pool, author_id int) ([]task, error) {
	rows, err := db.Query(ctx, `SELECT * FROM tasks WHERE author_id=($1);`, author_id)
	var tasks []task

	for rows.Next() {
		var t task
		err = rows.Scan(
			&t.id,
			&t.opened,
			&t.closed,
			&t.author_id,
			&t.assigned_id,
			&t.title,
			&t.content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)

	}
	return tasks, rows.Err()
}

func getTasksByLabel(db *pgxpool.Pool, label string) ([]task, error) {
	rows, err := db.Query(ctx, `SELECT tasks.*
		FROM tasks
		JOIN tasks_labels ON tasks_labels.task_id = tasks.id
		JOIN labels ON labels.id = tasks_labels.label_id
		WHERE labels.name = $1;`, label)
	var tasks []task

	for rows.Next() {
		var t task
		err = rows.Scan(
			&t.id,
			&t.opened,
			&t.closed,
			&t.author_id,
			&t.assigned_id,
			&t.title,
			&t.content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)

	}
	return tasks, rows.Err()
}
