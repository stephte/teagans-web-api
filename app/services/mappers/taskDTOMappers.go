package mappers

import (
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

func MapTaskDTOToTask(dto dtos.TaskDTO) models.Task {
	taskStatus, _ := enums.ParseTaskStatusString(dto.Status)
	taskPriority, _ := enums.ParseTaskPriorityString(dto.Priority)

	return models.Task{
		BaseModel: models.BaseModel{
			ID: dto.ID,
		},
		TaskCategoryID: dto.TaskCategoryID,
		Title: dto.Title,
		Details: dto.Details,
		Status: taskStatus,
		Priority: taskPriority,
		Effort: dto.Effort,
		Cleared: dto.Cleared,
		TaskNumber: dto.TaskNumber,
	}
}

func MapTaskToTaskDTO(task models.Task) dtos.TaskDTO {
	return dtos.TaskDTO{
		BaseDTO: dtos.BaseDTO{
			ID: task.ID,
		},
		TaskCategoryID: task.TaskCategoryID,
		Title: task.Title,
		Details: task.Details,
		Status: task.Status.String(),
		Priority: task.Priority.String(),
		Effort: task.Effort,
		Cleared: task.Cleared,
		TaskNumber: task.TaskNumber,
	}
}

func MapTasksToTaskDTOs(tasks []models.Task) []dtos.TaskDTO {
	rv := []dtos.TaskDTO{}

	for _, task := range tasks {
		rv = append(rv, MapTaskToTaskDTO(task))
	}

	return rv
}
