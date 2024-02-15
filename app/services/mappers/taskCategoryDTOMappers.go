package mappers

import (
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

func MapTaskCategoryDTOToTaskCategory(dto dtos.TaskCategoryDTO) models.TaskCategory {
	return models.TaskCategory{
		BaseModel: models.BaseModel{
			ID: dto.ID,
		},
		UserID: dto.UserID,
		Name: dto.Name,
		Position: dto.Position,
	}
}

func MapTaskCategoryToTaskCategoryDTO(tc models.TaskCategory) dtos.TaskCategoryDTO {
	return dtos.TaskCategoryDTO{
		BaseDTO: dtos.BaseDTO{
			ID: tc.ID,
		},
		UserID: tc.UserID,
		Name: tc.Name,
		Position: tc.Position,
	}
}

func MapTaskCategoriesToTaskCategoryDTOs(tcs []models.TaskCategory) []dtos.TaskCategoryDTO {
	rv := []dtos.TaskCategoryDTO{}

	for _, tc := range tcs {
		rv = append(rv, MapTaskCategoryToTaskCategoryDTO(tc))
	}

	return rv
}
