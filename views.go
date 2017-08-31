package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func RegisterGET(c echo.Context) error {

	data := NewTemplateGlobal()

	service := ValidateService(c)
	if service == nil {
		data.NotAuthorized = true
		return c.Render(http.StatusUnauthorized, "register", data)
	}
	data.CASRegister = GetCasRegisterUrl(service.String())
	data.CASLogin = GetCasLoginUrl(service.String())
	return c.Render(http.StatusOK, "register", data)
}

func RegisterPOST(c echo.Context) error {

	data := NewTemplateGlobal()

	service := ValidateService(c)
	if service == nil {
		data.NotAuthorized = true
		return c.Render(http.StatusUnauthorized, "register", data)
	}
	data.CASRegister = GetCasRegisterUrl(service.String())
	data.CASLogin = GetCasLoginUrl(service.String())

	fullname := strings.TrimSpace(c.FormValue("fullname"))
	email := strings.ToLower(strings.TrimSpace(c.FormValue("email")))
	password := c.FormValue("password")

	api_v1_url_for_register_user := "http://" + *OSFHost + "/api/v1/register/"
	jsonStr := "{\"email1\": \"" + email + "\", \"email2\": \"" + email + "\", \"password\": \"" + password + "\", \"fullName\": \"" + fullname + "\"}"
	byteStr := []byte(jsonStr)
	resp, err := http.Post(api_v1_url_for_register_user, "application/json", bytes.NewBuffer(byteStr))
	if err != nil {
		panic(err)
		data.ShowErrorMessages = true
		return c.Render(http.StatusOK, "register", data)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response Body:\t", string(body))
	if resp.StatusCode == http.StatusOK {
		data.RegisterSuccessful = true
	} else {
		data.ShowErrorMessages = true
	}

	return c.Render(http.StatusOK, "register", data)
}

func LoginPOST(c echo.Context) error {
	data := NewTemplateGlobal()

	service := ValidateService(c)
	if service == nil {
		data.NotAuthorized = true
		return c.Render(http.StatusUnauthorized, "login", data)
	}

	var isRegistered bool
	username := strings.ToLower(strings.TrimSpace(c.FormValue("username")))

	// fakeCAS does not check password
	err := DatabaseConnection.QueryRow(`
		SELECT is_registered
		FROM osf_osfuser
		WHERE username = $1
		OR EXISTS(SELECT * FROM osf_email WHERE osf_email.user_id = osf_osfuser.id AND osf_email.address = $1)
	`, username).Scan(&isRegistered)

	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		fmt.Println("User", username, "not found.")
		data.LoginForm = true
		data.NotExist = true
		data.CASLogin = GetCasLoginUrl(service.String())
		return c.Render(http.StatusOK, "login", data)
	}

	if !isRegistered {
		data.NotRegistered = true
		return c.Render(http.StatusOK, "login", data)
	}

	query := service.Query()
	query.Set("ticket", username)
	service.RawQuery = query.Encode()

	fmt.Println("Logging in and redirecting to", service)
	return c.Redirect(http.StatusFound, service.String())
}

