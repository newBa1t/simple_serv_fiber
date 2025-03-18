package repo

import (
	"fmt"
	"sync"
)

type MemoryRepo struct {
	tasks  map[int]Task
	nextID int
	mu     sync.Mutex
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (r *MemoryRepo) CreateTask(task Task) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	taskID := r.nextID
	task.ID = taskID
	r.tasks[taskID] = task
	r.nextID++
	return taskID, nil
}

func (r *MemoryRepo) GetTaskByID(taskID int) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[taskID]
	if !ok {
		return Task{}, fmt.Errorf("task not found")
	}
	return task, nil
}

func (r *MemoryRepo) GetAllTasks() ([]Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var taskList []Task
	for _, task := range r.tasks {
		taskList = append(taskList, task)
	}
	return taskList, nil
}

func (r *MemoryRepo) UpdateTask(task Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[task.ID]; !ok {
		return fmt.Errorf("task not found")
	}
	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryRepo) DeleteTask(taskID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[taskID]; !ok {
		return fmt.Errorf("task not found")
	}
	delete(r.tasks, taskID)
	return nil
}
