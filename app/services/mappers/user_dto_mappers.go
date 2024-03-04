package mappers

import (
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
)

func MapCreateUserDTOToUser(dto dtos.CreateUserDTO) models.User {
	userRole, _ := enums.NewUserRole(dto.Role)
	return models.User{
		FirstName: dto.FirstName,
		LastName: dto.LastName,
		Email: dto.Email,
		Role: userRole,
		Password: dto.Password,
	}
}

func MapUserInDTOToUser(dto dtos.UserInDTO) models.User {
	userRole, _ := enums.NewUserRole(dto.Role)
	return models.User{
		FirstName: dto.FirstName,
		LastName: dto.LastName,
		Email: dto.Email,
		Role: userRole,
	}
}

func MapUserToUserOutDTO(user models.User) dtos.UserOutDTO {
	return dtos.UserOutDTO{
		BaseDTO: dtos.BaseDTO{
			ID: user.ID,
		},
		UserInDTO: dtos.UserInDTO{
			FirstName: user.FirstName,
			LastName: user.LastName,
			Email: user.Email,
			Role: int64(user.Role),
		},
	}
}

func MapUsersToUserOutDTOs(users []models.User) ([]dtos.UserOutDTO) {
	rv := []dtos.UserOutDTO{}

	for _, user := range users {
		rv = append(rv, MapUserToUserOutDTO(user))
	}

	return rv
}
