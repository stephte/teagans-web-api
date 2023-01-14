package dtos

import (
	"chi-users-project/app/utilities"
	"strings"
	"errors"
	"fmt"
)

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


func genQueryStr(path string, limit int, page int, sort string) string {
	return fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", path, limit, page, sort)
}
