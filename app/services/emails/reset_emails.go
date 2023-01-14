package emails

import (
	"errors"
)

// implements EmailRequest interface
type PWResetEmail struct {
	BaseEmailRequest
	Token		string
}

func(this *PWResetEmail) GenerateAndSetMessage() error {
	msg, err := this.generateMessage(this, "app/templates/pw_reset.gohtml")
	if err != nil {
		return err
	}

	this.SetMessage(msg)

	return nil
}

func(this PWResetEmail) CheckReadyToExecute() error {
	if this.Token == "" {
		return errors.New("Token missing")
	}

	return this.readyToExecute()
}
