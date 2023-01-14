package dtos

import (
	"github.com/google/uuid"
)

type BaseDTO struct {
	Key				uuid.UUID	`json:"key"`
}
