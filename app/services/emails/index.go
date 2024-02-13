package emails

import (
	"teagans-web-api/app/utilities"
	"net/smtp"
	"errors"
	"bytes"
	"fmt"
	"os"
)

type EmailRequest interface {
	GenerateAndSetMessage()		error
	CheckReadyToExecute()		error
}

type BaseEmailRequest struct {
	toEmails		[]string
	from			string
	cc				[]string
	subject			string
	message			[]byte
}


func(this BaseEmailRequest) SendEmail() error {
	if !this.ReadyToSend() {
		err := errors.New("Email request missing required data; unable to send")
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Sending email to: %s\n", this.GetToEmailString())
	auth := smtp.PlainAuth("", os.Getenv("CHI_YT_EMAIL_UNAME"), os.Getenv("CHI_YT_EMAIL_PW"), os.Getenv("CHI_YT_EMAIL_HOST"))

	hostWithPort := fmt.Sprintf("%s:%s", os.Getenv("CHI_YT_EMAIL_HOST"), os.Getenv("CHI_YT_EMAIL_PORT"))

	err := smtp.SendMail(hostWithPort, auth, this.from, this.GetToEmails(), this.message)

	if err != nil {
		fmt.Printf("Error sending email: %s\n", err.Error())
	} else {
		fmt.Println("Email successfully sent")
	}

	return err
}


func(this BaseEmailRequest) ReadyToSend() bool {
	return len(this.GetToEmails()) >= 1 && this.from != "" && this.message != nil
}


func(this BaseEmailRequest) readyToExecute() error {
	errMsg := ""
	if this.toEmails == nil {
		errMsg = "To Email(s) missing"
	} else if this.from == "" {
		errMsg = "From email missing"
	}

	if errMsg != "" {
		return errors.New(errMsg)
	}

	return nil
}


func(this *BaseEmailRequest) generateMessage(request EmailRequest, files ...string) ([]byte, error) {
	execErr := request.CheckReadyToExecute()
	if execErr != nil {
		fmt.Printf("Request not ready to execute: %s\n", execErr.Error())
		return []byte{}, execErr
	}

	temp, err := parseFiles(files...)
	if err != nil {
		fmt.Printf("ERROR parsing template: %s\n", err.Error())
		return []byte{}, err
	}

	buf := new(bytes.Buffer)
	err = temp.Execute(buf, request)
	if err != nil {
		fmt.Printf("ERROR executing template: %s\n", err.Error())
		return []byte{}, err
	}

	return []byte(buf.String()), nil
}


// -------- Getters/Setters --------


func(this BaseEmailRequest) GetToEmailString() string {
	emails := this.GetToEmails()

	if len(emails) == 0 {
		return ""
	}

	rv := emails[0]
	for i := 1; i < len(emails); i++ {
		rv = fmt.Sprintf("%s,%s", rv, emails[i])
	}

	return rv
}

func(this BaseEmailRequest) GetToEmails() []string {
	return this.toEmails
}

func(this BaseEmailRequest) GetFromEmail() string {
	return this.from
}

func(this BaseEmailRequest) GetCCEmails() []string {
	return this.cc
}

func(this BaseEmailRequest) GetSubject() string {
	return this.subject
}

func(this *BaseEmailRequest) SetToEmails(emails []string) {
	this.toEmails = utilities.FilterValidEmails(emails)
}

func(this *BaseEmailRequest) SetCCEmails(emails []string) {
	this.cc = utilities.FilterValidEmails(emails)
}

func(this *BaseEmailRequest) SetMessage(msg []byte) {
	this.message = msg
}

func(this *BaseEmailRequest) SetSubject(subject string) {
	this.subject = subject
}

func(this *BaseEmailRequest) SetFrom(email string) {
	this.from = email
}
