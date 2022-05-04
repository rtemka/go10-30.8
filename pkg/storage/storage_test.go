package storage

import (
	"log"
	"os"
	"testing"
)

var db *Storage

func TestMain(m *testing.M) {
	var err error

	dbUrl := os.Getenv("DATABASE_TEST_URL")
	if dbUrl == "" {
		log.Fatal("environment variable DATABASE_TEST_URL must be set")
	}

	db, err = NewStorage(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()

	err = db.testCleanUp()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	os.Exit(exitCode)
}

func TestReadAllTasks(t *testing.T) {
	tasks, err := db.ReadAllTasks()
	if err != nil {
		t.Fatalf("expected to get tasks, got error [%v]\n", err)
	}

	if len(tasks) != 4 {
		t.Fatalf("expected to get 4 tasks in total, got [%d]\n", len(tasks))
	}
}

func TestCreateTasks(t *testing.T) {
	newTasks := []Task{
		{
			Opened:   0,
			Closed:   0,
			Author:   User{Id: 1},
			Assigned: User{Id: 2},
			Title:    "Test title1",
			Content:  "Test content1",
		},
		{
			Opened:   1,
			Closed:   2,
			Author:   User{Id: 3},
			Assigned: User{Id: 4},
			Title:    "Test title2",
			Content:  "Test content2",
		},
	}
	err := db.CreateTasks(newTasks)
	if err != nil {
		t.Fatalf("expected to create tasks, got error [%v]\n", err)
	}

	tasks, err := db.ReadAllTasks()
	if err != nil {
		t.Fatalf("expected to get tasks, got error [%v]\n", err)
	}

	if len(tasks) != 6 {
		t.Fatalf("expected to get 6 tasks in total, got [%d]\n", len(tasks))
	}
}

func TestReadTaskById(t *testing.T) {
	tasks, err := db.ReadTaskById(1)
	if err != nil {
		t.Fatalf("expected to get tasks, got error [%v]\n", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected to get one task, got [%d]\n", len(tasks))
	}

	expected := Task{
		Id:       1,
		Opened:   1648805555,
		Closed:   1649583155,
		Title:    "Cook meth",
		Content:  `Make 99.1% pure crystals`,
		Author:   User{Id: 3, Name: "Gus Fring"},
		Assigned: User{Id: 1, Name: "Walter White"},
	}

	if tasks[0] != expected {
		t.Fatalf("expected to get %v task, got %v\n", expected, tasks[0])
	}
}

func TestReadTaskByTag(t *testing.T) {
	tasks, err := db.ReadTaskByTag("meth")
	if err != nil {
		t.Fatalf("expected to get tasks, got error [%v]\n", err)
	}

	if len(tasks) != 3 {
		t.Fatalf("expected to get 3 tasks, got [%d]\n", len(tasks))
	}
}

func TestUpdateTaskById(t *testing.T) {
	newTask := Task{
		Id:       1,
		Opened:   0,
		Closed:   0,
		Author:   User{Id: 1, Name: "Walter White"},
		Assigned: User{Id: 2, Name: "Jesse Pinkman"},
		Title:    "Test update title",
		Content:  "Test update content",
	}
	err := db.UpdateTaskById(newTask.Id, newTask)
	if err != nil {
		t.Fatalf("expected to get tasks, got error [%v]\n", err)
	}

	tasks, err := db.ReadTaskById(newTask.Id)
	if err != nil {
		t.Fatalf("expected to get tasks, got error [%v]\n", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected to get one task, got [%d]\n", len(tasks))
	}

	if tasks[0] != newTask {
		t.Fatalf("expected to get updated task %v, got %v\n", newTask, tasks[0])
	}
}

func TestDeleteTaskById(t *testing.T) {
	err := db.DeleteTaskById(1)
	if err != nil {
		t.Fatalf("expected to get tasks, got error [%v]\n", err)
	}

	tasks, err := db.ReadTaskById(1)
	if err != nil {
		t.Fatalf("expected to get tasks, got error [%v]\n", err)
	}

	if len(tasks) != 0 {
		t.Fatalf("expected to get zero tasks, got [%d]\n", len(tasks))
	}
}
