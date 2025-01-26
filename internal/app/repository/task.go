package repository

import (
	"errors"

	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"

	"gorm.io/gorm"
)

type taskRepository struct{}

var TaskRepository = new(taskRepository)

func (t *taskRepository) GetAllTask() (result []model.Task) {
	db.Find(&result)
	return
}

func (t *taskRepository) GetTaskById(id int) (result model.Task, err error) {
	err = db.Model(&model.Task{}).Where(&model.Task{Id: int(id)}).First(&result).Error
	return
}

func (t *taskRepository) GetTaskByKey(key string) (result model.Task, err error) {
	err = db.Model(&model.Task{}).Where(&model.Task{Key: &key, ApiEnable: true}).First(&result).Error
	return
}

func (t *taskRepository) AddTask(data model.Task) (taskId int, err error) {
	err = db.Create(&data).Error
	taskId = data.Id
	return
}

func (t *taskRepository) DeleteTask(id int) (err error) {
	err = db.Delete(&model.Task{Id: id}).Error
	return
}

func (t *taskRepository) EditTask(data model.Task) (err error) {
	err = db.Model(&model.Task{}).Where(&model.Task{Id: data.Id}).First(&model.Task{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = db.Model(&model.Task{}).Where(&model.Task{Id: data.Id}).Updates(data).Error
	return
}

func (t *taskRepository) EditTaskEnable(id int, enable bool) (err error) {
	err = db.Model(&model.Task{}).Where(&model.Task{Id: id}).First(&model.Task{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = db.Model(&model.Task{}).Where(&model.Task{Id: id}).Update("enable", enable).Error
	return
}

func (t *taskRepository) GetAllTaskWithProcessName() (result []model.TaskVo) {
	db.Raw(`SELECT t.*, p.name AS process_name, p2.name AS target_name,p3.name AS trigger_name 
	FROM task t LEFT JOIN process p ON t.process_id = p.uuid LEFT JOIN process p2 ON t.operation_target = p2.uuid LEFT JOIN process p3 ON t.trigger_target = p3.uuid`).Scan(&result)
	return
}

func (t *taskRepository) GetTriggerTask(processName string, event constants.ProcessState) []model.Task {
	result := []model.Task{}
	db.Raw(`SELECT task.* FROM task left join process p  on p.uuid == task.trigger_target  WHERE trigger_event= ? AND p.name = ?`, event, processName).Scan(&result)
	return result
}
