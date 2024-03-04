package services

import (
	intrfaceUtils "teagans-web-api/app/utilities/intrface"
	"teagans-web-api/app/services/mappers"
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/utilities/uuid"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
	"strings"
)

type TaskCategoryService struct {
	*BaseService
	taskCategory	models.TaskCategory
}

func(this TaskCategoryService) CreateTaskCategory(dto dtos.TaskCategoryInDTO) (dtos.TaskCategoryOutDTO, dtos.ErrorDTO) {
	this.taskCategory = mappers.MapTaskCategoryInDTOToTaskCategory(dto)
	this.taskCategory.ID = uuid.UUID{} // dont let ID be set by user

	// default it to current user if not present
	if !this.taskCategory.UserID.Exists() {
		this.taskCategory.UserID = this.currentUser.ID
	}

	if this.taskCategory.UserID != this.currentUser.ID && !this.validateUserHasAccess(enums.SUPERADMIN) {
		return dtos.TaskCategoryOutDTO{}, dtos.AccessDeniedError(false)
	}

	if createErr := this.db.Create(&this.taskCategory).Error; createErr != nil {
		this.log.Warn().Msg(createErr.Error())
		return dtos.TaskCategoryOutDTO{}, dtos.CreateErrorDTO(createErr, 0, false)
	}
	
	rv := mappers.MapTaskCategoryToTaskCategoryOutDTO(this.taskCategory)

	return rv, dtos.ErrorDTO{}
}

func(this TaskCategoryService) UpdateTaskCategory(data map[string]interface{}, taskCategoryIdStr string) (dtos.TaskCategoryOutDTO, dtos.ErrorDTO) {
	err := this.setTaskCategory(taskCategoryIdStr)
	if err != nil {
		return dtos.TaskCategoryOutDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	if !this.validateUserHasAccess(enums.SUPERADMIN) && this.currentUser.ID != this.taskCategory.UserID {
		return dtos.TaskCategoryOutDTO{}, dtos.AccessDeniedError(false)
	}

	tcMap, mapErr := intrfaceUtils.ValidateMapWithStruct(data, dtos.TaskCategoryInDTO{})
	if mapErr != nil {
		return dtos.TaskCategoryOutDTO{}, dtos.CreateErrorDTO(mapErr, 0, false)
	}

	// update task category
	if updateErr := this.db.Model(&this.taskCategory).Omit("user_id").Updates(tcMap).Error; updateErr != nil {
		return dtos.TaskCategoryOutDTO{}, dtos.CreateErrorDTO(updateErr, 0, false)
	}

	rv := mappers.MapTaskCategoryToTaskCategoryOutDTO(this.taskCategory)

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

func(this TaskCategoryService) GetTaskCategoryTasks(categoryIdStr, statusQuery string, getCleared bool) (dtos.TaskListDTO, dtos.ErrorDTO) {
	err := this.setTaskCategory(categoryIdStr)
	if err != nil {
		return dtos.TaskListDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	if !this.validateUserHasAccess(enums.SUPERADMIN) && this.currentUser.ID != this.taskCategory.UserID {
		return dtos.TaskListDTO{}, dtos.AccessDeniedError(false)
	}

	statusList := genStatusList(statusQuery)

	var tasks []models.Task
	this.db.Where("task_category_id = ? AND cleared = ? AND status IN ?", this.taskCategory.ID, getCleared, statusList).Order("priority desc").Find(&tasks)

	rv := dtos.TaskListDTO{
		Tasks: mappers.MapTasksToTaskOutDTOs(tasks),
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

func genStatusList(statusQuery string) []enums.TaskStatus {
	var statusList []enums.TaskStatus
	if statusQuery == "" || strings.ToLower(statusQuery) == "all" {
		// create arr of all statuses
		max := enums.TaskStatusCount()
		for i := 1; i <= max; i++ {
			statusList = append(statusList, enums.TaskStatus(i))
		}
	} else {
		// process status query string:
		statuses := strings.Split(statusQuery, ",")
		for _, status := range statuses {
			statusEnum, exists := enums.ParseTaskStatusString(status)
			if exists {
				statusList = append(statusList, statusEnum)
			}
		}
	}

	return statusList
}
