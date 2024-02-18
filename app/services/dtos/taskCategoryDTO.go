package dtos

import (
	"teagans-web-api/app/utilities/uuid"
)

type TaskCategoryOutDTO struct {
	BaseDTO
	UserID		uuid.UUID	`json:"userId"`

	Name		string		`json:"name"`
	Position	int64		`json:"position"`
}

type TaskCategoryInDTO struct {
	UserID		uuid.UUID	`json:"userId"`

	Name		string		`json:"name"`
	Position	int64		`json:"position"`
}

type TaskCategoryListDTO struct {
	TaskCategories	[]TaskCategoryOutDTO	`json:"taskCategories"`
}
