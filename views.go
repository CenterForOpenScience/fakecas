package main

import (
	"fmt"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"strings"
)

func Login(c *echo.Context) error {
	redir, err := url.Parse(c.Form("service"))

	if err != nil {
		c.Error(err)
		return nil
	}

	query := redir.Query()
	query.Set("ticket", c.Form("username"))
	redir.RawQuery = query.Encode()

	fmt.Println("Logging in and redirecting to", redir)
	c.Redirect(http.StatusFound, redir.String())
	return nil
}

func Logout(c *echo.Context) error {
	fmt.Println("Logging out and redirecting to", c.Form("service"))
	c.Redirect(http.StatusFound, c.Form("service"))
	return nil
}

func ServiceValidate(c *echo.Context) error {
	result := User{}
	err := UserCollection.Find(bson.M{"emails": c.Form("ticket")}).One(&result)

	if err != nil {
		fmt.Println("User", c.Form("ticket"), "not found.")
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

func OAuth(c *echo.Context) error {
	result := User{}
	err := UserCollection.Find(bson.M{
		"_id": strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1),
	}).One(&result)

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(200, OAuthResponse{
		Id: result.Id,
		Attributes: OAuthAttributes{
			LastName:  result.FamilyName,
			FirstName: result.GivenName,
		},
	})
}
