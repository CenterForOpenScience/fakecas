package main

import (
	"fmt"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
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
	result := User{}
	username := strings.ToLower(strings.TrimSpace(c.FormValue("username")))

	// fakeCAS does not check password
	if err := UserCollection.Find(bson.M{"username": username}).One(&result); err != nil {
		fmt.Println("User", username, "not found.")
		data.LoginForm = true
		data.NotExist = true
		data.CASLogin = GetCasLoginUrl(service.String())
		return c.Render(http.StatusOK, "login", data)
	}

	if !result.IsRegistered {
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

	verification_key, err := url.Parse(c.QueryParam("verification_key"))
	if err != nil {
		c.Error(err)
		return nil
	}

	if username.String() == "" && verification_key.String() == "" {
		data.LoginForm = true
		data.CASLogin = GetCasLoginUrl(service.String())
		return c.Render(http.StatusOK, "login", data)
	}

	result := User{}
	if err = UserCollection.Find(bson.M{"username": c.FormValue("username")}).One(&result); err != nil {
		fmt.Println("User", c.FormValue("username"), "not found.")
		data.NotValid = true
		return c.Render(http.StatusNotFound, "login", data)
	}
	// fakeCAS will check verification key
	if result.VerificationKey != c.FormValue("verification_key") {
		fmt.Println("Invalid Verification Key\nExpecting: ", result.VerificationKey,
			"\nActual: ", c.FormValue("verification_key"))
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

	if err := UserCollection.Find(bson.M{"emails": c.FormValue("ticket")}).One(&result); err != nil {
		fmt.Println("User", c.FormValue("ticket"), "not found.")
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

	token := AccessToken{}
	tokenId := strings.Replace(c.Request().Header().Get("Authorization"), "Bearer ", "", 1)
	err := AccessTokenCollection.Find(bson.M{
		"token_id": tokenId,
	}).One(&token)

	userId := ""

	if err == nil {
		userId = token.Owner
	}
	if err != nil {
		fmt.Println("Access token", tokenId, "not found")
		userId = strings.Replace(c.Request().Header().Get("Authorization"), "Bearer ", "", 1)
	}

	result := User{}
	err = UserCollection.Find(bson.M{
		"_id": userId,
	}).One(&result)

	if err != nil {
		fmt.Println("User", userId, "not found")
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(200, OAuthResponse{
		Id: result.Id,
		Attributes: OAuthAttributes{
			LastName:  result.FamilyName,
			FirstName: result.GivenName,
		},
		Scope: strings.Split(token.Scopes, " "),
	})
}
