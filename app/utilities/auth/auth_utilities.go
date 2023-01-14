package auth

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func CreateHash(str string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), 10)

	if err != nil {
		return hash, err
	}

	return hash, nil
}

func CompareStringWithHash(hash []byte, strToComp string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(strToComp))

	return err == nil
}

// TODO: add more rigorous password validation
func ValidatePassword(password string) bool {
	valid := false

	valid = (len([]rune(password)) >= 8) || valid

	return valid
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	length := len(letters)

	rand.Seed(time.Now().UnixNano())

	rv := make([]rune, n)
	for i := range rv {
		rv[i] = letters[rand.Intn(length)]
	}

	return string(rv)
}

// time params in milliseconds
func KillSomeTime(min int, max int) {
	// add a random amount of time (for security purposes)
	rand.Seed(time.Now().UnixNano())

	amount := rand.Intn(max-min) + min

	time.Sleep(time.Duration(amount) * time.Millisecond)
}


// ---------- User access control data -----------


func GetUserRoles() []int {
	return []int{
		RegularAccess(), // "regular",
		AdminAccess(), // "admin",
		SuperAdminAccess(), // "super admin"
	}
}

func RegularAccess() int {
	return 1
}
func AdminAccess() int {
	return 2
}
func SuperAdminAccess() int {
	return 3
}
