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

	e.GET("/sudo", ah.flagsMiddleware(ah.AdminHandler))
	e.POST("/sudo", ah.flagsMiddleware(ah.AdminHandler))

	e.GET("/logout", ah.flagsMiddleware(ah.LogoutHandler))

	admingroup := e.Group("/su", ah.adminMiddleware)
	admingroup.GET("", ah.AdminPageHandler)
	admingroup.GET("/deleteteam/:id", ah.AdminDeleteTeam)
	admingroup.GET("/deletequestion/:id", ah.AdminDeleteQuestion)
	admingroup.GET("/question", ah.AdminQuestionHandler)
	admingroup.POST("/question", ah.AdminQuestionHandler)

	admingroup.GET("/hints", ah.AdminHintsHandler)
	admingroup.GET("/hints/new", ah.AdminHintNewHandler)
	admingroup.POST("/hints/new", ah.AdminHintNewHandler)

	admingroup.GET("/hints/delete/:id", ah.AdminDeleteHint)
}
