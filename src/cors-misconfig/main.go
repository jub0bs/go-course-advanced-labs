package main

import (
	"crypto/subtle"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	// CORS middleware (configured insecurely)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))
	// Basic Auth middleware
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte("admin")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("admin")) == 1 {
			return true, nil
		}
		return false, nil
	}))
	e.GET("/admin", admin)
	e.Logger.Fatal(e.Start(":8081"))
}

func admin(c echo.Context) error {
	return c.String(http.StatusOK, "Super secret API key: XXXX-XXXX-XXXX")
}
