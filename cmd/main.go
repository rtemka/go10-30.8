package main

import (
	"fmt"
	"log"
	"os"
	"tasks/pkg/storage"
)

type taskStorage interface {
	CreateTasks(tasks []storage.Task) error

	ReadAllTasks() ([]storage.Task, error)

	ReadTaskByTag(tag string) ([]storage.Task, error)

	UpdateTaskById(id int, t storage.Task) error

	DeleteTaskById(id int) error
}

func main() {
	var err error

	log.SetOutput(os.Stderr)

	connString := os.Getenv("DB_CONN_STRING")
	if connString == "" {
		log.Println("environment variable DB_CONN_STRING must be set")
		os.Exit(1)
	}

	var db taskStorage

	db, err = storage.NewStorage(connString)
	if err != nil {
		log.Printf("while establishing database connection [%v]\n", err)
		os.Exit(1)
	}

	tasks, err := db.ReadAllTasks()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tasks)

}
