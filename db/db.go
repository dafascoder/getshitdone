package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID          int
	Name        string
	Description string
	Project     string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Status int

const (
	Todo Status = iota + 1
	InProgress
	Done
)

type TaskDB struct {
	database *sql.DB
}

func OpenDB() (*TaskDB, error) {
	database, err := sql.Open("sqlite3", "tasks.db")
	if err != nil {
		return nil, err
	}

	t := TaskDB{database}
	t.CreateTable()

	return &t, nil
}

func (t *TaskDB) CreateTable() {

	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
	"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"name" TEXT,
	"description" TEXT,
	"status" TEXT,
	"project" TEXT,
	"createdAt" TIMESTAMP,
	"updatedAt" TIMESTAMP
	);`

	statement, err := t.database.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec()
}

func InitTaskDir(taskDir string) error {
	if _, err := os.Stat(taskDir); os.IsNotExist(err) {
		err := os.MkdirAll(taskDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TaskDB) Close() {
	t.database.Close()
}

func (s Status) String() string {
	return [...]string{"Not Started", "In Progress", "Done"}[s-1]
}

func (t *TaskDB) AddTask(name, project string, description string) error {
	// We don't care about the returned values, so we're using Exec. If we
	// wanted to reuse these statements, it would be more efficient to use
	// prepared statements. Learn more:
	// https://go.dev/doc/database/prepared-statements
	log.Println("Adding task to the database")
	_, err := t.database.Exec(
		"INSERT INTO tasks(name,project,description,status,createdAt,updatedAt) VALUES(?, ?, ?, ?, ?, ?)",
		name,
		description,
		project,
		Todo.String(),
		time.Now(),
		time.Now())
	return err
}

func (t *TaskDB) GetTask(id int) (Task, error) {
	row := t.database.QueryRow("SELECT * FROM tasks WHERE id=?", id)
	task := Task{}
	err := row.Scan(&task.ID, &task.Name, &task.Description, &task.Project, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	return task, err
}

func (t *TaskDB) UpdateTask(task Task) error {
	log.Println("Updating task in the database")
	_, err := t.database.Exec(
		"UPDATE tasks SET name=?, project=?, description=?, status=?, updatedAt=? WHERE id=?",
		task.Name,
		task.Project,
		task.Description,
		task.Status,
		time.Now(),
		task.ID)
	return err
}
func (t *TaskDB) DeleteTask(id int) error {
	log.Println("Deleting task from the database")
	_, err := t.database.Exec("DELETE FROM tasks WHERE id=?", id)
	return err
}

func (t *TaskDB) GetTasks() ([]Task, error) {
	var tasks []Task
	rows, err := t.database.Query("SELECT * FROM tasks")
	if err != nil {
		return tasks, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var task Task
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Project,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}
