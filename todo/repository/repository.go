package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	Done        bool   `json:"done"`
}

func Setup() {
	// Open the database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (id TEXT, title TEXT, description TEXT, dueDate TEXT, done BOOL)")
	if err != nil {
		log.Fatal(err)
	}
}

func GetTodos() ([]Todo, error) {
	// Connect to the database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query the database
	rows, err := db.Query("SELECT id, title, description, dueDate, done FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Read the rows
	todos := []Todo{}
	for rows.Next() {
		todo := Todo{}
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.DueDate, &todo.Done)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	err = rows.Err()

	return todos, err
}

func AddTodo(todo Todo) error {
	// Connect to the database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Query the database
	_, err = db.Exec("INSERT INTO todos (id, title, description, dueDate, done) VALUES (?, ?, ?, ?, ?)", todo.Id, todo.Title, todo.Description, todo.DueDate, todo.Done)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTodo(todo Todo) error {
	// Connect to the database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Query the database
	_, err = db.Exec("DELETE FROM todos WHERE id = ?", todo.Id)
	if err != nil {
		return err
	}

	return nil
}

func GetTodo(id string) (Todo, error) {
	// Connect to the database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return Todo{}, err
	}
	defer db.Close()

	// Query the database
	row := db.QueryRow("SELECT id, title, description, dueDate, done FROM todos WHERE id = ?", id)
	todo := Todo{}
	err = row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.DueDate, &todo.Done)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func UpdateTodo(todo Todo) error {
	// Connect to the database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Query the database
	_, err = db.Exec("UPDATE todos SET title = ?, description = ?, dueDate = ?, done = ? WHERE id = ?", todo.Title, todo.Description, todo.DueDate, todo.Done, todo.Id)
	if err != nil {
		return err
	}

	return nil
}
