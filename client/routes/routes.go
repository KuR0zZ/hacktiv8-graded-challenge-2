package routes

import (
	"graded-challenge-2-client/controller"

	_ "graded-challenge-2-client/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init(e *echo.Echo, sc controller.ServerController) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/users/register", sc.Register)
	e.POST("/users/login", sc.Login)
	e.POST("/books", sc.AddBook)
	e.GET("/books/:id", sc.GetBookByID)
	e.PUT("/books/:id", sc.UpdateBook)
	e.DELETE("/books/:id", sc.DeleteBook)
	e.POST("/books/:id", sc.BorrowBook)
}
