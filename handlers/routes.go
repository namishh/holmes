package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, ah *AuthHandler) {
	e.GET("/", ah.flagsMiddleware(ah.HomeHandler))

	// AUTH ROUTES
	e.GET("/register", ah.flagsMiddleware(ah.RegisterHandler))
	e.POST("/register", ah.flagsMiddleware(ah.RegisterHandler))

	e.GET("/login", ah.flagsMiddleware(ah.LoginHandler))
	e.POST("/login", ah.flagsMiddleware(ah.LoginHandler))

	e.GET("/logout", ah.flagsMiddleware(ah.LogoutHandler))
}
