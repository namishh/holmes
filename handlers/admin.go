package handlers

import (
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/namishh/holmes/services"
	"github.com/namishh/holmes/views/pages"
	"github.com/namishh/holmes/views/pages/panel"
)

func (ah *AuthHandler) adminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)
		if user, ok := sess.Values[user_type].(string); !ok || user != "admin" {
			c.Set("FROMPROTECTED", false)
			c.Set("ISADMIN", false)
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		c.Set("ISADMIN", true)

		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			c.Set("FROMPROTECTED", false)
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Please provide valid credentials")
		}

		if userId, ok := sess.Values[user_id_key].(int); ok && userId != 0 {
			c.Set(user_id_key, userId) // set the user_id in the context
		}

		if username, ok := sess.Values[user_name_key].(string); ok && len(username) != 0 {
			c.Set(user_name_key, username) // set the username in the context
		}

		if tzone, ok := sess.Values[tzone_key].(string); ok && len(tzone) != 0 {
			c.Set(tzone_key, tzone) // set the client's time zone in the context
		}

		// fromProtected = true
		c.Set("FROMPROTECTED", true)

		return next(c)
	}
}

func (ah *AuthHandler) AdminHandler(c echo.Context) error {

	sess, _ := session.Get(auth_sessions_key, c)
	if user, ok := sess.Values[user_type].(string); !ok || user == "admin" {
		return c.Redirect(http.StatusSeeOther, "/su")
	}

	errs := make(map[string]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {

		if c.FormValue("password") != os.Getenv("ADMIN_PASS") {
			c.Set("ISERROR", true)
			errs["pass"] = "Incorrect Password"

			adminLoginView := pages.AdminLogin(fromProtected, errs)
			c.Set("ISERROR", false)
			return renderView(c, pages.AdminLoginIndex(
				"Admin Panel",
				"admin",
				fromProtected,
				c.Get("ISERROR").(bool),
				adminLoginView,
			))
		} else {
			tzone := ""
			if len(c.Request().Header["X-Timezone"]) != 0 {
				tzone = c.Request().Header["X-Timezone"][0]
			}

			sess, _ := session.Get(auth_sessions_key, c)
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   60 * 60 * 24 * 7, // 1 week
				HttpOnly: true,
			}

			// Set user as authenticated, their username,
			// their ID and the client's time zone

			sess.Values = map[interface{}]interface{}{
				auth_key:      true,
				user_type:     "admin",
				user_id_key:   9999999,
				user_name_key: "admin",
				tzone_key:     tzone,
			}
			sess.Save(c.Request(), c.Response())

			return c.Redirect(http.StatusSeeOther, "/su")
		}

	}

	//sess, _ := session.Get(auth_sessions_key, c)
	// isError = false
	adminLoginView := pages.AdminLogin(fromProtected, errs)
	c.Set("ISERROR", false)
	return renderView(c, pages.AdminLoginIndex(
		"Admin Panel",
		"admin",
		fromProtected,
		c.Get("ISERROR").(bool),
		adminLoginView,
	))
}

func (ah *AuthHandler) AdminPageHandler(c echo.Context) error {
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	users := make([]services.User, 0)
	questions := make([]services.Question, 0)

	users, err := ah.UserServices.GetAllUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching users")
	}

	adminLoginView := panel.PanelHome(fromProtected, users, questions)
	c.Set("ISERROR", false)
	return renderView(c, panel.PanelIndex(
		"Admin Panel",
		"admin",
		fromProtected,
		c.Get("ISERROR").(bool),
		adminLoginView,
	))
}
