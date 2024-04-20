package mappers

import (
	"teagans-web-api/app/utilities/enums"
	"github.com/microcosm-cc/bluemonday"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

func MapTaskInDTOToTask(dto dtos.TaskInDTO) models.Task {
	taskStatus, _ := enums.NewTaskStatus(dto.Status)
	taskPriority, _ := enums.NewTaskPriority(dto.Priority)

	return models.Task{
		TaskCategoryID: dto.TaskCategoryID,
		Title: dto.Title,
		DetailHtml: bluemonday.UGCPolicy().Sanitize(dto.DetailHtml),
		DetailJson: dto.DetailJson,
		Status: taskStatus,
		Priority: taskPriority,
		Position: dto.Position,
		Effort: dto.Effort,
		Cleared: dto.Cleared,
	}
}

func MapTaskToTaskOutDTO(task models.Task) dtos.TaskOutDTO {
	return dtos.TaskOutDTO{
		BaseDTO: dtos.BaseDTO{
			ID: task.ID,
		},
		TaskInDTO: dtos.TaskInDTO{
			TaskCategoryID: task.TaskCategoryID,
			Title: task.Title,
			DetailHtml: task.DetailHtml,
			DetailJson: task.DetailJson,
			Status: int64(task.Status),
			Priority: int64(task.Priority),
			Position: task.Position,
			Effort: task.Effort,
			Cleared: task.Cleared,
		},
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