func LoginGET(c echo.Context) error {
	data := NewTemplateGlobal()

	service := ValidateService(c)
	if service == nil {
		data.NotAuthorized = true
		return c.Render(http.StatusUnauthorized, "login", data)
	}

	username, err := url.Parse(c.QueryParam("username"))
	if err != nil {
		c.Error(err)
		return nil
	}

	verificationKey, err := url.Parse(c.QueryParam("verification_key"))
	if err != nil {
		c.Error(err)
		return nil
	}

	if username.String() == "" && verificationKey.String() == "" {
		data.LoginForm = true
		data.CASLogin = GetCasLoginUrl(service.String())
		return c.Render(http.StatusOK, "login", data)
	}

	var verification string
	uname := strings.ToLower(strings.TrimSpace(c.FormValue("username")))
	err = DatabaseConnection.QueryRow(`
		SELECT verification_key
		FROM osf_osfuser
		WHERE username = $1
		OR EXISTS(SELECT * FROM osf_email WHERE osf_email.user_id = osf_osfuser.id AND osf_email.address = $1)
	`, uname).Scan(&verification)

	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		fmt.Println("User", uname, "not found.")
		data.NotValid = true
		return c.Render(http.StatusNotFound, "login", data)
	}

	// fakeCAS will check verification key
	if verification != c.FormValue("verification_key") {
		fmt.Println("Invalid Verification Key\nExpecting: ", verification, "\nActual: ", c.FormValue("verification_key"))
		data.NotValid = true
		return c.Render(http.StatusNotFound, "login", data)
	}

	query := service.Query()
	query.Set("ticket", c.FormValue("username"))
	service.RawQuery = query.Encode()

	fmt.Println("Logging in and redirecting to", service)
	return c.Redirect(http.StatusFound, service.String())
}

func Logout(c echo.Context) error {
	fmt.Println("Logging out and redirecting to", c.FormValue("service"))
	return c.Redirect(http.StatusFound, c.FormValue("service"))
}

func ServiceValidate(c echo.Context) error {
	result := User{}
	ticket := c.FormValue("ticket")

	var queryString = `
		SELECT DISTINCT
			osf_guid._id,
			osf_osfuser.username,
			osf_osfuser.given_name,
			osf_osfuser.family_name
		FROM osf_osfuser
			RIGHT JOIN osf_guid
				ON osf_guid.object_id = osf_osfuser.id
		WHERE osf_guid.content_type_id = (SELECT id FROM django_content_type WHERE model = 'osfuser' LIMIT 1)
		AND (
			EXISTS (SELECT * FROM osf_email WHERE osf_email.user_id = osf_osfuser.id AND osf_email.address = $1)
			OR (osf_osfuser.username = $1)
		);
	`
	err := DatabaseConnection.QueryRow(queryString, ticket).Scan(&result.Id, &result.Username, &result.GivenName, &result.FamilyName)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		fmt.Println("User", ticket, "not found.")
		return c.NoContent(http.StatusNotFound)
	}
	fmt.Println("User found: username =", result.Username, ", guid =", result.Id)

	response := ServiceResponse{
		Xmlns:       "http://www.yale.edu/tp/cas",
		User:        result.Id,
		NewLogin:    true,
		Date:        "Eh",
		GivenName:   result.GivenName,
		FamilyName:  result.FamilyName,
		AccessToken: result.Id,
		UserName:    result.Username,
	}
	return c.XML(http.StatusOK, response)
}

func OAuth(c echo.Context) error {
	var (
		scopes string
		result User
	)

	tokenId := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1)

	var queryString = `
		SELECT DISTINCT
			osf_guid._id,
			osf_osfuser.username,
			osf_osfuser.given_name,
			osf_osfuser.family_name,
			osf_apioauth2personaltoken.scopes
		FROM osf_guid
			LEFT JOIN django_content_type
				ON django_content_type.model = 'osfuser'
			JOIN osf_osfuser
				ON django_content_type.id = osf_guid.content_type_id AND object_id = osf_osfuser.id
			JOIN osf_apioauth2personaltoken
					ON osf_osfuser.id = osf_apioauth2personaltoken.owner_id
			WHERE osf_apioauth2personaltoken.token_id = $1
	`
	err := DatabaseConnection.QueryRow(queryString, tokenId).Scan(&result.Id, &result.Username, &result.GivenName, &result.FamilyName, &scopes)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		fmt.Println("Access token", tokenId, "not found")
		return c.NoContent(http.StatusNotFound)
	}
	fmt.Println("User found for token: username =", result.Username, ", guid =", result.Id)

	return c.JSON(200, OAuthResponse{
		Id: result.Id,
		Attributes: OAuthAttributes{
			LastName:  result.FamilyName,
			FirstName: result.GivenName,
		},
		Scope: strings.Split(scopes, " "),
	})
}
