package mappers

import (
	"teagans-web-api/app/utilities/enums"
	"github.com/microcosm-cc/bluemonday"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
	"time"
)

func MapTaskInDTOToTask(dto dtos.TaskInDTO) models.Task {
	taskStatus, _ := enums.NewTaskStatus(dto.Status)
	taskPriority, _ := enums.NewTaskPriority(dto.Priority)
	format := "2006-01-02"
	dd, _ := time.Parse(format, dto.DueDate)

	return models.Task{
		TaskCategoryID: dto.TaskCategoryID,
		Title: dto.Title,
		DetailHtml: bluemonday.UGCPolicy().Sanitize(dto.DetailHtml),
		DetailJson: dto.DetailJson,
		Status: taskStatus,
		Priority: taskPriority,
		Position: dto.Position,
		DueDate: &dd,
		Cleared: dto.Cleared,
	}
}

func MapTaskToTaskOutDTO(task models.Task) dtos.TaskOutDTO {
	var dd string
	if task.DueDate == nil {
		dd = ""
	} else {
		format := "2006-01-02"
		dd = task.DueDate.Format(format)
	}
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
			DueDate: dd,
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
