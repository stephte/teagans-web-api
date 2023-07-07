package services

import (
	"chi-users-project/app/utilities/enums"
	"chi-users-project/app/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"fmt"
)

type BaseService struct {
	db 				*gorm.DB
	log				zerolog.Logger
	currentUser		models.User
}


func (bs *BaseService) setCurrentUser(id uuid.UUID) error {
	user, findErr := bs.findUser(id)
	if findErr != nil {
		bs.log.Error().Msg(fmt.Sprintf("Can't find user with id: %s", id.String()))
		return findErr
	}

	bs.currentUser = user

	return nil
}


func (bs *BaseService) setCurrentUserByEmail(email string) error {
	user, findErr := bs.findUserByEmail(email)
	if findErr != nil {
		bs.log.Error().Msg(fmt.Sprintf("Can't find user with email: %s", email))
		return findErr
	}

	bs.currentUser = user

	return nil
}


func (bs BaseService) findUser(id uuid.UUID) (models.User, error) {
	user := models.User{}
	if findErr := bs.db.First(&user, id).Error; findErr != nil {
		return user, findErr
	}

	return user, nil
}


func (bs BaseService) findUserByEmail(userEmail string) (models.User, error) {
	user := models.User{}
	if findErr := bs.db.Where("Email = $1", userEmail).First(&user).Error; findErr != nil {
		return user, findErr
	}

	return user, nil
}


func (bs BaseService) validateUserHasAccess(accessNeeded enums.UserRole) bool {
	return bs.currentUser.Role >= accessNeeded
}
