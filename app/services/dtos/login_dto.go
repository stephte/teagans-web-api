package dtos

type LoginDTO struct {
	Email		string
	Password	string
}

type LoginTokenDTO struct {
	Token		string		`json:"jwt"`
	CSRF		string		`json:"csrf"`
}

type EmailDTO struct {
	Email		string
}

type ConfirmResetTokenDTO struct {
	Email		string
	Token		string
}

type ResetPWDTO struct {
	Password	string
}
