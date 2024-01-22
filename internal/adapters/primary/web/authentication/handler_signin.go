package authentication

import (
	"net/http"
	"time"

	view "github.com/NordGus/fnncr/internal/adapters/primary/web/views/authentication"
	"github.com/NordGus/fnncr/internal/core/services/authentication"
	"github.com/labstack/echo/v4"
)

func (h *handler) SignInHandlerFunc(c echo.Context) error {
	var (
		username = c.FormValue("username")
		password = c.FormValue("password")

		req = authentication.SignInUserReq{
			Username: username,
			Password: password,
		}

		form = view.FormLogin{
			ActionURL: "/sign_in",
			Username:  username,
			Password:  password,
			Failed:    true,
		}

		cookie = &http.Cookie{
			Name:    "_session_fnncr",
			Path:    "/",
			Domain:  "localhost",
			Expires: time.Now().Add(7 * 24 * time.Hour),

			MaxAge:   int((7 * 24 * time.Hour).Seconds()),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
	)

	resp, err := h.api.SignInUser(c.Request().Context(), req)
	if err != nil {
		c.Logger().Error(err)

		return view.SignIn(form).Render(c.Request().Context(), c.Response())
	}

	cookie.Value = resp.SessionID
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}