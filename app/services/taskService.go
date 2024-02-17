package services

import (
	"teagans-web-api/app/services/mappers"
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/utilities/uuid"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

type TaskService struct {
	*BaseService
	task	models.Task
}

func(this TaskService) CreateTask(dto dtos.TaskDTO) (dtos.TaskDTO, dtos.ErrorDTO) {
	this.task = mappers.MapTaskDTOToTask(dto)
	this.task.ID = uuid.UUID{} // don't allow them to set the ID when creating

	// find task category its associating to
	taskCategory := this.getTaskCategory()

	// unsure if TaskCategory is loadable before saved
	if taskCategory.UserID != this.currentUser.ID && !this.validateUserHasAccess(enums.SUPERADMIN) {
		return dtos.TaskDTO{}, dtos.AccessDeniedError(false)
	}

	if createErr := this.db.Create(&this.task).Error; createErr != nil {
		this.log.Warn().Msg(createErr.Error())
		return dtos.TaskDTO{}, dtos.CreateErrorDTO(createErr, 0, false)
	}

	rv := mappers.MapTaskToTaskDTO(this.task)

	return rv, dtos.ErrorDTO{}
}

func(this TaskService) DeleteTask(taskIdStr string) dtos.ErrorDTO {
	err := this.setTask(taskIdStr)
	if err != nil {
		return dtos.CreateErrorDTO(err, 0, false)
	}

	if !this.validateUserHasAccess(enums.SUPERADMIN) && this.currentUser.ID != this.task.TaskCategory.UserID {
		return dtos.AccessDeniedError(false)
	}

	if deleteErr := this.db.Delete(&this.task).Error; deleteErr != nil {
		return dtos.CreateErrorDTO(deleteErr, 0, false)
	}

	return dtos.ErrorDTO{}
}

// ----- private -----

func(this *TaskService) setTask(categoryIdStr string) error {
	id, parseErr := uuid.Parse(categoryIdStr)
	if parseErr != nil {
		return parseErr
	}

	task, findErr := this.findTask(id)
	if findErr != nil {
		return findErr
	}

	this.task = task

	return nil
}

func(this TaskService) findTask(id uuid.UUID) (models.Task, error) {
	var rv models.Task
	if findErr := this.db.Preload("TaskCategory").First(&rv, id).Error; findErr != nil {
		return rv, findErr
	}

	return rv, nil
}

func(this TaskService) getTaskCategory() (models.TaskCategory) {
	var rv models.TaskCategory
	this.db.First(&rv, this.task.TaskCategoryID)
	return rv
}
