package main

import (
	"fmt"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"strings"
)

func Login(c echo.Context) error {
	redir, err := url.Parse(c.FormValue("service"))

	if err != nil {
		c.Error(err)
		return nil
	}

	result := User{}

	if err = UserCollection.Find(bson.M{"username": c.FormValue("username")}).One(&result); err != nil {
		fmt.Println("User", c.FormValue("ticket"), "not found.")
		return c.NoContent(http.StatusNotFound)
	}

	if !result.IsRegistered {
		return c.HTML(200, UNREGISTERED)
	}

	query := redir.Query()
	query.Set("ticket", c.FormValue("username"))
	redir.RawQuery = query.Encode()

	fmt.Println("Logging in and redirecting to", redir)
	return c.Redirect(http.StatusFound, redir.String())
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
	tokenId :=  strings.Replace(c.Request().Header().Get("Authorization"), "Bearer ", "", 1)
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
