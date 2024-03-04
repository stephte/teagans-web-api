package services

import (
	intrfaceUtils "teagans-web-api/app/utilities/intrface"
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

func(this TaskService) CreateTask(dto dtos.TaskInDTO) (dtos.TaskOutDTO, dtos.ErrorDTO) {
	this.task = mappers.MapTaskInDTOToTask(dto)

	// find task category its associating to
	taskCategory := this.getTaskCategory()

	// unsure if TaskCategory is loadable before saved
	if taskCategory.UserID != this.currentUser.ID && !this.validateUserHasAccess(enums.SUPERADMIN) {
		return dtos.TaskOutDTO{}, dtos.AccessDeniedError(false)
	}

	if createErr := this.db.Create(&this.task).Error; createErr != nil {
		this.log.Warn().Msg(createErr.Error())
		return dtos.TaskOutDTO{}, dtos.CreateErrorDTO(createErr, 0, false)
	}

	rv := mappers.MapTaskToTaskOutDTO(this.task)

	return rv, dtos.ErrorDTO{}
}

func(this TaskService) GetTask(taskIdStr string) (dtos.TaskOutDTO, dtos.ErrorDTO) {
	findErr := this.setTask(taskIdStr)
	if findErr != nil {
		return dtos.TaskOutDTO{}, dtos.CreateErrorDTO(findErr, 0, false)
	}

	if !this.validateUserHasAccess(enums.SUPERADMIN) && this.currentUser.ID != this.task.TaskCategory.UserID {
		return dtos.TaskOutDTO{}, dtos.AccessDeniedError(false)
	}

	return mappers.MapTaskToTaskOutDTO(this.task), dtos.ErrorDTO{}
}

func(this TaskService) UpdateTask(data map[string]interface{}, taskIdStr string) (dtos.TaskOutDTO, dtos.ErrorDTO) {
	err := this.setTask(taskIdStr)
	if err != nil {
		return dtos.TaskOutDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	if !this.validateUserHasAccess(enums.SUPERADMIN) && this.currentUser.ID != this.task.TaskCategory.UserID {
		return dtos.TaskOutDTO{}, dtos.AccessDeniedError(false)
	}

	// convert dto to a map
	taskMap, mapErr := intrfaceUtils.ValidateMapWithStruct(data, dtos.TaskInDTO{})
	if mapErr != nil {
		return dtos.TaskOutDTO{}, dtos.CreateErrorDTO(mapErr, 0, false)
	}

	if updateErr := this.db.Model(&this.task).Updates(taskMap).Error; updateErr != nil {
		return dtos.TaskOutDTO{}, dtos.CreateErrorDTO(updateErr, 0, false)
	}

	rv := mappers.MapTaskToTaskOutDTO(this.task)

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

	if deleteErr := this.db.Unscoped().Delete(&this.task).Error; deleteErr != nil {
		return dtos.CreateErrorDTO(deleteErr, 0, false)
	}

	return dtos.ErrorDTO{}
}

// ----- private -----

func(this *TaskService) setTask(taskIdStr string) error {
	id, parseErr := uuid.Parse(taskIdStr)
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
