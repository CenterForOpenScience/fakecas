package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
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

	tokenId := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1)

	// Find the user that owns the token
	var result User
	queryString := `
		SELECT DISTINCT
			osf_guid._id,
			osf_osfuser.username,
			osf_osfuser.given_name,
			osf_osfuser.family_name
		FROM osf_guid
			LEFT JOIN django_content_type
				ON django_content_type.model = 'osfuser'
			JOIN osf_osfuser
				ON django_content_type.id = osf_guid.content_type_id AND object_id = osf_osfuser.id
			JOIN osf_apioauth2personaltoken
				ON osf_osfuser.id = osf_apioauth2personaltoken.owner_id
		WHERE osf_apioauth2personaltoken.token_id = $1
	`
	err := DatabaseConnection.QueryRow(queryString, tokenId).Scan(&result.Id, &result.Username, &result.GivenName, &result.FamilyName)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		fmt.Printf("Access token %s not found\n", tokenId)
		return c.NoContent(http.StatusNotFound)
	}
	fmt.Printf("User found for token: username = %s , guid =%s\n", result.Username, result.Id)

	// Find all the scopes associated with the token
	fmt.Printf("Reading scopes ... ")
	queryString = `
		SELECT DISTINCT osf_apioauth2scope.name
		FROM osf_apioauth2personaltoken_scopes
			JOIN osf_apioauth2personaltoken
		   		on osf_apioauth2personaltoken_scopes.apioauth2personaltoken_id = osf_apioauth2personaltoken.id
			JOIN osf_apioauth2scope
		   		on osf_apioauth2personaltoken_scopes.apioauth2scope_id = osf_apioauth2scope.id
		WHERE osf_apioauth2personaltoken.token_id = $1
	`
	rows, err := DatabaseConnection.Query(queryString, tokenId)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		fmt.Printf("No scope is found for access token %s\n", tokenId)
		return c.NoContent(http.StatusNotFound)
	}
	defer rows.Close()
	scopes := make([]string, 0)
	var scope string
	for rows.Next() {
		err = rows.Scan(&scope)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s, ", scope)
		scopes = append(scopes, scope)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Printf("... %d scopes in total.\n", len(scopes))

	// Return 200 with user information and scopes
	return c.JSON(200, OAuthResponse{
		Id: result.Id,
		Attributes: OAuthAttributes{
			LastName:  result.FamilyName,
			FirstName: result.GivenName,
		},
		Scope: scopes,
	})
}

func OAuthRevoke(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
