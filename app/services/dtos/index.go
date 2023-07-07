package dtos

import (
	"chi-users-project/app/utilities/uuid"
)

type BaseDTO struct {
	ID		uuid.UUID	`json:"id"`
}
