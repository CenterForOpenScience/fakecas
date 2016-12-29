package main

import (
	"database/sql"
	"flag"
	"fmt"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"html/template"
	"os"
)

var Version string

var (
	Host               = flag.String("host", "localhost:8080", "The host to bind to")
	OSFHost            = flag.String("osfhost", "localhost:5000", "The osf host to bind to")
	DatabaseName       = flag.String("dbname", "osf", "The name of your OSF database")
	DatabaseAddress    = flag.String("dbaddress", "postgres://postgres@localhost:5432/osf?sslmode=disable", "The address of your postgres instance. ie: postgres://user:pass@127.0.0.1/dbname?other=args")
	DatabaseConnection *sql.DB
)

func main() {
	fmt.Printf("Starting FakeCAS %s\n", Version)
	flag.Parse()
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${response_time} ${response_size}\n",
		Output: os.Stdout,
	}))
	e.Use(middleware.Recover())

  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:   []string{"*"},
		AllowMethods:   []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:   []string{"Range", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:   []string{"Range", "Content-Type", "Authorization", "X-Requested-With"},
	}))

	t, err := template.New("login").Parse(LOGINPAGE)
	if err != nil {
		panic(err)
	}
	temp := &Template{templates: t}
  e.Renderer = temp

	e.GET("/login", LoginGET)
	e.POST("/login", LoginPOST)
	e.GET("/logout", Logout)
	e.GET("/oauth2/profile", OAuth)
	e.GET("/p3/serviceValidate", ServiceValidate)

	fmt.Println("Expecting database", *DatabaseName, "to be running at", *DatabaseAddress)

	DatabaseConnection, err = sql.Open("postgres", *DatabaseAddress)

	if err != nil {
		panic(err)
	}

	defer DatabaseConnection.Close()

  e.Start(*Host)
}
