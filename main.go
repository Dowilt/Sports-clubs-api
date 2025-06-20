package main

import (
	"sports-clubs-api/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/clubs", handlers.CreateClub)
	e.PUT("/clubs/:id", handlers.UpdateClub)
	e.DELETE("/clubs/:id", handlers.DeleteClub)

	// Роуты API
	e.GET("/clubs", handlers.GetClubs)

	// Обслуживание статических файлов (index.html и т.д.)
	e.Static("/", "public")

	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}
