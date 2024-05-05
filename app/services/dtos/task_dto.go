package dtos

import (
	"teagans-web-api/app/utilities/uuid"
	"time"
)

type TaskInDTO struct {
	TaskCategoryID		uuid.UUID			`json:"taskCategoryID"`

	Title				string				`json:"title"`
	DetailHtml			string				`json:"detailHtml"`
	DetailJson			string				`json:"detailJson"`
	Status				int64				`json:"status" enum:"TaskStatus"`
	Priority			int64				`json:"priority" enum:"TaskPriority"`
	Position			int64				`json:"position"`
	DueDate				*time.Time			`json:"dueDate"`
	Cleared				bool				`json:"cleared"`
}

type TaskOutDTO struct {
	BaseDTO
	TaskInDTO

	TaskNumber			int64				`json:"taskNumber"`
}

type TaskListOutDTO struct {
	Tasks	[]TaskOutDTO	`json:"tasks"`
}

type TaskListInDTO struct {
	Tasks	[]map[string]interface{}	`json:"tasks"`
}
