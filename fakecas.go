package main

import (
	"flag"
	"fmt"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
)

var (
	Host            = flag.String("host", "localhost:8080", "The host to bind to")
	DatabaseName    = flag.String("dbname", "osf20130903", "The name of your OSF database")
	DatabaseAddress = flag.String("dbaddress", "localhost:27017", "The address of your mongodb. ie: localhost:27017")
	DatabaseSession mgo.Session
	UserCollection  *mgo.Collection
)

func main() {
	flag.Parse()
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(CorsMiddleWare())

	e.Post("/login", Login)
	e.Get("/logout", Logout)
	e.Get("/oauth2/profile", OAuth)
	e.Get("/p3/serviceValidate", ServiceValidate)

	fmt.Println("Expecting database", *DatabaseName, " to be running at", *DatabaseAddress)
	fmt.Println("Listening on", *Host)

	DatabaseSession, err := mgo.Dial(*DatabaseAddress)
	if err != nil {
		panic(err)
	}
	defer DatabaseSession.Close()

	UserCollection = DatabaseSession.DB(*DatabaseName).C("user")

	e.Run(*Host)
}
