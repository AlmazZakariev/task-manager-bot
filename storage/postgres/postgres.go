package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tasks-manager-bot/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	//TODO: может быть ошибка когда бд ещё не запущена
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, task *storage.Task) error {
	q := `INSERT INTO tasks (user_id, chat_id, text, date, user_name) VALUES ($1,$2,$3,$4,$5)`

	_, err := s.db.ExecContext(
		ctx,
		q,
		task.UserId,
		task.ChatId,
		task.Text,
		task.Date,
		task.UserName,
	)

	if err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}
	return nil

}
func (s *Storage) GetTasksToExecute(ctx context.Context) (*[]storage.Task, error) {
	q := `SELECT chat_id, user_id, user_name, text, date FROM tasks WHERE date <= current_timestamp`

	tasks, err := s.selectQuery(s.db.QueryContext(ctx, q))
	if err != nil {
		return nil, fmt.Errorf("can't get tasks to execute: %w", err)
	}
	for _, val := range *tasks {
		s.Remove(ctx, &val)
	}
	return tasks, nil
}

func (s *Storage) Remove(ctx context.Context, task *storage.Task) error {
	q := `DELETE FROM tasks WHERE chat_id = $1 AND text = $2 AND date = $3`

	if _, err := s.db.ExecContext(ctx, q, task.ChatId, task.Text, task.Date); err != nil {
		return fmt.Errorf("can't remove tasks")
	}
	return nil
}

func (s *Storage) IsExists(ctx context.Context, task *storage.Task) (bool, error) {
	//TODO
	return false, nil
}

func (s *Storage) GetAllTasks(ctx context.Context, userID int) (*[]storage.Task, error) {
	q := `SELECT chat_id, user_id, user_name, text, date FROM tasks WHERE user_id = $1 ORDER BY date ASC`
	return s.selectQuery(s.db.QueryContext(ctx, q, userID))
}

func (s *Storage) selectQuery(rows *sql.Rows, err error) (*[]storage.Task, error) {
	var tasks []storage.Task

	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("can't pick tasks to execute: %w", err)
	}

	if rows.Err() == sql.ErrNoRows {
		return nil, nil
	}

	for rows.Next() {
		var task storage.Task

		err := rows.Scan(&task.ChatId, &task.UserId, &task.UserName, &task.Text, &task.Date)
		if err != nil {
			return nil, fmt.Errorf("can't scan task to execute: %w", err)
		}
		tasks = append(tasks, task)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("can't pick2 task to execute: %w", err)
	}
	return &tasks, nil
}
func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS tasks (user_id INTEGER, chat_id INTEGER, user_name TEXT, text TEXT, date timestamp)`
	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}
	return nil
}
