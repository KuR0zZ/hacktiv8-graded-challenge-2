package routes

import (
	"graded-challenge-2-client/controller"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, sc controller.ServerController) {
	e.POST("/users/register", sc.Register)
	e.POST("/users/login", sc.Login)
	e.POST("/books", sc.AddBook)
	e.GET("/books/:id", sc.GetBookByID)
	e.PUT("/books/:id", sc.UpdateBook)
	e.DELETE("/books/:id", sc.DeleteBook)
}
