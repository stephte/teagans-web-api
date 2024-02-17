package services

import (
	"teagans-web-api/app/services/mappers"
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/utilities/uuid"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

type TaskCategoryService struct {
	*BaseService
	taskCategory	models.TaskCategory
}

func(this TaskCategoryService) CreateTaskCategory(dto dtos.TaskCategoryDTO) (dtos.TaskCategoryDTO, dtos.ErrorDTO) {
	this.taskCategory = mappers.MapTaskCategoryDTOToTaskCategory(dto)
	this.taskCategory.ID = uuid.UUID{} // dont let ID be set by user

	// default it to current user if not present
	if !this.taskCategory.UserID.Exists() {
		this.taskCategory.UserID = this.currentUser.ID
	}

	if this.taskCategory.UserID != this.currentUser.ID && !this.validateUserHasAccess(enums.SUPERADMIN) {
		return dtos.TaskCategoryDTO{}, dtos.AccessDeniedError(false)
	}

	if createErr := this.db.Create(&this.taskCategory).Error; createErr != nil {
		this.log.Warn().Msg(createErr.Error())
		return dtos.TaskCategoryDTO{}, dtos.CreateErrorDTO(createErr, 0, false)
	}
	
	rv := mappers.MapTaskCategoryToTaskCategoryDTO(this.taskCategory)

	return rv, dtos.ErrorDTO{}
}

func(this TaskCategoryService) DeleteTaskCategory(taskCategoryIdStr string) dtos.ErrorDTO {
	err := this.setTaskCategory(taskCategoryIdStr)
	if err != nil {
		return dtos.CreateErrorDTO(err, 0, false)
	}

	if !this.validateUserHasAccess(enums.SUPERADMIN) && this.currentUser.ID != this.taskCategory.UserID {
		return dtos.AccessDeniedError(false)
	}

	if deleteErr := this.db.Delete(&this.taskCategory).Error; deleteErr != nil {
		return dtos.CreateErrorDTO(deleteErr, 0, false)
	}

	return dtos.ErrorDTO{}
}

func(this TaskCategoryService) GetTaskCategoryTasks(categoryIdStr string) (dtos.TaskListDTO, dtos.ErrorDTO) {
	err := this.setTaskCategory(categoryIdStr)
	if err != nil {
		return dtos.TaskListDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	if !this.validateUserHasAccess(enums.SUPERADMIN) && this.currentUser.ID != this.taskCategory.UserID {
		return dtos.TaskListDTO{}, dtos.AccessDeniedError(false)
	}

	var tasks []models.Task
	this.db.Model(&this.taskCategory).Order("priority desc").Association("Tasks").Find(&tasks)

	rv := dtos.TaskListDTO{
		Tasks: mappers.MapTasksToTaskDTOs(tasks),
	}

	return rv, dtos.ErrorDTO{}
}

// ----- private -----

func(this *TaskCategoryService) setTaskCategory(categoryIdStr string) error {
	id, parseErr := uuid.Parse(categoryIdStr)
	if parseErr != nil {
		return parseErr
	}

	taskCategory, findErr := this.findTaskCategory(id)
	if findErr != nil {
		return findErr
	}

	this.taskCategory = taskCategory

	return nil
}

func (this TaskCategoryService) findTaskCategory(id uuid.UUID) (models.TaskCategory, error) {
	var rv models.TaskCategory
	if findErr := this.db.First(&rv, id).Error; findErr != nil {
		return rv, findErr
	}

	return rv, nil
}
