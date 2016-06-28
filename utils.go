package main

import (
	"github.com/labstack/echo"
	"io"
	"net/url"
)

func NewTemplateGlobal() *TemplateGlobal {
	templateGlobal := new(TemplateGlobal)
	templateGlobal.CASLogin = GetCasLoginUrl("http://" + *OSFHost + "/dashboard")
	templateGlobal.OSFCreateAccount = GetOsfUrl("/register")
	templateGlobal.OSFDomain = GetOsfUrl("/")
	templateGlobal.OSFForgotPassword = GetOsfUrl("/forgotpassword")
	templateGlobal.OSFInstitutionLogin = GetOsfUrl("/login?campaign=institution")
	templateGlobal.OSFResendConfirmation = GetOsfUrl("/resend")
	return templateGlobal
}

func GetOsfUrl(path string) string {
	osfUrl, err := url.Parse("http://" + *OSFHost + path)
	if err != nil {
		panic(err)
	}
	return osfUrl.String()
}

func GetCasLoginUrl(service string) string {
	casLogin, err := url.Parse("http://" + *Host + "/login?service=" + service)
	if err != nil {
		panic(err)
	}
	return casLogin.String()
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
