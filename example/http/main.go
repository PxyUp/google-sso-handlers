package main

import (
	"fmt"
	mcm_google_sso "github.com/PxyUp/google-sso-handlers"
	"log"
	"net/http"
)

type config struct {
}

func NewConfig() *config {
	return &config{}
}

func (s *config) GetHost() string {
	return ""
}

func (s *config) GetClientId() string {
	return ""
}

func (s *config) GetClientSecret() string {
	return ""
}

func (s *config) GetRandomBytesLength() int {
	return 16
}

type redirects struct {
}

func NewRedirects() *redirects {
	return &redirects{}
}

func (r *redirects) GetSuccessRedirectUrl(token string) string {
	return fmt.Sprintf("/login?success=true&token=%s", token)
}

func (r *redirects) GetFailedRedirectUrl(errCode int, err error) string {
	return fmt.Sprintf("/login?success=false&err=%s", err.Error())
}

func (r *redirects) GetCallbackUrl() string {
	return "/google/callback"
}

type user struct {
}

func NewUser() *user {
	return &user{}
}

func (u *user) UserInfoFn(user *mcm_google_sso.GoogleOauthUser) (string, error) {
	return user.Email, nil
}

func initRouter() {
	redirect := NewRedirects()
	oauth, err := mcm_google_sso.NewGoogleOAuth(NewConfig(), redirect, NewUser()).GetGoogleAuthHandler()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/google/login", oauth.LoginHandler)
	http.HandleFunc(redirect.GetCallbackUrl(), oauth.CallbackHandler) // callback url from redirects
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	initRouter()
}