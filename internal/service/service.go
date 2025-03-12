package service

//
//import (
//	"encoding/json"
//	"github.com/gofiber/fiber/v2"
//	"go.uber.org/zap"
//	"simple-service/internal/dto"
//	"simple-service/internal/repo"
//	"simple-service/pkg/validator"
//	"strconv"
//)
//
//// Слой бизнес-логики. Тут должна быть основная логика сервиса
//
//// Service - интерфейс для бизнес-логики
//type Service interface {
//	CreateTask(ctx *fiber.Ctx) error
//	GetTaskByID(ctx *fiber.Ctx) error // Новый метод
//	GetAllTasks(ctx *fiber.Ctx) error
//	UpdateTask(ctx *fiber.Ctx) error
//	DeleteTask(ctx *fiber.Ctx) error
//}
//
//type service struct {
//	repo repo.Repository
//	log  *zap.SugaredLogger
//}
//
//// NewService - конструктор сервиса
//func NewService(repo repo.Repository, logger *zap.SugaredLogger) Service {
//	return &service{
//		repo: repo,
//		log:  logger,
//	}
//}
//
//// CreateTask - обработчик запроса на создание задачи
//func (s *service) CreateTask(ctx *fiber.Ctx) error {
//	var req TaskRequest
//
//	// Десериализация JSON-запроса
//	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
//		s.log.Error("Invalid request body", zap.Error(err))
//		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
//	}
//
//	// Валидация входных данных
//	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
//		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
//	}
//
//	// Вставка задачи в БД через репозиторий
//	task := repo.Task{
//		Title:       req.Title,
//		Description: req.Description,
//		Status:      "new",
//	}
//	taskID, err := s.repo.CreateTask(ctx.Context(), task)
//
//	//// 500 Internal Server Error
//	if err != nil {
//		s.log.Error("Failed to insert task", zap.Error(err))
//		return dto.InternalServerError(ctx)
//	}
//
//	// Формирование ответа 200
//	response := dto.Response{
//		Status: "success",
//		Data:   map[string]int{"task_id": taskID},
//	}
//
//	return ctx.Status(fiber.StatusOK).JSON(response)
//}
//
//// GetTaskByID - обработчик запроса на получение задачи по ID
//func (s *service) GetTaskByID(ctx *fiber.Ctx) error {
//	id := ctx.Params("id") // Получаем ID из параметров маршрута
//
//	// Преобразуем ID в int
//	taskID, err := strconv.Atoi(id)
//
//	if err != nil {
//		s.log.Error("Invalid task ID", zap.Error(err))
//		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID")
//	}
//
//	// Получаем задачу из репозитория
//	task, err := s.repo.GetTaskByID(ctx.Context(), taskID)
//
//	// 500 Internal Server Error
//	if err != nil {
//		s.log.Error("Failed to get task by ID", zap.Error(err))
//		return dto.InternalServerError(ctx)
//	}
//
//	// 404
//	if task == nil {
//		s.log.Error("Task not found", zap.String("id", id))
//		return ctx.Status(fiber.StatusNotFound).JSON(dto.Response{
//			Status: "error",
//			Data:   map[string]string{"message": "Task not found"},
//		})
//	}
//
//	// Формируем ответ
//	response := dto.Response{
//		Status: "success",
//		Data:   task, // Отправляем найденную задачу в ответе
//	}
//
//	return ctx.Status(fiber.StatusOK).JSON(response)
//}
//
//func (s *service) GetAllTasks(ctx *fiber.Ctx) error {
//	tasks, err := s.repo.GetAllTasks(ctx.Context())
//
//	// 500 Internal Server Error
//	if err != nil {
//		s.log.Error("Failed to get tasks", zap.Error(err))
//		return dto.InternalServerError(ctx)
//	}
//
//	//200
//	response := dto.Response{
//		Status: "success",
//		Data:   tasks,
//	}
//	return ctx.Status(fiber.StatusOK).JSON(response)
//}
//
//func (s *service) UpdateTask(ctx *fiber.Ctx) error {
//	var req TaskRequest
//
//	if err := json.Unmarshal(ctx.Body(), &req); err != nil || validator.Validate(ctx.Context(), req) != nil {
//		s.log.Error("Invalid request", zap.Error(err))
//		//404
//		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request")
//	}
//
//	taskID, err := strconv.Atoi(ctx.Params("id"))
//	if err != nil {
//		s.log.Error("Invalid task ID", zap.Error(err))
//		//404
//		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID")
//	}
//
//	existingTask, err := s.repo.GetTaskByID(ctx.Context(), taskID)
//	if err != nil {
//		s.log.Error("Failed to get task by ID", zap.Error(err))
//		//505
//		return dto.InternalServerError(ctx)
//	}
//	if existingTask == nil {
//		//404
//		s.log.Warn("Task not found", zap.Int("id", taskID))
//		return ctx.Status(fiber.StatusNotFound).JSON(dto.Response{
//			Status: "error",
//			Data:   map[string]string{"message": "Task not found"},
//		})
//	}
//
//	if err := s.repo.UpdateTask(ctx.Context(), repo.Task{
//		ID:          taskID,
//		Title:       req.Title,
//		Description: req.Description,
//		Status:      "Update data",
//	});
//	// 500 Internal Server Error
//	err != nil {
//		s.log.Error("Failed to update task", zap.Error(err))
//		return dto.InternalServerError(ctx)
//	}
//
//	//200
//	return ctx.Status(fiber.StatusOK).JSON(dto.Response{Status: "Update data", Data: map[string]string{"message": "Task updated successfully"}})
//}
//
//func (s *service) DeleteTask(ctx *fiber.Ctx) error {
//	id := ctx.Params("id")
//
//	//400
//	taskID, err := strconv.Atoi(id)
//	if err != nil {
//		s.log.Error("Invalid task ID", zap.Error(err))
//		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID")
//	}
//
//	existingTask, err := s.repo.GetTaskByID(ctx.Context(), taskID)
//	// 500
//	if err != nil {
//		s.log.Error("Failed to check if task exists", zap.Error(err))
//		return dto.InternalServerError(ctx)
//	}
//
//	// 404
//	if existingTask == nil {
//		s.log.Warn("Task not found", zap.Int("id", taskID))
//		return ctx.Status(fiber.StatusNotFound).JSON(dto.Response{
//			Status: "error",
//			Data:   map[string]string{"message": "Task not found"},
//		})
//	}
//
//	//500
//	err = s.repo.DeleteTask(ctx.Context(), taskID)
//	if err != nil {
//		s.log.Error("Failed to delete task", zap.Error(err))
//		return dto.InternalServerError(ctx)
//	}
//
//	//200
//	response := dto.Response{
//		Status: "success",
//		Data:   map[string]string{"message": "Task deleted successfully"},
//	}
//	return ctx.Status(fiber.StatusOK).JSON(response)
//}
