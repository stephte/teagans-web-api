package dtos

import (
	"teagans-web-api/app/utilities/uuid"
)

type TaskDTO struct {
	BaseDTO
	TaskCategoryID		uuid.UUID			`json:"userId"`

	Title				string				`json:"title"`
	Details				string				`json:"details"`
	Status				string				`json:"status"`
	Priority			string				`json:"priority"`
	Effort				int64				`json:"effort"`
	Cleared				bool				`json:"cleared"`
	// TaskNumber			int64				`json:"taskNumber"`
}

type TaskListDTO struct {
	Tasks	[]TaskDTO	`json:"tasks"`
}
