package models

import (
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/utilities/uuid"
)

type Task struct {
	BaseModel

	TaskCategoryID		uuid.UUID			`gorm:"type:uuid;"`

	Name				string
	Details				string
	Status				enums.TaskStatus
	Priority			enums.TaskPriority
	Effort				int64
	Cleared				bool				`gorm:"default:false;"`
	TaskNumber			int64				
}
