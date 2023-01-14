package services

import (
	"chi-users-project/app/services/mappers"
	"chi-users-project/app/services/emails"
	"chi-users-project/app/utilities/auth"
	"chi-users-project/app/services/dtos"
	"chi-users-project/app/models"
	"github.com/google/uuid"
	"errors"
)

type UserService struct {
	*BaseService
	user						models.User
}


func(this UserService) GetUser(userKeyStr string) (dtos.UserDTO, dtos.ErrorDTO) {
	err := this.setUserByKeyStr(userKeyStr)
	if err != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	if !this.validateUserHasAccess(auth.AdminAccess()) && this.currentUser.ID != this.user.ID {
		return dtos.UserDTO{}, dtos.AccessDeniedError()
	}

	return mappers.MapUserToUserDTO(this.user), dtos.ErrorDTO{}
}


func (this UserService) GetUsers(dto dtos.PaginationDTO, path string) (dtos.PageResponseDTO, dtos.ErrorDTO) {
	if !this.validateUserHasAccess(auth.AdminAccess()) {
		return dto.GetPageResponse(), dtos.AccessDeniedError()
	}

	// first get count of total rows
	var count int64
	countErr := this.db.Model(&models.User{}).Count(&count).Error
	if countErr != nil {
		return dto.GetPageResponse(), dtos.CreateErrorDTO(countErr, 500, false)
	}

	dto.SetTotalRows(count)
	dto.GenAndSetData(path) // can call this once TotalRows is set

	var users []models.User
	if err := this.db.Limit(dto.GetLimit()).Offset(dto.GetOffset()).Order(dto.GetSort()).Find(&users).Error; err != nil {
		return dto.GetPageResponse(), dtos.CreateErrorDTO(err, 500, false)
	}

	userDTOs := mappers.MapUsersToUserDTOs(users)
	dto.SetRows(userDTOs)

	return dto.GetPageResponse(), dtos.ErrorDTO{}
}


// takes in CreateUserDTO, returns UserDTO
func (this UserService) CreateUser(dto dtos.CreateUserDTO) (dtos.UserDTO, dtos.ErrorDTO) {
	if dto.Role > auth.RegularAccess()  {
		if !this.validateUserHasAccess(auth.SuperAdminAccess()) {
			this.log.Error().Msg("Invalid Role: User create with admin attempted by non super-admin")
			return dtos.UserDTO{}, dtos.CreateErrorDTO(errors.New("Invalid create params"), 0, false)
		}
	}

	user := mappers.MapCreateUserDTOToUser(dto)

	if createErr := this.db.Create(&user).Error; createErr != nil {
		this.log.Warn().Msg(createErr.Error())
		return dtos.UserDTO{}, dtos.CreateErrorDTO(createErr, 0, false)
	}

	go this.sendSignupEmail(user.Email, user.FirstName)

	rv := mappers.MapUserToUserDTO(user)

	return rv, dtos.ErrorDTO{}
}


// saved via a Map thats validated
func (this UserService) UpdateUser(userKeyStr string, data map[string]interface{}) (dtos.UserDTO, dtos.ErrorDTO) {
	// validate User update data
	validatedData, dataErr := dtos.ValidateUserMap(data)
	if dataErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(dataErr, 0, false)
	}

	err := this.setUserByKeyStr(userKeyStr)
	if err != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	// check if role exists in data; else resume as if its equal to users current role
	var role int
	_, exists := validatedData["Role"]
	if exists {
		roleFloat, isFloat := validatedData["Role"].(float64)
		if !isFloat {
			return dtos.UserDTO{}, dtos.CreateErrorDTO(errors.New("Role is not a float64"), 0, false)
		}

		role = int(roleFloat)
	} else {
		role = this.user.Role
	}

	if !this.validateUserHasAccess(auth.SuperAdminAccess()) && (role != this.user.Role || this.currentUser.ID != this.user.ID) {
		return dtos.UserDTO{}, dtos.AccessDeniedError()
	}

	if updateErr := this.db.Model(&this.user).Updates(validatedData).Error; updateErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(updateErr, 0, false)
	}

	return mappers.MapUserToUserDTO(this.user), dtos.ErrorDTO{}
}


// saves via UserDTO thats converted to a User model
func (this UserService) UpdateUserOG(userKeyStr string, dto dtos.UserDTO) (dtos.UserDTO, dtos.ErrorDTO) {
	err := this.setUserByKeyStr(userKeyStr)
	if err != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	// handle validation (only super admins can update Role)
	if !this.validateUserHasAccess(auth.SuperAdminAccess()) && (dto.Role != this.user.Role || this.currentUser.ID != this.user.ID) {
		return dto, dtos.AccessDeniedError()
	}

	updatedUser := mappers.MapUserDTOToUser(dto)

	// will have issue updating Role to 0 (GORM only updates non-zero fields when updating with struct)
	if updateErr := this.db.Model(&this.user).Updates(updatedUser).Error; updateErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(updateErr, 0, false)
	}

	return mappers.MapUserToUserDTO(this.user), dtos.ErrorDTO{}
}


func(this UserService) DeleteUser(userKeyStr string) dtos.ErrorDTO {
	err := this.setUserByKeyStr(userKeyStr)
	if err != nil {
		return dtos.CreateErrorDTO(err, 0, false)
	}

	if !this.validateUserHasAccess(auth.SuperAdminAccess()) && this.currentUser.ID != this.user.ID {
		return dtos.AccessDeniedError()
	}

	// Unscoped actually deletes the User, without it it just sets the 'DeletedAt' field
	if deleteErr := this.db.Unscoped().Delete(&this.user).Error; deleteErr != nil {
		return dtos.CreateErrorDTO(deleteErr, 0, false)
	}

	return dtos.ErrorDTO{}
}


// ---------- Private ---------


func(this *UserService) setUserByKeyStr(userKeyStr string) error {
	key, parseErr := uuid.Parse(userKeyStr)
	if parseErr != nil {
		return parseErr
	}

	user, findErr := this.findUserByKey(key)
	if findErr != nil {
		return findErr
	}

	this.user = user

	return nil
}

func(this UserService) sendSignupEmail(email, firstName string) error {
	this.log.Debug().Msg("Sending Signup confirmation email")

	request := emails.SignupEmail {
		BaseEmailRequest: emails.InitBaseRequest(),
		FirstName: firstName,
	}

	request.SetToEmails([]string{email})
	request.SetSubject("Teagans App Signup Confirmation")

	// generate html for email
	err := request.GenerateAndSetMessage()
	if err != nil {
		return err
	}

	// send email
	return request.SendEmail()
}
