package dtos

import (
	"teagans-web-api/app/utilities/interfaceUtils"
	"strings"
	"errors"
	"fmt"
)


type UserInDTO struct {
	FirstName 	string		`json:"firstName"`
	LastName	string		`json:"lastName"`
	Email		string		`json:"email"`
	Role		int64		`json:"role" enum:"UserRole"`
}

type UserOutDTO struct {
	BaseDTO
	UserInDTO
}

// How to get password so we can retrieve it, but not send it??
type CreateUserDTO struct {
	UserInDTO
	Password	string
}


// // uses UserDTO to validate the data passed before updating the user with it
// func ValidateUserMap(data map[string]interface{}) (map[string]interface{}, error) {
// 	rv := map[string]interface{}{}
// 	var strct UserInDTO

// 	for key, value := range data {
// 		capitalKey := strings.Title(key)

// 		// ensure they can't manipulate BaseDTO data
// 		if strings.Contains(capitalKey, "BaseDTO") {
// 			return rv, errors.New(fmt.Sprintf("Unpermitted data: %s", key))
// 		} else if strings.Contains(strings.ToLower(key), "id") {
// 			continue
// 		}

// 		typ := interfaceUtils.GetType(value)

// 		// custom check for role, since it is an enum
// 		if capitalKey == "Role" {
// 			if typ == "float64" || typ == "string" {
// 				rv[capitalKey] = value
// 				continue
// 			}
// 		}

// 		validKey, typErr := interfaceUtils.CheckKeyValue(strct, capitalKey, typ)

// 		if typErr != nil {
// 			return rv, typErr
// 		} else if !validKey {
// 			return rv, errors.New(fmt.Sprintf("Unpermitted data: %s", key))
// 		}
		
// 		rv[capitalKey] = value
// 	}

// 	// if interface is empty, return error
// 	if len(rv) == 0 {
// 		return rv, errors.New("No valid user data found")
// 	}

// 	return rv, nil
// }
