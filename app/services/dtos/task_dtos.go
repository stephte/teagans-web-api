package dtos

import (
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/utilities/uuid"
)

type TaskDTO struct {
	TaskCategoryID		uuid.UUID			`json:"userId"`
	Name				string				`json:"name"`
	Details				string				`json:"details"`
	Priority			enums.TaskPriority	`json:"priority"`
	Effort				int64				`json:"effort"`
}

type TaskListDTO struct {
	Tasks	[]TaskDTO	`json:"tasks"`
}


type TaskCategoryDTO struct {
	BaseDTO
	UserID		uuid.UUID	`json:"userId"`
	Name		string		`json:"name"`
}

type TaskCategoryListDTO struct {
	TaskCategories	[]TaskCategoryDTO	`json:"taskCategories"`
}
