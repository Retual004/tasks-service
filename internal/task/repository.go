package task

import (
	"fmt"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetAllTask() ([]Task, error)
	UpdateTaskByID(id uint, task Task) (Task, error)
	DeleteTaskByID(id uint) error
	GetTasksUserUserId(id uint) ([]Task, error)

	GetTaskByID(id uint) (Task, error)
}

type taskRepository struct{
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *taskRepository{
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task Task) (Task, error){
	if task.UserID == 0 {
		return Task{}, fmt.Errorf("user_id cannot be zero")
	}
	
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
} 

func (r *taskRepository) GetAllTask() ([]Task, error){
	var tasks []Task
err := r.db.Find(&tasks).Error
if err != nil {
	return []Task{}, err
}
return tasks, err
}

func (r *taskRepository) UpdateTaskByID(id uint, requestBody Task) (Task, error){
	var task Task
	result:= r.db.Find(&task, id)
	if result.Error != nil {
		return task , result.Error
	} 

	if requestBody.Task != "" {
		task.Task = requestBody.Task
	}
	if requestBody.IsDone != task.IsDone {
		task.IsDone = requestBody.IsDone
	}

	result = r.db.Save(&task)
	if result.Error != nil {
		return task , result.Error
	}

	return task , nil 
}

func (r *taskRepository) DeleteTaskByID(id uint) error{
	var task Task
result := r.db.First(&task, id)
if result.Error != nil{
	return result.Error
}
result = r.db.Delete(&task)
if result.Error != nil {
	return result.Error
}
return nil
}


func (r *taskRepository) GetTasksUserUserId(userID uint) ([]Task, error) {
	var tasks []Task
	// Предполагаем, что у задачи есть поле UserID, которое связывает задачу с пользователем
	if err := r.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) GetTaskByID(id uint) (Task, error) {
    var t Task
    if err := r.db.First(&t, id).Error; err != nil {
        return Task{}, err
    }
    return t, nil
}