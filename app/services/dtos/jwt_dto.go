package dtos

import (
	"time"
)

// ----- JWT Header -----

type JWTHeaderDTO struct {
	Algorithm		string		`json:"alg"`
	Type			string		`json:"typ"`
}

func(this JWTHeaderDTO) Exists() bool {
	return this.Algorithm != "" && this.Type != ""
}


// ----- JWT Payload ------

type JWTPayloadDTO struct {
	ID				string		`json:"id"`
	CreatedAt		int64		`json:"cre"`	// probs dont need
	Expiration		int64		`json:"exp"`
	Issuer			string		`json:"iss"`
	PRT				bool		`json:"prt"`	// is password reset
	CSRF			string		`json:"csrf"`
}

func(this JWTPayloadDTO) Exists() bool {
	return this.ID != "" && this.Expiration != 0 && this.Issuer != ""
}

func(this JWTPayloadDTO) IsActive() bool {
	return time.Now().Unix() < this.Expiration
}
