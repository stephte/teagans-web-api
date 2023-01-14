package utilities

import(
	"reflect"
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
			if typ == mapValType {
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
	typ := fmt.Sprintf("%T", value)

	floatValue, ok := value.(float64)
	if ok && floatValue == float64(int(floatValue)) {
		typ = "int"
	}

	return typ
}
