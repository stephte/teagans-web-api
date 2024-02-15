package services

import (
	"teagans-web-api/app/services/mappers"
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

type TaskCategoryService struct {
	*BaseService
	taskCategory	models.TaskCategory
}

func(this TaskCategoryService) CreateTaskCategory(dto dtos.TaskCategoryDTO) (dtos.TaskCategoryDTO, dtos.ErrorDTO) {
	this.taskCategory = mappers.MapTaskCategoryDTOToTaskCategory(dto)

	// check that this can't happen, put this in case it can
	// this.taskCategory.ID = uuid.UUID{}

	// unsure if need this
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
