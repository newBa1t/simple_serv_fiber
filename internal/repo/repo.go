package repo

//
//import (
//	"context"
//	"fmt"
//	"github.com/jackc/pgx/v5"
//
//	"github.com/jackc/pgx/v5/pgxpool"
//	"github.com/pkg/errors"
//
//	"simple-service/internal/config"
//)
//
//// Слой репозитория, здесь должны быть все методы, связанные с базой данных
//
//// SQL-запрос на вставку задачи
//const (
//	insertTaskQuery     = `INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id;`
//	selectTaskByIDQuery = `SELECT id, title, description FROM tasks WHERE id = $1;`      // Новый запрос
//	selectAllTasksQuery = `SELECT id, title, description FROM tasks;`                    // Новый запрос для получения всех задач
//	updateTaskQuery     = `UPDATE tasks SET title = $1, description = $2 WHERE id = $3;` // Новый запрос для обновления задачи
//	deleteTaskQuery     = `DELETE FROM tasks WHERE id = $1;`                             // Новый запрос для удаления задачи
//)
//
//type repository struct {
//	pool *pgxpool.Pool
//}
//
//// Repository - интерфейс с методом создания задачи
//type Repository interface {
//	CreateTask(ctx context.Context, task Task) (int, error) // Создание задачи
//	GetTaskByID(ctx context.Context, id int) (*Task, error) // Получение по id
//	GetAllTasks(ctx context.Context) ([]*Task, error)       // Получение всех задач
//	UpdateTask(ctx context.Context, task Task) error        // Обновление
//	DeleteTask(ctx context.Context, id int) error           // Удаление
//}
//
//// NewRepository - создание нового экземпляра репозитория с подключением к PostgreSQL
//func NewRepository(ctx context.Context, cfg config.PostgreSQL) (Repository, error) {
//	// Формируем строку подключения
//	connString := fmt.Sprintf(
//		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s
//        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
//		cfg.User,
//		cfg.Password,
//		cfg.Host,
//		cfg.Port,
//		cfg.Name,
//		cfg.SSLMode,
//		cfg.PoolMaxConns,
//		cfg.PoolMaxConnLifetime.String(),
//		cfg.PoolMaxConnIdleTime.String(),
//	)
//
//	// Парсим конфигурацию подключения
//	config, err := pgxpool.ParseConfig(connString)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
//	}
//
//	// Оптимизация выполнения запросов (кеширование запросов)
//	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe
//
//	// Создаём пул соединений с базой данных
//	pool, err := pgxpool.NewWithConfig(ctx, config)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
//	}
//
//	return &repository{pool}, nil
//}
//
//// CreateTask - вставка новой задачи в таблицу tasks
//func (r *repository) CreateTask(ctx context.Context, task Task) (int, error) {
//	var id int
//	err := r.pool.QueryRow(ctx, insertTaskQuery, task.Title, task.Description).Scan(&id)
//	if err != nil {
//		return 0, errors.Wrap(err, "failed to insert task")
//	}
//	return id, nil
//}
//
//// GetTaskByID - получение задачи по ID
//func (r *repository) GetTaskByID(ctx context.Context, id int) (*Task, error) {
//	var task Task
//	err := r.pool.QueryRow(ctx, selectTaskByIDQuery, id).Scan(&task.ID, &task.Title, &task.Description)
//	if err != nil {
//		if err == pgx.ErrNoRows {
//			return nil, nil // Задача не найдена
//		}
//		return nil, errors.Wrap(err, "failed to get task by ID")
//	}
//	return &task, nil
//}
//
//func (r *repository) GetAllTasks(ctx context.Context) ([]*Task, error) {
//	rows, err := r.pool.Query(ctx, selectAllTasksQuery)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to get all tasks")
//	}
//	defer rows.Close()
//
//	var tasks []*Task
//	for rows.Next() {
//		var task Task
//		if err := rows.Scan(&task.ID, &task.Title, &task.Description); err != nil {
//			return nil, errors.Wrap(err, "failed to get all tasks")
//		}
//		tasks = append(tasks, &task)
//	}
//	return tasks, nil
//}
//
//func (r *repository) UpdateTask(ctx context.Context, task Task) error {
//	_, err := r.pool.Exec(ctx, updateTaskQuery, task.Title, task.Description, task.ID)
//	if err != nil {
//		return errors.Wrap(err, "failed to update task")
//	}
//	return nil
//}
//
//func (r *repository) DeleteTask(ctx context.Context, id int) error {
//	_, err := r.pool.Exec(ctx, deleteTaskQuery, id)
//	if err != nil {
//		return errors.Wrap(err, "failed to delete task")
//	}
//	return nil
//}
