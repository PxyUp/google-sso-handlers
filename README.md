# GoLang Google SSO handlers

![CircleCI](https://circleci.com/gh/PxyUp/google-sso-handlers/tree/master.svg?style=svg)

Package for creating authorization with Google on GoLang, can be use in any frameworks.
Have protection from CSRF attack.

# Usage

`go get github.com/PxyUp/google-sso-handlers`

```go

type ConfigController interface {
	// Get Host address
	GetHost() string 
	// from google console
	GetClientId() string 
	// from google console
	GetClientSecret() string 
	// for generate secret storage
	GetRandomBytesLength() int 
}

type RedirectsController interface {
	// get redirect url without host
	GetSuccessRedirectUrl(token string) string 
	// get redirect url without host
	GetFailedRedirectUrl(errCode int, err error) string 
	// google will be call that url after success login, without host (need provided in google console with host)
	GetCallbackUrl() string 
}

type UserController interface {
	// function must return token(jwt, any access token, and we be passed to RedirectsController.GetSuccessRedirectUrl function
	UserInfoFn(user *mcm_google_sso.GoogleOauthUser) (string, error)
}

func initRouter() {
	oauth, err := mcm_google_sso.NewGoogleOAuth(ConfigController, RedirectsController, UserController).GetGoogleAuthHandler()
    	if err != nil {
    		log.Fatal(err)
    	}
        http.HandleFunc("/google/login", oauth.LoginHandler)
	http.HandleFunc("/google/callback", oauth.CallbackHandler) // callback url from redirects
        log.Fatal(http.ListenAndServe(":8080", nil))
}

```


# Examples
- [HTTP](https://github.com/PxyUp/google-sso-handlers/tree/master/example/http/main.go)
- [GIN](https://github.com/PxyUp/google-sso-handlers/tree/master/example/gin/main.go)
