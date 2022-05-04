package storage

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Id   int
	Name string
}

type Task struct {
	Id       int
	Opened   int64
	Closed   int64
	Author   User
	Assigned User
	Title    string
	Content  string
}

type Storage struct {
	db *pgxpool.Pool
}

// NewStorage выполняет подключение
// и возвращает объект для взаимодействия с БД
func NewStorage(connString string) (*Storage, error) {

	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

// CreateTasks пакетно создает задачи в БД
func (s *Storage) CreateTasks(tasks []Task) error {
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	batch := new(pgx.Batch)

	stmt := `
		INSERT INTO tasks(opened, closed, title, content, author_id, assigned_id)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	for i := range tasks {
		batch.Queue(stmt, tasks[i].Opened, tasks[i].Closed, tasks[i].Title,
			tasks[i].Content, tasks[i].Author.Id, tasks[i].Assigned.Id)
	}

	res := tx.SendBatch(context.Background(), batch)
	err = res.Close()
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

// ReadAllTasks возвращает список всех задач
func (s *Storage) ReadAllTasks() ([]Task, error) {

	stmt := `
		SELECT 
			t.id, 
			t.opened, 
			t.closed, 
			t.title, 
			t.content, 
			u.name,
			u.id,  
			u2.name,
			u2.id
		FROM
			tasks AS t INNER JOIN users AS u ON t.author_id = u.id 
			INNER JOIN users AS u2 ON t.assigned_id = u2.id;
	`

	return s.readTasks(stmt)
}

// ReadTaskById возвращает список задач по id
func (s *Storage) ReadTaskById(id int) ([]Task, error) {
	stmt := `
		SELECT 
			t.id, 
			t.opened, 
			t.closed, 
			t.title, 
			t.content, 
			u.name,
			u.id,  
			u2.name,
			u2.id
		FROM
			tasks AS t INNER JOIN users AS u ON t.author_id = u.id 
			INNER JOIN users AS u2 ON t.assigned_id = u2.id
		WHERE t.id = $1;
	`

	return s.readTasks(stmt, id)
}

// ReadTaskByTag возвращает список задач по метке
func (s *Storage) ReadTaskByTag(tag string) ([]Task, error) {

	stmt := `
		SELECT 
			t.id, 
			t.opened, 
			t.closed, 
			t.title, 
			t.content, 
			u.name,
			u.id,  
			u2.name,
			u2.id
		FROM
			labels as l INNER JOIN tasks_labels ON tasks_labels.label_id = l.id
			INNER JOIN tasks as t ON t.id = tasks_labels.task_id
			INNER JOIN users AS u ON t.author_id = u.id 
			INNER JOIN users AS u2 ON t.assigned_id = u2.id
		WHERE l.name = $1;
	`

	return s.readTasks(stmt, tag)
}

// readTasks вспомогательная функция для чтения задач из БД,
// выполняет запрос согласно переданному тексту запроса и аргументам,
// возвращает список задач
func (s *Storage) readTasks(stmt string, args ...interface{}) ([]Task, error) {
	var tasks []Task

	rows, err := s.db.Query(context.Background(), stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.Id, &t.Opened, &t.Closed, &t.Title, &t.Content,
			&t.Author.Name, &t.Author.Id, &t.Assigned.Name, &t.Assigned.Id)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

// UpdateTaskById обновялет задачу по переданному id
func (s *Storage) UpdateTaskById(id int, t Task) error {

	stmt := `
		UPDATE tasks
		SET opened = $2,
			closed = $3,
			title = $4,
			content = $5,
			author_id = $6,
			assigned_id = $7
		WHERE id = $1;
	`
	ctx := context.Background()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, stmt, id, t.Opened, t.Closed,
		t.Title, t.Content, t.Author.Id, t.Assigned.Id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeleteTaskById удаляет задачу по переданному id
func (s *Storage) DeleteTaskById(id int) error {

	// запрос для удаления информации
	// из ссылочной таблицы
	stmtLabels := `
		DELETE FROM tasks_labels
		WHERE tasks_labels.task_id = $1;
	`
	// основной запрос для удаления
	stmtTask := `
		DELETE FROM tasks
		WHERE tasks.id = $1;
	`
	ctx := context.Background()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// удаляем ссылки
	_, err = tx.Exec(ctx, stmtLabels, id)
	if err != nil {
		return err
	}

	// удаляем задачу
	_, err = tx.Exec(ctx, stmtTask, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Storage) testCleanUp() error {

	//fp := filepath.Join("..", "pkg", "storage", "testCleanUp.sql")

	b, err := os.ReadFile("testCleanUp.sql")
	if err != nil {
		return err
	}

	ctx := context.Background()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, string(b))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
