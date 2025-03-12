package repo

import (
	"fmt"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
}

type MemoryRepo struct {
	tasks  map[int]Task
	nextID int
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (r *MemoryRepo) CreateTask(task Task) (int, error) {
	taskID := r.nextID
	task.ID = taskID
	r.tasks[taskID] = task
	r.nextID++
	return taskID, nil
}

func (r *MemoryRepo) GetTaskByID(taskID int) (Task, error) {
	task, ok := r.tasks[taskID]
	if !ok {
		return Task{}, fmt.Errorf("task not found")
	}
	return task, nil
}

func (r *MemoryRepo) GetAllTasks() ([]Task, error) {
	var taskList []Task
	for _, task := range r.tasks {
		taskList = append(taskList, task)
	}
	return taskList, nil
}

func (r *MemoryRepo) UpdateTask(task Task) error {
	if _, ok := r.tasks[task.ID]; !ok {
		return fmt.Errorf("task not found")
	}
	r.tasks[task.ID] = task
	return nil
}

func (r *MemoryRepo) DeleteTask(taskID int) error {
	if _, ok := r.tasks[taskID]; !ok {
		return fmt.Errorf("task not found")
	}
	delete(r.tasks, taskID)
	return nil
}
