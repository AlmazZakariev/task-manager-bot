package storage

import (
	"time"
)

// Task is a type that stores user requests
type Task struct {
	ChatId   int
	UserId   int
	UserName string
	Text     string
	Date     time.Time
}

/*func (t *Task) Id() uint32 {
	return t.id
}
func (t *Task) UserId() uint64 {
	return t.userId
}
func (t *Task) Date() time.Time {
	return t.date
}
func (t *Task) Text() string {
	return t.text
}
func (t *Task) SetId(id uint32) error {
	if false {
		return errors.New("invalid id")
	}
	t.id = id
	return nil
}
func (t *Task) SetUserId(userId uint64) error {
	if false {
		return errors.New("invalid user id")
	}
	t.userId = userId
	return nil
}
func (t *Task) SetDate(date time.Time) error {
	if false {
		return errors.New("invalid date")
	}
	t.date = date
	return nil
}
func (t *Task) SetText(text string) error {
	if false {
		return errors.New("invalid text")
	}
	t.text = text
	return nil
}*/
