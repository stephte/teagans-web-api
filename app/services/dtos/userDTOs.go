package dtos

import (
	"chi-users-project/app/utilities/enums"
	"chi-users-project/app/utilities"
	"strings"
	"errors"
	"fmt"
)

type UserDTO struct {
	BaseDTO
	FirstName 	string				`json:"firstName"`
	LastName	string				`json:"lastName"`
	Email		string				`json:"email"`
	Role		enums.UserRole		`json:"role"`
}

// How to get password so we can retrieve it, but not send it??
type CreateUserDTO struct {
	FirstName 	string
	LastName	string
	Email		string
	Role		enums.UserRole
	Password	string
}


// uses UserDTO to validate the data passed before updating the user with it
func ValidateUserMap(data map[string]interface{}) (map[string]interface{}, error) {
	rv := map[string]interface{}{}
	var strct UserDTO

	for key, value := range data {
		capitalKey := strings.Title(key)

		// ensure they can't manipulate BaseDTO data
		if strings.Contains(capitalKey, "BaseDTO") {
			return rv, errors.New(fmt.Sprintf("Unpermitted data: %s", key))
		}

		typ := utilities.GetType(value)

		// custom check for role, since it is an enum
		if capitalKey == "Role" {
			if typ == "float64" || typ == "string" {
				rv[capitalKey] = value
				continue
			}
		}

		validKey, typErr := utilities.CheckKeyValue(strct, capitalKey, typ)

		if typErr != nil {
			return rv, typErr
		} else if !validKey {
			return rv, errors.New(fmt.Sprintf("Unpermitted data: %s", key))
		}
		
		rv[capitalKey] = value
	}

	// if interface is empty, return error
	if len(rv) == 0 {
		return rv, errors.New("No valid user data found")
	}

	return rv, nil
}
