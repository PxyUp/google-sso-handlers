package mcm_google_sso

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"net/http"
)

var (
	SECURE_ERROR = errors.New("Length of the secure bytes must be more than 10")
)

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
	UserInfoFn(user *GoogleOauthUser) (string, error)
}

type GoogleOauthUser struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

type GoogleAuthHandlers struct {
	LoginHandler    func(w http.ResponseWriter, r *http.Request)
	CallbackHandler func(w http.ResponseWriter, r *http.Request)
}

type GoogleOAuth struct {
	config    ConfigController
	redirects RedirectsController
	user      UserController
}

func NewGoogleOAuth(c ConfigController, r RedirectsController, u UserController) *GoogleOAuth {
	return &GoogleOAuth{
		redirects: r,
		user:      u,
		config:    c,
	}
}

func (c *GoogleOAuth) GetGoogleAuthHandler() (*GoogleAuthHandlers, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")

	if err != nil {
		return nil, err
	}

	config := oauth2.Config{
		ClientID:     c.config.GetClientId(),
		ClientSecret: c.config.GetClientSecret(),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		RedirectURL:  fmt.Sprintf("%s%s", c.config.GetHost(), c.redirects.GetCallbackUrl()),
	}
	length := c.config.GetRandomBytesLength()
	if length < 10 {
		return nil, SECURE_ERROR
	}
	bytes, err := generateRandomBytes(length)
	if err != nil {
		return nil, err
	}
	cookieStore := sessions.NewCookieStore(bytes)
	cookieStore.Options = &sessions.Options{
		HttpOnly: true,
	}

	return &GoogleAuthHandlers{
		LoginHandler: func(w http.ResponseWriter, r *http.Request) {
			sessionID := uuid.New().String()
			oauthSession, err := cookieStore.New(r, sessionID)
			if err != nil {
				http.Redirect(w, r, c.redirects.GetFailedRedirectUrl(http.StatusInternalServerError, err), http.StatusMovedPermanently)
				return
			}
			oauthSession.Options.MaxAge = 10 * 60
			err = oauthSession.Save(r, w)
			if err != nil {
				http.Redirect(w, r, c.redirects.GetFailedRedirectUrl(http.StatusInternalServerError, err), http.StatusMovedPermanently)
				return
			}
			http.Redirect(w, r, config.AuthCodeURL(sessionID), http.StatusMovedPermanently)
		},
		CallbackHandler: func(w http.ResponseWriter, r *http.Request) {
			_, err := cookieStore.Get(r, r.URL.Query().Get("state"))

			if err != nil {
				http.Redirect(w, r, c.redirects.GetFailedRedirectUrl(http.StatusInternalServerError, err), http.StatusMovedPermanently)
				return
			}

			oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))

			if err != nil {
				http.Redirect(w, r, c.redirects.GetFailedRedirectUrl(http.StatusInternalServerError, err), http.StatusMovedPermanently)
				return
			}

			userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
			if err != nil {
				http.Redirect(w, r, c.redirects.GetFailedRedirectUrl(http.StatusInternalServerError, err), http.StatusMovedPermanently)
				return
			}

			userGoogle := &GoogleOauthUser{}

			err = userInfo.Claims(userGoogle)

			if err != nil {
				http.Redirect(w, r, c.redirects.GetFailedRedirectUrl(http.StatusInternalServerError, err), http.StatusMovedPermanently)
				return
			}

			token, err := c.user.UserInfoFn(userGoogle)

			if err != nil {
				http.Redirect(w, r, c.redirects.GetFailedRedirectUrl(http.StatusInternalServerError, err), http.StatusMovedPermanently)
				return
			}

			http.Redirect(w, r, c.redirects.GetSuccessRedirectUrl(token), http.StatusMovedPermanently)
		},
	}, nil
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}
