package enums

import (
	"strings"
)

type UserRole int

const (
	REGULAR		UserRole = iota + 1
	ADMIN
	SUPERADMIN
)

var UserRolesArr []UserRole = []UserRole{REGULAR, ADMIN, SUPERADMIN}

var roleStrings = []string{"regular", "admin", "superadmin"}

func(role UserRole) String() string {
	return roleStrings[role]
}

func(role UserRole) IsValid() bool {
	return int(role) <= len(roleStrings) && int(role) > 0
}

var roleMap = map[string]UserRole {
	"regular":		REGULAR,
	"admin":		ADMIN,
	"superadmin":	SUPERADMIN,
}

func ParseRoleString(roleStr string) (UserRole, bool) {
	role, ok := roleMap[strings.ToLower(roleStr)]
	return role, ok
}

func ValToRole(val interface{}) (UserRole, bool) {
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
		return ParseRoleString(roleStr)
	}

	return REGULAR, false
}
