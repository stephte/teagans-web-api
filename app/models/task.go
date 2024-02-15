package models

import (
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/utilities/uuid"
	"gorm.io/gorm"
)

type Task struct {
	BaseModel

	TaskCategoryID		uuid.UUID			`gorm:"type:uuid;uniqueIndex:cattasknumndx;"`

	Title				string				`gorm:"default:null;not null;"`
	Details				string
	Status				enums.TaskStatus	`gorm:"default:null;not null;"`
	Priority			enums.TaskPriority
	Effort				int64
	Cleared				bool				`gorm:"default:false;"`
	TaskNumber			int64				`gorm:"uniqueIndex:cattasknumndx;"`
}

func(this *Task) BeforeCreate(tx *gorm.DB) error {
	var count int64
	tx.Where(&Task{TaskCategoryID: this.TaskCategoryID}).Count(&count)
	// .Where(&Task{}"task_category_id = ?", this.TaskCategoryID.String())

	this.TaskNumber = count + 1

	return nil
}
