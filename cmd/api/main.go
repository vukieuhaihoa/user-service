package main

import (
	"github.com/vukieuhaihoa/user-service/internal/infrastructure"
)

// @title User API for Bookmark Management Backend(DDD Version)
// @version 1.2
// @description This is the API documentation for the User service in the Bookmark Management system.
// @host localhost:8080
// @BasePath /
// @schemes http
// @SecurityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
//
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email vukieuhaihoa@gmail.com
func main() {
	app := infrastructure.CreateAPI()

	app.Start()
}
