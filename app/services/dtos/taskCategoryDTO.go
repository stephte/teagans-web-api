package dtos

import (
	"teagans-web-api/app/utilities/uuid"
)

type TaskCategoryInDTO struct {
	UserID		uuid.UUID	`json:"userId"`

	Name		string		`json:"name"`
	Position	int64		`json:"position"`
}

type TaskCategoryOutDTO struct {
	BaseDTO
	TaskCategoryInDTO
}

type TaskCategoryListDTO struct {
	TaskCategories	[]TaskCategoryOutDTO	`json:"taskCategories"`
}
