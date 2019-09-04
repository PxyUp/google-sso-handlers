# GoLang Google SSO handlers

Package for creating authorization with Google on GoLang, can be use in any frameworks.
Have protection from csrf attack.

# Usage

`go get github.com/PxyUp/google-sso-handlers`

```go

type ConfigController interface {
	GetHost() string
	GetClientId() string
	GetClientSecret() string
	GetRandomBytesLength() int
}

type RedirectsController interface {
	GetSuccessRedirectUrl(token string) string
	GetFailedRedirectUrl(errCode int, err error) string
	GetCallbackUrl() string
}

type UserController interface {
	UserInfoFn(user *mcm_google_sso.GoogleOauthUser) (string, error)
}

func initRouter() error {
	oauth, err := mcm_google_sso.NewGoogleOAuth(ConfigController, RedirectsController, UserController).GetGoogleAuthHandler()
    	if err != nil {
    		return nil, err
    	}
        http.HandleFunc("/google/login", oauth.LoginHandler)
	http.HandleFunc("/google/callback", oauth.CallbackHandler) // callback url from redirects
        log.Fatal(http.ListenAndServe(":8080", nil))
}

```


# Examples
- [HTTP](https://github.com/PxyUp/google-sso-handlers/examples/http/main.go)
- [GIN](https://github.com/PxyUp/google-sso-handlers/examples/gin/main.go)