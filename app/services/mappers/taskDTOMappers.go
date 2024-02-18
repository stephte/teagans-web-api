package mappers

import (
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

func MapTaskInDTOToTask(dto dtos.TaskInDTO) models.Task {
	taskStatus, _ := enums.ParseTaskStatusString(dto.Status)
	taskPriority, _ := enums.ParseTaskPriorityString(dto.Priority)

	return models.Task{
		TaskCategoryID: dto.TaskCategoryID,
		Title: dto.Title,
		Details: dto.Details,
		Status: taskStatus,
		Priority: taskPriority,
		Effort: dto.Effort,
		Cleared: dto.Cleared,
	}
}

func MapTaskToTaskOutDTO(task models.Task) dtos.TaskOutDTO {
	return dtos.TaskOutDTO{
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

func MapTasksToTaskOutDTOs(tasks []models.Task) []dtos.TaskOutDTO {
	rv := []dtos.TaskOutDTO{}

	for _, task := range tasks {
		rv = append(rv, MapTaskToTaskOutDTO(task))
	}

	return rv
}
