package dtos

import (
	"teagans-web-api/app/utilities/uuid"
)

type TaskCategoryDTO struct {
	BaseDTO
	UserID		uuid.UUID	`json:"userId"`

	Name		string		`json:"name"`
	Position	int64		`json:"position"`
}

type TaskCategoryListDTO struct {
	TaskCategories	[]TaskCategoryDTO	`json:"taskCategories"`
}
