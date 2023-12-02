package dtos

import (
	"youtube-downloader/app/utilities/uuid"
)

type BaseDTO struct {
	ID		uuid.UUID	`json:"id"`
}
