package uuid

import (
	"github.com/google/uuid"
	"database/sql/driver"
	"encoding/json"
)

type UUID uuid.UUID

func New() UUID {
	return UUID(uuid.New())
}

func(this *UUID) Scan(value interface{}) error {
	var u uuid.UUID
	if value == nil {
		*this = UUID(u)
	} else {
		if err := u.Scan(value); err != nil {
        	return err
    	}

		*this = UUID(u)
	}

	return nil
}

func(this UUID) Value() (driver.Value, error) {
	if this == UUID(uuid.UUID{}) {
		return nil, nil
	}

	return uuid.UUID(this).Value()
}

func(this UUID) String() string {
	if this == UUID(uuid.UUID{}) {
		return ""
	}
	return uuid.UUID(this).String()
}

func(this UUID) Exists() bool {
	return this.String() != ""
}

func Parse(uuidStr string) (UUID, error) {
	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return New(), err
	}

	return UUID(u), nil
}

func(this UUID) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.String())
}

func(this *UUID) UnmarshalJSON(byt []byte) error {
	var str string
	err := json.Unmarshal(byt, &str)
	if err != nil {
		return err
	}

    u, parsErr := uuid.Parse(str)
    if parsErr != nil {
        return parsErr
    }

    *this = UUID(u)

    return nil
}
