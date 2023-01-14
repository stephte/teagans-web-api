package mappers

import (
	"chi-users-project/app/services/dtos"
	"chi-users-project/app/models"
)

func MapCreateUserDTOToUser(dto dtos.CreateUserDTO) models.User {
	return models.User{
		FirstName: dto.FirstName,
		LastName: dto.LastName,
		Email: dto.Email,
		Role: dto.Role,
		Password: dto.Password,
	}
}

func MapUserDTOToUser(dto dtos.UserDTO) models.User {
	return models.User{
		BaseModel: models.BaseModel{
			Key: dto.Key,
		},
		FirstName: dto.FirstName,
		LastName: dto.LastName,
		Email: dto.Email,
		Role: dto.Role,
	}
}

func MapUserToUserDTO(user models.User) dtos.UserDTO {
	return dtos.UserDTO{
		BaseDTO: dtos.BaseDTO{
			Key: user.Key,
		},
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Role: user.Role,
	}
}

func MapUsersToUserDTOs(users []models.User) ([]dtos.UserDTO) {
	rv := []dtos.UserDTO{}

	for _, user := range users {
		rv = append(rv, MapUserToUserDTO(user))
	}

	return rv
}
