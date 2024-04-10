package dtos

import (
	"teagans-web-api/app/utilities/uuid"
)

type TaskInDTO struct {
	TaskCategoryID		uuid.UUID			`json:"taskCategoryID"`

	Title				string				`json:"title"`
	DetailHtml			string				`json:"detailHtml"`
	DetailJson			string				`json:"detailJson"`
	Status				int64				`json:"status" enum:"TaskStatus"`
	Priority			int64				`json:"priority" enum:"TaskPriority"`
	Position			int64				`json:"position"`
	Effort				int64				`json:"effort"`
	Cleared				bool				`json:"cleared"`
}

type TaskOutDTO struct {
	BaseDTO
	TaskInDTO

	TaskNumber			int64				`json:"taskNumber"`
}

type TaskListDTO struct {
	Tasks	[]TaskOutDTO	`json:"tasks"`
}
