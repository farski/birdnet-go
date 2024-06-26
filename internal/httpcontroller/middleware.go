package httpcontroller

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// configureMiddleware sets up middleware for the server.
func (s *Server) configureMiddleware() {
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level:     6,
		MinLength: 2048,
	}))
	// Apply the Cache Control Middleware
	s.Echo.Use(CacheControlMiddleware())
	// Apply a middleware to set the Vary: HX-Request header for all responses
	s.Echo.Use(VaryHeaderMiddleware())
}

// CacheControlMiddleware applies cache control headers for specified routes.
func CacheControlMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			// Apply cache control for assets and clips paths
			if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/clips/") {
				c.Response().Header().Set("Cache-Control", "no-store, max-age=0") // 1 day
			} else {
				// No cache for other routes
				c.Response().Header().Set("Cache-Control", "no-store, max-age=0")
			}
			return next(c)
		}
	}
}

// VaryHeaderMiddleware sets the "Vary: HX-Request" header for all responses.
func VaryHeaderMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Example check: only set Vary header for specific routes or under certain conditions
			if c.Request().Header.Get("HX-Request") != "" {
				c.Response().Header().Set("Vary", "HX-Request")
			}
			return next(c)
		}
	}
}
