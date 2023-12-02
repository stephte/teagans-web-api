package models

import (
	"youtube-downloader/app/utilities/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID			uuid.UUID		`gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	DeletedAt	gorm.DeletedAt	`gorm:"index"`
}
