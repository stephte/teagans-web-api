package emails

import (
	"html/template"
	"os"
)


// first file passed should be main template
func parseFiles(files ...string) (*template.Template, error) {
	filesArr := append(files, "app/templates/smtpInfo.gohtml")
	return template.ParseFiles(filesArr...)
}


func InitBaseRequest() BaseEmailRequest {
	bs := BaseEmailRequest{
		from: os.Getenv("CHI_YT_APP_EMAIL_ADDR"),
	}

	return bs
}
