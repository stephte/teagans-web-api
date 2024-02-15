package models

import (
	"teagans-web-api/app/utilities/uuid"
)

type TaskCategory struct {
	BaseModel

	UserID		uuid.UUID	`gorm:"type:uuid;"`

	User		User
	Tasks		[]Task		`gorm:"constraint:OnDelete:CASCADE;"`

	Name		string		`gorm:"default:null;not null;"`
	Position	int64		`gorm:"default:1;"`
}
