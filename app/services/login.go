package services

import (
	"chi-users-project/app/services/emails"
	"chi-users-project/app/utilities/uuid"
	"chi-users-project/app/utilities/auth"
	"chi-users-project/app/services/dtos"
	"chi-users-project/app/models"
	"errors"
	"time"
)

type LoginService struct {
	*BaseService
}


func(this *LoginService) LoginUser(credentials dtos.LoginDTO, killTime bool) (dtos.LoginTokenDTO, int64, dtos.ErrorDTO) {
	// help protect against brute force attack
	if killTime {
		auth.KillSomeTime(967, 2978)
	}

	findErr := this.setCurrentUserByEmail(credentials.Email)
	if findErr != nil {
		return dtos.LoginTokenDTO{}, 0, dtos.CreateErrorDTO(errors.New("Email or Password Incorrect"), 401, false)
	}

	if !this.currentUser.CheckPassword(credentials.Password) {
		return dtos.LoginTokenDTO{}, 0, dtos.CreateErrorDTO(errors.New("Email or Password Incorrect"), 401, false)
	}

	// then create JWT token and return it
	token, csrf, maxAge, tokenErrDTO := this.genToken(false)

	return dtos.LoginTokenDTO{Token: token, CSRF: csrf}, maxAge, tokenErrDTO
}


func(this LoginService) StartPWReset(dto dtos.EmailDTO) (dtos.ErrorDTO) {
	findErr := this.setCurrentUserByEmail(dto.Email)
	if findErr != nil {
		return dtos.ErrorDTO{}
	}

	randToken := auth.RandomString(10)
	tokenHash, hashErr := auth.CreateHash(randToken)
	if hashErr != nil {
		return dtos.CreateErrorDTO(hashErr, 500, false)
	}

	// reset token expires in 1 hour
	expirationTS := time.Now().Add(time.Hour * 1).Unix()

	if updateErr := this.db.Model(&this.currentUser).Updates(models.User{PasswordResetToken: tokenHash, PasswordResetExpiration: expirationTS}).Error; updateErr != nil {
		return dtos.CreateErrorDTO(updateErr, 500, false)
	}

	// now spin off goroutine and send email with token
	go this.sendPWResetToken(randToken, dto.Email)

	return dtos.ErrorDTO{}
}


func(this *LoginService) ConfirmResetToken(dto dtos.ConfirmResetTokenDTO) (dtos.LoginTokenDTO, dtos.ErrorDTO) {
	findErr := this.setCurrentUserByEmail(dto.Email) 
	if findErr != nil {
		return dtos.LoginTokenDTO{}, dtos.AccessDeniedError()
	}

	if !this.currentUser.CheckPWResetToken(dto.Token) {
		return dtos.LoginTokenDTO{}, dtos.AccessDeniedError()
	}

	tokenExpired := this.currentUser.PasswordResetExpiration < time.Now().Unix()

	// clear out the User's Reset token
	if updateErr := this.db.Model(&this.currentUser).Select("PasswordResetToken", "PasswordResetExpiration").Updates(models.User{PasswordResetToken: nil, PasswordResetExpiration: 0}).Error; updateErr != nil {
		return dtos.LoginTokenDTO{}, dtos.CreateErrorDTO(updateErr, 500, false)
	}

	if tokenExpired {
		return dtos.LoginTokenDTO{}, dtos.CreateErrorDTO(errors.New("Token expired"), 0, false)
	}

	// create JWT token with PRT set to true, expiration in 1 hour (or less)
	// TODO: add csrf/maxAge handling (if using cookies for PW reset... which may not do)
	token, _, _, tokenErrDTO := this.genToken(true)

	return dtos.LoginTokenDTO{Token: token}, tokenErrDTO
}


func(this LoginService) UpdateUserPassword(dto dtos.ResetPWDTO) (dtos.LoginTokenDTO, int64, dtos.ErrorDTO) {
	this.currentUser.Password = dto.Password
	if saveErr := this.db.Save(&this.currentUser).Error; saveErr != nil {
		this.log.Error().Msg(saveErr.Error())
		return dtos.LoginTokenDTO{}, 0, dtos.CreateErrorDTO(saveErr, 0, false)
	}

	// then create new JWT token and return it
	token, csrf, maxAge, tokenErrDTO := this.genToken(false)

	return dtos.LoginTokenDTO{Token: token, CSRF: csrf}, maxAge, tokenErrDTO
}


// ---------- Private Methods ----------

//												 jwt	 csrf	maxAge   any error
func(this LoginService) genToken(pwReset bool) (string, string, int64, dtos.ErrorDTO) {
	header := dtos.JWTHeaderDTO{
		Algorithm: "HS256",
		Type: "JWT",
	}

	csrf := uuid.New().String()

	payload := dtos.JWTPayloadDTO{
		ID: this.currentUser.ID.String(),
		Issuer: "chi-users-project",
		CreatedAt: time.Now().Unix(),
		CSRF: csrf,
	}

	now := time.Now()
	if pwReset {
		payload.PRT = true
		payload.Expiration = time.Now().Add(time.Minute * 20).Unix()
	} else {
		payload.Expiration = now.Add(time.Hour * 4).Unix()
	}

	maxAge := payload.Expiration - now.Unix()

	authService := AuthService{this.BaseService}
	jwt, jerr := authService.GenerateJWT(header, payload)

	return jwt, csrf, maxAge, jerr
}


func(this LoginService) sendPWResetToken(token string, email string) error {
	this.log.Debug().Msg("Password reset email flow triggered")

	// set up email request
	request := emails.PWResetEmail {
		BaseEmailRequest: emails.InitBaseRequest(),
		Token: token,
	}
	request.SetToEmails([]string{email})
	request.SetSubject("Chi App Password Reset")

	// generate html for email
	err := request.GenerateAndSetMessage()
	if err != nil {
		return err
	}

	// send email
	return request.SendEmail()
}
