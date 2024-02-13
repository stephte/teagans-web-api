package mappers

import (
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

func MapTaskDTOToTask(dto dtos.TaskDTO) models.Task {
	return models.Task{
		BaseModel: models.BaseModel{
			ID: dto.ID,
		},
		TaskCategoryID: dto.TaskCategoryID,
		Name: dto.Name,
		Details: dto.Details,
		Priority: dto.Priority,
		Effort: dto.Effort,
	}
}

func MapTaskToTaskDTO(task models.Task) dtos.TaskDTO {
	return dtos.TaskDTO{
		ID: task.ID,
		TaskCategoryID: task.TaskCategoryID,
		Name: task.Name,
		Details: task.Details,
		Priority: task.Priority,
		Effort: task.Effort,
	}
}

func MapTasksToTaskDTOs(tasks []models.Task) []dtos.TaskDTO {
	rv := []dtos.TaskDTO{}

	for _, task := range tasks {
		rv = append(rv, MapTaskToTaskDTO(task))
	}

	return rv
}


func MapTaskCategoryDTOToTaskCategory(dto dtos.TaskCategoryDTO) models.TaskCategory {
	return models.Task{
		BaseModel: models.BaseModel{
			ID: dto.ID,
		},
		UserID: dto.UserID,
		Name: dto.Name,
	}
}

func MapTaskCategoryToTaskCategoryDTO(tc models.TaskCategory) dtos.TaskCategoryDTO {
	return models.Task{
		ID: dto.ID,
		UserID: tc.UserID,
		Name: tc.Name,
	}
}
