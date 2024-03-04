package enums

import (
	"teagans-web-api/app/utilities"
	"strings"
)

type UserRole int

const (
	REGULAR		UserRole = iota + 1
	ADMIN
	SUPERADMIN
)

var roleStrings = []string{"regular", "admin", "superadmin"}

func(role UserRole) String() string {
	if role.IsValid() {
		return strings.Title(roleStrings[role-1])
	}

	return ""
}

func(role UserRole) IsValid() bool {
	return int(role) > 0 && int(role) <= len(roleStrings)
}

func NewUserRole(num int64) (UserRole, bool) {
	rv := UserRole(num)

	if rv.IsValid() {
		return rv, true
	} else {
		return REGULAR, false
	}
}

func ParseUserRoleString(str string) (UserRole, bool) {
	ndx := utilities.IndexOf(roleStrings, strings.ToLower(str))

	if ndx >= 0 {
		return UserRole(ndx + 1), true
	} else {
		return REGULAR, false
	}
}

func ValToUserRole(val interface{}) (UserRole, bool) {
	var roleInt int
	var roleFloat32 float32
	var roleFloat64 float64
	var roleStr		string
	var ok bool

	roleInt, ok = val.(int)
	if ok {
		return UserRole(roleInt), true
	}
	roleFloat32, ok = val.(float32)
	if ok {
		return UserRole(roleFloat32), true
	}
	roleFloat64, ok = val.(float64)
	if ok {
		return UserRole(roleFloat64), true
	}
	roleStr, ok = val.(string)
	if ok {
		return ParseUserRoleString(roleStr)
	}

	return REGULAR, false
}
