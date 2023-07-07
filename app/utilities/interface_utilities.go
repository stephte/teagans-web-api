package utilities

import(
	"reflect"
	"strings"
	"errors"
	"fmt"
)


// returns true if mapKey is a field in struct, and if they are of the same data type
func CheckKeyValue(strct interface{}, mapKey string, mapValType string) (bool, error) {
	structVal := reflect.ValueOf(strct)

	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		name := structVal.Type().Field(i).Name
		typ := field.Type().String()

		if name == mapKey {
			if typ == mapValType || (IsNumberType(mapValType) && IsNumberType(typ)) {
				return true, nil
			} else {
				errStr := fmt.Sprintf("%s is of incorrect data type; is %s, should be %s", name, mapValType, typ)
				return false, errors.New(errStr)
			}
		}
	}

	return false, nil
}


// returns value type as string
func GetType(value interface{}) string {
	return fmt.Sprintf("%T", value)
}


func IsNumberType(typStr string) bool {
	if strings.Contains(strings.ToLower(typStr), "float") || strings.Contains(strings.ToLower(typStr), "int") {
		return true
	}

	return false
}
