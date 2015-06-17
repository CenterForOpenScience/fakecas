package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type OAuthAttributes struct {
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
}

type OAuthResponse struct {
	Id         string          `json:"id"`
	Attributes OAuthAttributes `json:"attributes"`
}

type User struct {
	Id         string   `bson:"_id"`
	Username   string   `bson:"username"`
	Emails     []string `bson:"emails"`
	Fullname   string   `bson:"fullname"`
	GivenName  string   `bson:"given_name"`
	FamilyName string   `bson:"family_name"`
}

type ServiceResponse struct {
	Xmlns        string   `xml:"xmlns:cas,attr"`
	XMLName      xml.Name `xml:"cas:serviceResponse"`
	User         string   `xml:"cas:authenticationSuccess>cas:user"`
	NewLogin     bool     `xml:"cas:authenticationSuccess>cas:attributes>cas:isFromNewLogin"`
	Date         string   `xml:"cas:authenticationSuccess>cas:attributes>cas:authenticationDate"`
	GivenName    string   `xml:"cas:authenticationSuccess>cas:attributes>cas:givenName"`
	FamilyName   string   `xml:"cas:authenticationSuccess>cas:attributes>cas:familyName"`
	LongTermAuth bool     `xml:"cas:authenticationSuccess>cas:attributes>cas:longTermAuthenticationRequestTokenUsed"`
	AccessToken  string   `xml:"cas:authenticationSuccess>cas:attributes>accessToken"`
	UserName     string   `xml:"cas:authenticationSuccess>cas:attributes>username"`
}

var (
	host            = flag.String("host", "localhost:8080", "The host to bind to")
	databasename    = flag.String("dbname", "osf20130903", "The name of your OSF database")
	databaseaddress = flag.String("dbaddress", "localhost:27017", "The address of your mongodb. ie: localhost:27017")
)

func main() {
	flag.Parse()

	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/oauth2/profile", oauth)
	http.HandleFunc("/p3/serviceValidate", serviceValidate)

	fmt.Println("Expecting database", *databasename, " to be running at", *databaseaddress)
	fmt.Println("Listening on", *host)

	err := http.ListenAndServe(*host, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	redir, err := url.Parse(r.FormValue("service"))

	if err != nil {
		log.Fatal(err)
	}

	query := redir.Query()
	query.Set("ticket", r.FormValue("username"))
	redir.RawQuery = query.Encode()

	fmt.Println("Logging in and redirecting to", redir)
	http.Redirect(w, r, redir.String(), http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging out and redirecting to", r.FormValue("service"))
	http.Redirect(w, r, r.FormValue("service"), http.StatusFound)
}

func serviceValidate(w http.ResponseWriter, r *http.Request) {

	session, err := mgo.Dial(*databaseaddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(*databasename).C("user")

	result := User{}
	err = c.Find(bson.M{"emails": r.FormValue("ticket")}).One(&result)

	if err != nil {
		fmt.Println("User", r.FormValue("ticket"), "not found.")
		http.NotFound(w, r)
		return
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

	x, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

func oauth(w http.ResponseWriter, r *http.Request) {

	session, err := mgo.Dial(*databaseaddress)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(*databasename).C("user")

	result := User{}
	err = c.Find(bson.M{"_id": strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)}).One(&result)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	js, err := json.Marshal(OAuthResponse{
		Id: result.Id,
		Attributes: OAuthAttributes{
			LastName:  result.FamilyName,
			FirstName: result.GivenName,
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
