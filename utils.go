package main

import (
	"github.com/labstack/echo"
	"io"
	"net/url"
)

func ValidateService(c echo.Context) *url.URL {
	service, err := url.Parse(c.QueryParam("service"))
	// TODO: add service validation if needed
	if err != nil {
		service = nil
	}
	return service
}

func NewTemplateGlobal() *TemplateGlobal {
	templateGlobal := new(TemplateGlobal)
	templateGlobal.CASLogin = GetCasLoginUrl("http://" + *OSFHost + "/dashboard")
	templateGlobal.CASRegister = GetCasRegisterUrl("http://" + *OSFHost + "/dashboard")
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
	return "/login?service=" + url.QueryEscape(service)
}

func GetCasRegisterUrl(service string) string {
	return "/account/register?service=" + url.QueryEscape(service)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
