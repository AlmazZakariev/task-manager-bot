package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(context.Context, *Task) error

	//get and delete
	GetTasksToExecute(context.Context) (*[]Task, error)

	Remove(context.Context, *Task) error
	IsExists(context.Context, *Task) (bool, error)

	GetAllTasks(context.Context, int) (*[]Task, error)
}

var ErrNoSavedTasks = errors.New("no saved tasks")
