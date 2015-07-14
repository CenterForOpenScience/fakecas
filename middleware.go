package main

import "github.com/labstack/echo"

func CorsMiddleWare() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Response().Header().Add("Access-Control-Allow-Origin", "*")
			c.Response().Header().Add("Access-Control-Allow-Headers", "Range, Content-Type, Authorization, Cache-Control, X-Requested-With")
			c.Response().Header().Add("Access-Control-Expose-Headers", "Range, Content-Type, Authorization, Cache-Control, X-Requested-With")
			c.Response().Header().Add("Cache-control", "no-store, no-cache, must-revalidate, max-age=0")

			if c.Request().Method == "OPTIONS" {
				c.Response().Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
				return c.NoContent(204)
			}
			h(c)
			return nil
		}
	}
}
