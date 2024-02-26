package memory

import (
	"errors"
	"sort"
	"tasks-manager-bot/storage"
	"time"
)

type MemoryStorage struct {
	data []storage.Task
}

func New() *MemoryStorage {
	return &MemoryStorage{}
}
func (m MemoryStorage) Len() int {
	return len(m.data)
}
func (m MemoryStorage) Less(i, j int) bool {
	a := m.data[i].Date.Compare(m.data[j].Date)
	return a == -1
}
func (m MemoryStorage) Swap(i, j int) {
	temp := m.data[i]

	m.data[i] = m.data[j]

	m.data[j] = temp

}
func (m *MemoryStorage) Save(task *storage.Task) error {
	m.data = append(m.data, *task)
	return nil
}
func (m MemoryStorage) GetTenFresh() (*[]storage.Task, error) {
	tasks := make([]storage.Task, 0)
	sort.Sort(m)
	for i := 0; i < len(m.data) && len(tasks) <= 10; i++ {
		tasks = append(tasks, m.data[i])
	}
	return &tasks, nil
}
func (m *MemoryStorage) Remove(task *storage.Task) error {
	index := -1
	for i, value := range m.data {
		if value == *task {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("task now found")
	}

	m.data[index] = m.data[len(m.data)-1]
	m.data = m.data[:len(m.data)-1]
	return nil
}
func (m MemoryStorage) IsExists(task *storage.Task) (bool, error) {
	//not need implementation in memory storage
	return false, nil
}
func (m *MemoryStorage) GetTasksToExecute() (*[]storage.Task, error) {
	var tasks []storage.Task

	for _, task := range m.data {

		if task.Date.Before(time.Now()) {
			tasks = append(tasks, task)
			m.Remove(&task)
		}
	}

	return &tasks, nil

}
