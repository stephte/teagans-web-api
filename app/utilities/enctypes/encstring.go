package enctypes

import (
	"teagans-web-api/app/utilities/ciphers"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
)

type EncString string

func(this *EncString) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		*this = EncString("")
		return nil
	case string:
		if len(value) == 0 {
			*this = EncString("")
			return nil
		}
		// decrypt it
		if false {
			ciphertext, err := hex.DecodeString(value)
			if err != nil {
				return err
			}

			rv, err := ciphers.DecryptFromAES(ciphertext)
			if err != nil {
				return err
			}

			*this = EncString(rv)
		} else {
			*this = EncString(value)
		}
	case []byte:
		if value == nil {
			*this = EncString("")
			return nil 
		}

		return this.Scan(string(value))
	default:
		return fmt.Errorf("Unable to scan type %T into EncString", value)
	}

	return nil
}

func(this EncString) Value() (driver.Value, error) {
	str := string(this)
	if str == "" {
		return "", nil
	}
	// encrypt the string
	rv, err := ciphers.EncryptWithAES([]byte(str))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rv), nil
}

// func(this UUID) String() string {
// 	if this == UUID(uuid.UUID{}) {
// 		return ""
// 	}
// 	return uuid.UUID(this).String()
// }
