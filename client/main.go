package main

import (
	"graded-challenge-2-client/controller"
	custom_middleware "graded-challenge-2-client/middleware"
	"graded-challenge-2-client/pb"
	"graded-challenge-2-client/routes"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title           Book Management API
// @version         1.0
// @description     This is a book management API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, skipping....")
	}

	e := echo.New()

	e.Validator = custom_middleware.NewValidate(validator.New())

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	skipper := func(c echo.Context) bool {
		return c.Path() == "/users/login" || c.Path() == "/users/register" || c.Path() == "/swagger/*"
	}

	e.Use(custom_middleware.CustomJwtMiddleware(skipper))

	conn, err := grpc.NewClient(os.Getenv("SERVER_URI"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewServerServiceClient(conn)

	serverController := controller.NewServerController(client)

	routes.Init(e, *serverController)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
