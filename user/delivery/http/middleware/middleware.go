package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

type UserMiddleware struct {
}

func (m *UserMiddleware) ExecTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Content-Type", "application/json")
		start := time.Now()
		c.Response().Header().Add("Trailer", "Execution-Time")
		defer c.Response().Header().Set("Execution-Time", time.Since(start).String())
		return next(c)
	}
}

func InitMiddleware() *UserMiddleware {
	return &UserMiddleware{}
}
