package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
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

func NewTemplateGlobal(request engine.Request) *TemplateGlobal {
	templateGlobal := new(TemplateGlobal)
	templateGlobal.CASLogin = GetCasLoginUrl(request, "http://" + request.Host() + "/dashboard")
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

func GetCasLoginUrl(request engine.Request, service string) string {
	host := *Host

	if request != nil {
		host = request.Host();
	}

	casLogin, err := url.Parse("http://" + host + "/login?service=" + service)

	if err != nil {
		panic(err)
	}
	return casLogin.String()
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
