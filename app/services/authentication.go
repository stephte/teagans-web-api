package services

import (
	"chi-users-project/app/services/dtos"
	"github.com/google/uuid"
	"encoding/base64"
	"encoding/json"
	"crypto/sha256"
	"crypto/hmac"
	"strings"
	"errors"
	"fmt"
	"os"
)

type AuthService struct {
	*BaseService
}


func(this AuthService) GenerateJWT(header dtos.JWTHeaderDTO, payload dtos.JWTPayloadDTO) (string, dtos.ErrorDTO) {
	if !header.Exists() || !payload.Exists() {
		this.log.Warn().Msg(fmt.Sprintf("Invalid JWT Params!\nHeaderDTO: %v;\nPayloadDTO: %v", header, payload))
		return "", dtos.CreateErrorDTO(errors.New("Error!"), 500, false)
	}
	
	payloadJSON, pMarshalErr := json.Marshal(payload)
	if pMarshalErr != nil {
		this.log.Error().Err(pMarshalErr).Msg("")
		return "", dtos.CreateErrorDTO(errors.New("Error!"), 500, false)
	}

	headerJSON, hMarshalErr := json.Marshal(header)
	if hMarshalErr != nil {
		this.log.Error().Err(hMarshalErr).Msg("")
		return "", dtos.CreateErrorDTO(errors.New("Error!"), 500, false)
	}

	encodedHeader := base64.StdEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.StdEncoding.EncodeToString(payloadJSON)

	signature := this.generateSignature(encodedHeader, encodedPayload)

	token := fmt.Sprintf("%s.%s.%s", encodedHeader, encodedPayload, signature)

	return token, dtos.ErrorDTO{}
}


func(this *AuthService) ValidateJWT(jwt, csrf string, isPWReset bool) (bool, dtos.ErrorDTO) {
	if jwt == "" {
		this.log.Warn().Msg("No token found")
		return this.invalidTokenErr(true)
	}

	splitJWT := strings.Split(jwt, ".")

	if len(splitJWT) != 3 {
		this.log.Warn().Msg(fmt.Sprintf("Token not in 3 parts: %s\n", jwt))
		return this.invalidTokenErr(false)
	}

	encodedHeader := splitJWT[0]
	encodedPayload := splitJWT[1]
	jwtSignature := splitJWT[2]

	payloadJSON, decodeErr := base64.StdEncoding.DecodeString(encodedPayload)

	if decodeErr != nil {
		this.log.Error().Err(decodeErr).Msg("")
		return this.invalidTokenErr(false)
	}

	var payload dtos.JWTPayloadDTO
	marshalErr := json.Unmarshal(payloadJSON, &payload)
	if marshalErr != nil {
		this.log.Error().Err(marshalErr).Msg("")
		return this.invalidTokenErr(false)
	}

	// not using CSRF for PW reset tokens for now
	if !isPWReset {
		if payload.CSRF != csrf {
			this.log.Error().Msg("Invalid CSRF")
			return this.invalidTokenErr(false)
		}
	}

	userId, parseErr := uuid.Parse(payload.ID)
	if parseErr != nil {
		this.log.Error().Msg(fmt.Sprintf("Error parsing UUID: %s\n", payload.ID))
		return this.invalidTokenErr(false)
	}

	findErr := this.setCurrentUser(userId)
	if findErr != nil {
		this.log.Error().Msg(fmt.Sprintf("User not found: %s\n", payload.ID))
		return this.invalidTokenErr(false)
	}

	signature := this.generateSignature(encodedHeader, encodedPayload)

	if jwtSignature != signature {
		this.log.Warn().Msg(fmt.Sprintf("Signatures do not match: %s::%s\n", jwtSignature, signature))
		return this.invalidTokenErr(false)
	}

	if isPWReset != payload.PRT {
		this.log.Warn().Msg("Resets don't match")
		return this.invalidTokenErr(false)
	} else if !payload.IsActive() {
		this.log.Warn().Msg("Token expired")
		return false, dtos.ErrorDTO{}
	}

	return true, dtos.ErrorDTO{}
}


func(this *AuthService) generateSignature(header string, payload string) string {
	strToSign := fmt.Sprintf("%s%s", header, payload)

	signingKey := fmt.Sprintf("%s%s", os.Getenv("AUTH_SECRET"), this.currentUser.EncryptedPassword)

	hasher := hmac.New(sha256.New, []byte(signingKey))
	hasher.Write([]byte(strToSign))

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}


func(this AuthService) invalidTokenErr(relogin bool) (bool, dtos.ErrorDTO) {
	return false, dtos.AccessDeniedError(relogin)
}
