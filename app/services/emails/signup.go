package emails

type SignupEmail struct {
	BaseEmailRequest
	FirstName			string
}

func(this *SignupEmail) GenerateAndSetMessage() error {
	msg, err := this.generateMessage(this, "app/templates/signup.gohtml")
	if err != nil {
		return err
	}

	this.SetMessage(msg)

	return nil
}

func(this SignupEmail) CheckReadyToExecute() error {
	return this.readyToExecute()
}

