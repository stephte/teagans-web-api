package utilities

import(
	"teagans-web-api/app/utilities/uuid"
	"reflect"
	"strings"
	"errors"
	"fmt"
)


func ValidateMapWithStruct(mapToValidate map[string]interface{}, strct interface{}) (map[string]interface{}, error) {
	rv := map[string]interface{}{}

	strctFieldList := GetStructFields(strct)

	for key, value := range mapToValidate {
		// if key exists in the strct, continue on
		if ndx := StructFieldIndexOf(strctFieldList, key); ndx >= 0 {
			structField := strctFieldList[ndx]
			keyName := structField.Name

			typeMatch := reflect.TypeOf(value) == structField.Type
			if !typeMatch {
				if IsNumberType(structField.Type.String()) {
					typeMatch = IsNumberType(GetType(value))
				} else if structField.Type == reflect.TypeOf(uuid.New()) {
					_, typeMatch = value.(string)
				}
			}

			if typeMatch {
				rv[keyName] = value
			} else {
				return map[string]interface{}{}, errors.New(fmt.Sprintf("Invalid data type passed: key '%s' should be of type '%s', but is type '%T'", key, structField.Type.String(), value))
			}
		}
	}

	return rv, nil
}


func StructToMap(strct interface{}) (map[string]interface{}, error) {
	rv := map[string]interface{}{}

	structVal := reflect.ValueOf(strct)
	if structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	}

	if structVal.Kind() != reflect.Struct {
		return rv, errors.New("Param must be a struct")
	}

	typ := structVal.Type()
	for i := 0; i < structVal.NumField(); i++ {
		key := typ.Field(i).Name
		if key != "" {
			rv[key] = structVal.Field(i).Interface()
		}
	}

	return rv, nil
}


// return slice with key names (does this method work on map[string]interface?)
func GetStructFields(strct interface{}) ([]reflect.StructField) {
	var rv []reflect.StructField
	structVal := reflect.ValueOf(strct)
	if structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	}

	typ := structVal.Type()
	for i := 0; i < structVal.NumField(); i++ {
		// fmt.Println("typ.field:", typ.Field(i))
		field := typ.Field(i)
		if field.Name != "" {
			rv = append(rv, field)
		}
	}

	return rv
}


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
