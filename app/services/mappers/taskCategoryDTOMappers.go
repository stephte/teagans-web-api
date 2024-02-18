package mappers

import (
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

func MapTaskCategoryInDTOToTaskCategory(dto dtos.TaskCategoryInDTO) models.TaskCategory {
	return models.TaskCategory{
		UserID: dto.UserID,
		Name: dto.Name,
		Position: dto.Position,
	}
}

func MapTaskCategoryToTaskCategoryOutDTO(tc models.TaskCategory) dtos.TaskCategoryOutDTO {
	return dtos.TaskCategoryOutDTO{
		BaseDTO: dtos.BaseDTO{
			ID: tc.ID,
		},
		UserID: tc.UserID,
		Name: tc.Name,
		Position: tc.Position,
	}
}

func MapTaskCategoriesToTaskCategoryOutDTOs(tcs []models.TaskCategory) []dtos.TaskCategoryOutDTO {
	rv := []dtos.TaskCategoryOutDTO{}

	for _, tc := range tcs {
		rv = append(rv, MapTaskCategoryToTaskCategoryOutDTO(tc))
	}

	return rv
}
