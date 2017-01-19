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

	// find the user by email
	var queryString = `
		SELECT 
			id
			, username
			, given_name
			, family_name
		FROM osf_osfuser
		WHERE username = $1 OR $1 = ANY(emails)
	`
	err := DatabaseConnection.QueryRow(queryString, ticket).Scan(&result.Id, &result.Username, &result.GivenName, &result.FamilyName)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		fmt.Println("User", ticket, "not found.")
		return c.NoContent(http.StatusNotFound)
	}
	fmt.Println("User Found:", result.Username)

	// find the guid by user
	var guid string
	queryString = `
		SELECT DISTINCT _id
		FROM osf_guid
		LEFT JOIN django_content_type
			ON django_content_type.model = 'osfuser'
		JOIN osf_osfuser
			ON django_content_type.id = osf_guid.content_type_id AND object_id = osf_osfuser.id
			WHERE osf_osfuser.id = $1
		`
    err1 := DatabaseConnection.QueryRow(queryString, result.Id).Scan(&guid)
    if err != nil {
    	if err != sql.ErrNoRows {
    		panic(err1)
    	} 
    	fmt.Println("GUID not found")
    	return c.NoContent(http.StatusNotFound)
    }
    fmt.Println("GUID found:", guid)

    // build and return response to OSF
	response := ServiceResponse{
		Xmlns:       "http://www.yale.edu/tp/cas",
		User:        guid,
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

	// find the user and token scope by token
	var queryString = `
		SELECT
			osf_osfuser.id
			, osf_osfuser.username
			, osf_osfuser.given_name
			, osf_osfuser.family_name
			, osf_apioauth2personaltoken.scopes
		FROM osf_apioauth2personaltoken
		JOIN osf_osfuser
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
	fmt.Println("User", result.Username, "found for token")

	// find the guid for the user
	var guid string
	queryString = `
		SELECT DISTINCT _id
		FROM osf_guid
		LEFT JOIN django_content_type
			ON django_content_type.model = 'osfuser'
		JOIN osf_osfuser
			ON django_content_type.id = osf_guid.content_type_id AND object_id = osf_osfuser.id
			WHERE osf_osfuser.id = $1
	`
    err1 := DatabaseConnection.QueryRow(queryString, result.Id).Scan(&guid)
    if err != nil {
    	if err != sql.ErrNoRows {
    		panic(err1)
    	} 
    	fmt.Println("GUID not found")
    	return c.NoContent(http.StatusNotFound)
    }
    fmt.Println("GUID found:", guid)

    // return guid with attributes to API
	return c.JSON(200, OAuthResponse{
		Id: guid,
		Attributes: OAuthAttributes{
			LastName:  result.FamilyName,
			FirstName: result.GivenName,
		},
		Scope: strings.Split(scopes, " "),
	})
}
