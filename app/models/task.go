package models

import (
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/utilities/uuid"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	BaseModel

	TaskCategoryID		uuid.UUID			`gorm:"type:uuid;uniqueIndex:cattasknumndx;not null;"`
	TaskCategory		TaskCategory

	Title				string				`gorm:"default:null;not null;"`

	DetailHtml			string
	DetailJson			string

	Status				enums.TaskStatus
	Priority			enums.TaskPriority
	Position			int64
	DueDate				*time.Time			`gorm:"type:timestamp;"`
	Cleared				bool				`gorm:"default:false;"`
	TaskNumber			int64				`gorm:"uniqueIndex:cattasknumndx;"`
}

func(this *Task) BeforeCreate(tx *gorm.DB) error {
	var lastTask Task
	tx.Model(&Task{}).Unscoped().Order("task_number desc").Where("task_category_id = ?", this.TaskCategoryID).First(&lastTask)

	this.TaskNumber = lastTask.TaskNumber + 1

	return nil
}
