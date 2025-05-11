package grpc

import (
	"context"
	"fmt"

	taskpb "github.com/Retual004/project-protos/proto/task"
	userpb "github.com/Retual004/project-protos/proto/user"
	"github.com/Retual004/tasks-service/internal/task"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handler реализует gRPC-сервер для TaskService, внутри обращаясь к TaskService и UserServiceClient.
type Handler struct {
	svc        *task.TaskService
	userClient userpb.UserServiceClient
	taskpb.UnimplementedTaskServiceServer
}

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(svc *task.TaskService, uc userpb.UserServiceClient) *Handler {
	return &Handler{svc: svc, userClient: uc}
}

// CreateTask проверяет существование пользователя и создаёт новую задачу.
func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	// Проверяем, что пользователь существует
	if _, err := h.userClient.GetUser(ctx, &userpb.User{Id: req.UserId}); err != nil {
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	// Внутренняя модель использует поле Task для названия
	t, err := h.svc.CreateTask(task.Task{UserID: uint(req.UserId), Task: req.Title, IsDone: req.IsDone})
	if err != nil {
		return nil, err
	}

	// Возвращаем Protobuf-модель
	return &taskpb.CreateTaskResponse{Task: &taskpb.Task{
		Id:     uint32(t.ID),
		Title:  t.Task,
		IsDone: t.IsDone,
		UserId: uint32(t.UserID),
	}}, nil
}

// GetTask возвращает задачу по её ID.
func (h *Handler) GetTask(ctx context.Context, req *taskpb.Task) (*taskpb.Task, error) {
	// Извлекаем по ID
	t, err := h.svc.GetTaskByID(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &taskpb.Task{
		Id:     uint32(t.ID),
		Title:  t.Task,
		IsDone: t.IsDone,
		UserId: uint32(t.UserID),
	}, nil
}

// ListTasks возвращает все задачи.
func (h *Handler) ListTasks(ctx context.Context, _ *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	list, err := h.svc.GetAllTask()
	if err != nil {
		return nil, err
	}

	resp := &taskpb.ListTasksResponse{}
	for _, t := range list {
		resp.Tasks = append(resp.Tasks, &taskpb.Task{
			Id:     uint32(t.ID),
			Title:  t.Task,
			IsDone: t.IsDone,
			UserId: uint32(t.UserID),
		})
	}
	return resp, nil
}

// UpdateTask проверяет пользователя и обновляет задачу.
func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.Task, error) {
	// Проверка существования пользователя
	if _, err := h.userClient.GetUser(ctx, &userpb.User{Id: req.Task.UserId}); err != nil {
		return nil, fmt.Errorf("user %d not found: %w", req.Task.UserId, err)
	}

	// Обновляем через сервис
	tUpdated, err := h.svc.UpdateTaskByID(uint(req.Task.Id), task.Task{UserID: uint(req.Task.UserId), Task: req.Task.Title, IsDone: req.Task.IsDone})
	if err != nil {
		return nil, err
	}
	return &taskpb.Task{
		Id:     uint32(tUpdated.ID),
		Title:  tUpdated.Task,
		IsDone: tUpdated.IsDone,
		UserId: uint32(tUpdated.UserID),
	}, nil
}

// DeleteTask проверяет владельца и удаляет задачу.
func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*emptypb.Empty, error) {
	// Сначала получаем задачу, чтобы узнать UserID
	t, err := h.svc.GetTaskByID(uint(req.Id))
	if err != nil {
		return nil, err
	}

	// Проверяем, что пользователь все ещё существует
	if _, err := h.userClient.GetUser(ctx, &userpb.User{Id: uint32(t.UserID)}); err != nil {
		return nil, fmt.Errorf("user %d not found: %w", t.UserID, err)
	}

	// Удаляем задачу
	if err := h.svc.DeleteTaskByID(uint(req.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
