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

type TaskCategoryListOutDTO struct {
	TaskCategories	[]TaskCategoryOutDTO	`json:"taskCategories"`
}

type TaskCategoryListInDTO struct {
	TaskCategories	[]map[string]interface{}	`json:"taskCategories"`
}
