package intrface

import (
	"teagans-web-api/app/utilities/enums"
	"teagans-web-api/app/utilities"
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
		if ndx := utilities.StructFieldIndexOf(strctFieldList, key); ndx >= 0 {
			structField := strctFieldList[ndx]
			keyName := structField.Name
			typeMatch := false

			// handle case where its an enum passed in
			enumTyp := structField.Tag.Get("enum")
			if enumTyp != "" {
				var methodsMap map[string]interface{}
				_, stringOk := value.(string)
				dubVal, doubleOk := value.(float64)

				if stringOk {
					methodsMap = enums.GetParseMethodsMap()
				} else if doubleOk {
					value = int64(dubVal)
					methodsMap = enums.GetNewMethodsMap()
				} else {
					return map[string]interface{}{}, errors.New(fmt.Sprintf("Invalid data type passed: key '%s' should be of type 'string', but is type '%T'", key, value))
				}

				enumInt, enumValid := getEnumValue(value, enumTyp, methodsMap)
				if enumValid {
					value = enumInt
					typeMatch = true
				} else {
					return map[string]interface{}{}, errors.New(fmt.Sprintf("Invalid enum value passed for '%s'", key))
				}
			}

			typeMatch = typeMatch || reflect.TypeOf(value) == structField.Type
			if !typeMatch {
				strTypes := []string{"uuid.UUID"}

				if IsNumberType(structField.Type.String()) {
					typeMatch = IsNumberType(GetType(value))
				} else if structField.Type.String() == "*time.Time" {
					if value != nil {
						_, typeMatch = value.(string)
					} else {
						typeMatch = true
					}
				} else if utilities.ArrContains(strTypes, structField.Type.String()) {
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

// ---------- private ----------

func getEnumValue(value interface{}, enumTyp string, methodsMap map[string]interface{}) (int64, bool) {
	meth, ok := methodsMap[enumTyp]
	if !ok {
		fmt.Println(fmt.Sprintf("Enum parse function %s not found in MethodsMap", enumTyp))
		return 0, false
	}

	fnc := reflect.ValueOf(meth)
	res := fnc.Call([]reflect.Value{ reflect.ValueOf(value) })

	enumVal := res[0]
	enumBool := res[1]

	return enumVal.Int(), enumBool.Bool()
}
