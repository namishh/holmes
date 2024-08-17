package handlers

import (
	"fmt"
	"net/http"

	"github.com/namishh/holmes/views/errors"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)

	var errorPage func(fp bool) templ.Component

	switch code {
	case 401:
		errorPage = errors.Error401
	case 404:
		errorPage = errors.Error404
	case 500:
		errorPage = errors.Error500
	}

	// isError = true
	c.Set("ISERROR", true)

	renderView(c, errors.ErrorIndex(
		fmt.Sprintf("Error (%d)", code),
		c.Get("FROMPROTECTED").(bool),
		errorPage(c.Get("FROMPROTECTED").(bool)),
	))
}

func RouteNotFoundHandler(c echo.Context) error {
	// Hardcoded parameters

	return renderView(c, errors.ErrorIndex(
		fmt.Sprintf("Error (%d)", 404),
		false,
		errors.Error404(false),
	))
}
