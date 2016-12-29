package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
	"strings"
)

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
	err := DatabaseConnection.QueryRow("SELECT is_registered FROM osf_osfuser WHERE username = $1 OR $1 = ANY(emails)", username).Scan(&isRegistered)

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
	err = DatabaseConnection.QueryRow("SELECT verification_key FROM osf_osfuser WHERE username = $1 OR $1 = ANY(emails)", uname).Scan(&verification)

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

	err := DatabaseConnection.QueryRow(`
		SELECT
			id
			, username
			, given_name
			, family_name
		FROM osf_osfuser
		WHERE username = $1 OR $1 = ANY(emails)
	`, ticket).Scan(&result.Id, &result.Username, &result.GivenName, &result.Username)

	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		fmt.Println("User", ticket, "not found.")
		return c.NoContent(http.StatusNotFound)
	}

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

	err := DatabaseConnection.QueryRow(`
		SELECT
			user.id
			, user.given_name
			, user.family_name
			, token.scopes
		FROM osf_apioauth2personaltoken AS token
		JOIN osf_osfuser AS user ON user.id = token.owner_id
		WHERE token_id = $1
	`, tokenId).Scan(&result.Id, &result.Username, &result.GivenName, scopes)

	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		fmt.Println("Access token", tokenId, "not found")
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(200, OAuthResponse{
		Id: result.Id,
		Attributes: OAuthAttributes{
			LastName:  result.FamilyName,
			FirstName: result.GivenName,
		},
		Scope: strings.Split(scopes, " "),
	})
}
