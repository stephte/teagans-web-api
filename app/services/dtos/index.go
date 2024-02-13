package dtos

import (
	"teagans-web-api/app/utilities/uuid"
)

type BaseDTO struct {
	ID		uuid.UUID	`json:"id"`
}
