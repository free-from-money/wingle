package wingle

import (
	"encoding/xml"
	"net/http"
	"strings"
)

type SesTok struct {
	SesInfo string `xml:"SesInfo"`
	TokInfo string `xml:"TokInfo"`
}

func (st SesTok) GetCookie() *http.Cookie {
	cookieParsed := strings.Split(st.SesInfo, "=")
	c := http.Cookie{Name: cookieParsed[0], Value: cookieParsed[1]}
	return &c
}

type LoginRequest struct {
	XMLName      xml.Name `xml:"request"`
	Text         string   `xml:",chardata"`
	Username     string   `xml:"Username"`
	Password     string   `xml:"Password"`
	PasswordType string   `xml:"password_type"`
}

func NewLoginRequest(token string) LoginRequest {
	const (
		username = "admin"
		pw       = "12345678"
		pwt      = "4"
	)
	return LoginRequest{
		Username:     username,
		Password:     PwCipher(username, pw, token),
		PasswordType: pwt,
	}
}

func PwCipher(name, pw, token string) string {
	return Base64Encode(Sha256(name + Base64Encode(Sha256(pw)) + token))
}
