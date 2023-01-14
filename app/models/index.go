package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	Key						uuid.UUID	`gorm:"type:uuid;uniqueIndex;not null;default:uuid_generate_v4()"`
}
