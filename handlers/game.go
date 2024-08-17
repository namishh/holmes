package handlers

import (
	"errors"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/namishh/holmes/services"
	"github.com/namishh/holmes/views/pages"
	"github.com/namishh/holmes/views/pages/hunt"
)

func (ah *AuthHandler) HomeHandler(c echo.Context) error {
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	if !ok {
		return errors.New("invalid type for key 'ISADMIN'")
	}
	sess, _ := session.Get(auth_sessions_key, c)
	// isError = false
	homeView := pages.Home(fromProtected)
	c.Set("ISERROR", false)
	if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
		return renderView(c, pages.HomeIndex(
			"Home",
			"",
			fromProtected,
			c.Get("ISERROR").(bool),
			homeView,
		))
	}

	return renderView(c, pages.HomeIndex(
		"Home",
		sess.Values[user_name_key].(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		homeView,
	))
}

func (ah *AuthHandler) Hunt(c echo.Context) error {
	questions := make([]services.Question, 0)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	quizview := hunt.Hunt(fromProtected, questions)
	c.Set("ISERROR", false)
	return renderView(c, hunt.HuntIndex(
		"Home",
		c.Get(user_name_key).(string),
		fromProtected,
		c.Get("ISERROR").(bool),
		quizview,
	))
}
