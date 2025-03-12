package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"simple-service/internal/dto"
	"simple-service/internal/repo"
	"simple-service/pkg/validator"
	"strconv"
)

type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTaskByID(ctx *fiber.Ctx) error
	GetAllTasks(ctx *fiber.Ctx) error
	UpdateTask(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
}
type Service_new struct {
	log  *zap.Logger
	repo *repo.MemoryRepo
}

func NewService(log *zap.Logger, repo *repo.MemoryRepo) *Service_new {
	return &Service_new{
		repo: repo,
		log:  log,
	}
}

func (s *Service_new) CreateTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	// Десериализация JSON-запроса
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      "new",
	}

	taskID, err := s.repo.CreateTask(task)
	if err != nil {
		s.log.Error("Failed to create task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Ответ
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": taskID},
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *Service_new) GetTaskByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		s.log.Error("Invalid task ID", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID")
	}

	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	if (task == repo.Task{}) {
		s.log.Warn("Task not found", zap.Int("id", taskID))
		return ctx.Status(fiber.StatusNotFound).JSON(dto.Response{
			Status: "error",
			Data:   map[string]string{"message": "Task not found"},
		})
	}

	response := dto.Response{
		Status: "success",
		Data:   task,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *Service_new) GetAllTasks(ctx *fiber.Ctx) error {
	tasks, err := s.repo.GetAllTasks()
	if err != nil {
		s.log.Error("Failed to get tasks", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   tasks,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *Service_new) UpdateTask(ctx *fiber.Ctx) error {
	var req TaskRequest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil || validator.Validate(ctx.Context(), req) != nil {
		s.log.Error("Invalid request", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request")
	}

	taskID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		s.log.Error("Invalid task ID", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID")
	}

	existingTask, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		s.log.Error("Failed to get task by ID", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	if (existingTask == repo.Task{}) {
		s.log.Warn("Task not found", zap.Int("id", taskID))
		return ctx.Status(fiber.StatusNotFound).JSON(dto.Response{
			Status: "error",
			Data:   map[string]string{"message": "Task not found"},
		})
	}

	updatedTask := repo.Task{
		ID:          taskID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "updated",
	}

	if err := s.repo.UpdateTask(updatedTask); err != nil {
		s.log.Error("Failed to update task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.Response{
		Status: "success",
		Data:   map[string]string{"message": "Task updated successfully"},
	})
}

func (s *Service_new) DeleteTask(ctx *fiber.Ctx) error {
	taskID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		s.log.Error("Invalid task ID", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID")
	}

	existingTask, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		s.log.Error("Failed to check if task exists", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	if (existingTask == repo.Task{}) {
		s.log.Warn("Task not found", zap.Int("id", taskID))
		return ctx.Status(fiber.StatusNotFound).JSON(dto.Response{
			Status: "error",
			Data:   map[string]string{"message": "Task not found"},
		})
	}

	if err := s.repo.DeleteTask(taskID); err != nil {
		s.log.Error("Failed to delete task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	return ctx.Status(fiber.StatusOK).JSON(dto.Response{
		Status: "success",
		Data:   map[string]string{"message": "Task deleted successfully"},
	})
}
