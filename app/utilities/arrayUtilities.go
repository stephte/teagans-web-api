package utilities

import (
	"reflect"
)


// generic indexOf to be used for any type
func IndexOf[T comparable](arr []T, value T) int {
	for ndx, el := range arr {
		if value == el {
			return ndx
		}
	}
	return -1
}

func ArrContains[T comparable](arr []T, value T) bool {
	return IndexOf(arr, value) >= 0
}

func StructFieldIndexOf(arr []reflect.StructField, value string) int {
	for ndx, el := range arr {
		if value == el.Tag.Get("json") || value == el.Name {
			return ndx
		}
	}

	return -1
}
