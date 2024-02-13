package models

import (
	"teagans-web-api/app/utilities/uuid"
)

type TaskCategory struct {
	BaseModel

	UserID		uuid.UUID	`gorm:"type:uuid;"`

	User		User
	Tasks		[]Task

	Name		string
	Deleted		bool		`gorm:"default:false;"`
}
