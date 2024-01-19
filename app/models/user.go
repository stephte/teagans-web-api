package models

import (
	"chi-users-project/app/utilities/enums"
	"chi-users-project/app/utilities/auth"
	"chi-users-project/app/utilities"
	"gorm.io/gorm"
	"strings"
	"errors"
	"time"
	"fmt"
)

type User struct {
	BaseModel
	FirstName					string				`gorm:"not null"`
	LastName					string				`gorm:"not null"`
	Email						string				`gorm:"uniqueIndex;not null"`
	Role						enums.UserRole		`gorm:"default:0"`
	PasswordResetToken			[]byte
	PasswordResetExpiration		int64
	Password					string				`gorm:"-"`
	EncryptedPassword			[]byte				`gorm:"not null"`
	PasswordLastUpdated			int64
}


func(u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}


func(this *User) BeforeCreate(tx *gorm.DB) error {
	this.Email = strings.ToLower(this.Email)

	if this.Role == 0 {
		this.Role = 1
	}

	if pwErr := this.handlePassword(); pwErr != nil {
		return pwErr
	}

	return nil
}


func(this *User) BeforeUpdate(tx *gorm.DB) (err error) {
	typ := utilities.GetType(tx.Statement.Dest)
	
	// normal User update is assumed to be with a map
	if typ == "map[string]interface {}" {
		mp, ok := tx.Statement.Dest.(map[string]interface{})
		if ok {
			return this.beforeSaveWithMap(mp, tx)
		}

		return errors.New("Internal Error")
	} else if typ == "models.User" {
		usr, ok := tx.Statement.Dest.(User)
		if ok {
			return this.beforeSaveWithModel(usr, tx)
		}

		return errors.New("Internal Error")
	}

	// Password update only works with db.Save(&user) (on purpose, not an issue)
	// this works for password handling since we Save the user with the password already on the User
	if this.Password != "" {
		if pwErr := this.handlePassword(); pwErr != nil {
			return pwErr
		}
	}

	return nil
}


// pure validation checks should go in AfterSave
func(this *User) AfterSave(tx *gorm.DB) (err error) {
	return this.IsValid()
}


func(this User) CheckPassword(givenPW string) bool {
	return auth.CompareStringWithHash(this.EncryptedPassword, givenPW)
}


func(this User) CheckPWResetToken(givenToken string) bool {
	return auth.CompareStringWithHash(this.PasswordResetToken, givenToken)
}

// returns error if not valid, nil if is valid
func(this User) IsValid() error {
	if !this.Role.IsValid() {
		return errors.New("Invalid User Role")
	}

	if !utilities.IsValidEmail(this.Email) {
		return errors.New("Invalid User Email")
	}

	return nil
}


// ---------- Private ----------


func(this *User) beforeSaveWithMap(data map[string]interface{}, tx *gorm.DB) error {
	// first check if email key exists
	genericEmail, exists := data["Email"]
	if exists {
		// then get email string from map and update email with lowercase email
		email, isString := genericEmail.(string)
		if isString {
			tx.Statement.SetColumn("Email", strings.ToLower(email))
		} else {
			return errors.New("Email must be a string")
		}
	}

	return nil
}


func(this *User) beforeSaveWithModel(data User, tx *gorm.DB) error {
	if data.Email != "" {
		tx.Statement.SetColumn("Email", strings.ToLower(data.Email))
	}

	return nil
}


func(this *User) handlePassword() error {
	if !auth.ValidatePassword(this.Password) {
		return errors.New("Password must be at least 8 characters")
	}

	if auth.CompareStringWithHash(this.EncryptedPassword, this.Password) {
		return errors.New("New password cannot be the same as the current password")
	}

	hash, err := auth.CreateHash(this.Password)
	if err != nil {
		return err
	}

	this.EncryptedPassword = hash
	this.PasswordLastUpdated = time.Now().Unix()

	return nil
}
