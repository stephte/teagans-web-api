package services

import (
	"teagans-web-api/app/services/mappers"
	"teagans-web-api/app/services/emails"
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/models"
	"github.com/google/uuid"
	"errors"
)

type UserService struct {
	*BaseService
	user			models.User
}


func(this UserService) GetUser(userIdStr string) (dtos.UserDTO, dtos.ErrorDTO) {
	// if they ask for current, get current user
	if userIdStr == "current" {
		this.user = this.currentUser
	} else {
		err := this.setUser(userIdStr)
		if err != nil {
			return dtos.UserDTO{}, dtos.CreateErrorDTO(err, 0, false)
		}
	}

	if !this.validateUserHasAccess(enums.ADMIN) && this.currentUser.ID != this.user.ID {
		return dtos.UserDTO{}, dtos.AccessDeniedError(false)
	}

	return mappers.MapUserToUserDTO(this.user), dtos.ErrorDTO{}
}


func (this UserService) GetUsers(dto dtos.PaginationDTO, path string) (dtos.PageResponseDTO, dtos.ErrorDTO) {
	if !this.validateUserHasAccess(enums.ADMIN) {
		return dto.GetPageResponse(), dtos.AccessDeniedError(false)
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
	if dto.Role > enums.REGULAR  {
		if !this.validateUserHasAccess(enums.SUPERADMIN) {
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
func (this UserService) UpdateUser(userIdStr string, data map[string]interface{}) (dtos.UserDTO, dtos.ErrorDTO) {
	// validate User update data
	validatedData, dataErr := dtos.ValidateUserMap(data)
	if dataErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(dataErr, 0, false)
	}

	err := this.setUser(userIdStr)
	if err != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	// check if role exists in data; else resume as if its equal to users current role
	var role enums.UserRole
	_, exists := validatedData["Role"]
	if exists {
		var success bool
		role, success = enums.ValToRole(validatedData["Role"])
		if !success {
			return dtos.UserDTO{}, dtos.CreateErrorDTO(errors.New("Not able to convert role to UserRole type"), 0, false)
		}
	} else {
		role = this.user.Role
	}

	if !((this.validateUserHasAccess(enums.ADMIN) && this.currentUser.Role >= role && (this.validateUserHasAccess(enums.SUPERADMIN) || this.user.Role < enums.ADMIN)) || (role == this.user.Role && this.currentUser.ID == this.user.ID)) {
		return dtos.UserDTO{}, dtos.AccessDeniedError(false)
	}

	if updateErr := this.db.Model(&this.user).Updates(validatedData).Error; updateErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(updateErr, 0, false)
	}

	return mappers.MapUserToUserDTO(this.user), dtos.ErrorDTO{}
}


// saves via UserDTO thats converted to a User model
func (this UserService) UpdateUserOG(userIdStr string, dto dtos.UserDTO) (dtos.UserDTO, dtos.ErrorDTO) {
	err := this.setUser(userIdStr)
	if err != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(err, 0, false)
	}

	// handle validation (only super admins can update Role)
	if !this.validateUserHasAccess(enums.SUPERADMIN) && (dto.Role != this.user.Role || this.currentUser.ID != this.user.ID) {
		return dto, dtos.AccessDeniedError(false)
	}

	updatedUser := mappers.MapUserDTOToUser(dto)

	// will have issue updating Role to 0 (GORM only updates non-zero fields when updating with struct)
	if updateErr := this.db.Model(&this.user).Updates(updatedUser).Error; updateErr != nil {
		return dtos.UserDTO{}, dtos.CreateErrorDTO(updateErr, 0, false)
	}

	return mappers.MapUserToUserDTO(this.user), dtos.ErrorDTO{}
}


func(this UserService) DeleteUser(userIdStr string) dtos.ErrorDTO {
	err := this.setUser(userIdStr)
	if err != nil {
		return dtos.CreateErrorDTO(err, 0, false)
	}

	if !this.validateUserHasAccess(enums.SUPERADMIN) && !(this.validateUserHasAccess(enums.ADMIN) && this.user.Role == enums.REGULAR) && this.currentUser.ID != this.user.ID {
		return dtos.AccessDeniedError(false)
	}

	// Unscoped actually deletes the User, without it it just sets the 'DeletedAt' field
	if deleteErr := this.db.Unscoped().Delete(&this.user).Error; deleteErr != nil {
		return dtos.CreateErrorDTO(deleteErr, 0, false)
	}

	return dtos.ErrorDTO{}
}


// ---------- Private ---------


func(this *UserService) setUser(userIdStr string) error {
	id, parseErr := uuid.Parse(userIdStr)
	if parseErr != nil {
		return parseErr
	}

	user, findErr := this.findUser(id)
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
	request.SetSubject("Teagan's WebApp Signup Confirmation")

	// generate html for email
	err := request.GenerateAndSetMessage()
	if err != nil {
		return err
	}

	// send email
	return request.SendEmail()
}
