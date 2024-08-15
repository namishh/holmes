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

	protectedGroup := e.Group("/su", ah.adminMiddleware)
	protectedGroup.GET("", ah.AdminPageHandler)
	protectedGroup.GET("/deleteteam/:id", ah.AdminDeleteTeam)
	protectedGroup.GET("/deletequestion/:id", ah.AdminDeleteQuestion)
	protectedGroup.GET("/question", ah.AdminQuestionHandler)
	protectedGroup.POST("/question", ah.AdminQuestionHandler)
}
