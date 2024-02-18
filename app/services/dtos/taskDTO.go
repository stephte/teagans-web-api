package dtos

import (
	"teagans-web-api/app/utilities/uuid"
)

type TaskOutDTO struct {
	BaseDTO
	TaskCategoryID		uuid.UUID			`json:"taskCategoryID"`

	Title				string				`json:"title"`
	Details				string				`json:"details"`
	Status				string				`json:"status"`
	Priority			string				`json:"priority"`
	Effort				int64				`json:"effort"`
	Cleared				bool				`json:"cleared"`
	TaskNumber			int64				`json:"taskNumber"`
}

type TaskInDTO struct {
	TaskCategoryID		uuid.UUID			`json:"taskCategoryID"`

	Title				string				`json:"title"`
	Details				string				`json:"details"`
	Status				string				`json:"status"`
	Priority			string				`json:"priority"`
	Effort				int64				`json:"effort"`
	Cleared				bool				`json:"cleared"`
}

type TaskListDTO struct {
	Tasks	[]TaskOutDTO	`json:"tasks"`
}
