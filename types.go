package main

import (
	"encoding/xml"
	"html/template"
)

type OAuthAttributes struct {
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
}

type OAuthResponse struct {
	Id         string          `json:"id"`
	Attributes OAuthAttributes `json:"attributes"`
	Scope      []string        `json:"scope"`
}

type User struct {
	Id              string `bson:"_id"`
	Username        string `bson:"username"`
	Fullname        string `bson:"fullname"`
	GivenName       string `bson:"given_name"`
	FamilyName      string `bson:"family_name"`
	IsRegistered    bool   `bson:"is_registered"`
	Password        string `bson:"password"`
	VerificationKey string `bson:"verification_key"`
}

type ServiceResponse struct {
	Xmlns            string   `xml:"xmlns:cas,attr"`
	XMLName          xml.Name `xml:"cas:serviceResponse"`
	User             string   `xml:"cas:authenticationSuccess>cas:user"`
	NewLogin         bool     `xml:"cas:authenticationSuccess>cas:attributes>cas:isFromNewLogin"`
	Date             string   `xml:"cas:authenticationSuccess>cas:attributes>cas:authenticationDate"`
	GivenName        string   `xml:"cas:authenticationSuccess>cas:attributes>cas:givenName"`
	FamilyName       string   `xml:"cas:authenticationSuccess>cas:attributes>cas:familyName"`
	LongTermAuth     bool     `xml:"cas:authenticationSuccess>cas:attributes>cas:longTermAuthenticationRequestTokenUsed"`
	AccessToken      string   `xml:"cas:authenticationSuccess>cas:attributes>accessToken"`
	AccessTokenScope string   `xml:"cas:authenticationSuccess>cas:attributes>accessTokenScope"`
	UserName         string   `xml:"cas:authenticationSuccess>cas:attributes>username"`
}

type AccessToken struct {
	Id      string `bson:"_id"`
	Owner   string `bson:"owner"`
	TokenId string `bson:"token_id"`
	Scopes  string `bson:"scopes"`
}

type Template struct {
	templates *template.Template
}

type TemplateGlobal struct {
	// login flow
	LoginForm          bool
	RegisterForm       bool
	NotExist           bool
	NotValid           bool
	NotAuthorized      bool
	NotRegistered      bool
	RegisterSuccessful bool
	ShowErrorMessages  bool
	// cas url
	CASLogin    string
	CASRegister string
	// osf url
	OSFDomain             string
	OSFForgotPassword     string
	OSFInstitutionLogin   string
	OSFResendConfirmation string
}
