package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Routes struct {
}

func (r Routes) ActivateRoutes(e *echo.Echo) {
	ChatRoute(e)
	HealthyCheckRoutes(e)
	e.Any("*", silly)
}

func HealthyCheckRoutes(e *echo.Echo) {
	g := e.Group("/healthy")
	g.GET("/simple", healthyChecker)
	g.GET("/full", healthyFullChecker)
}

func healthyChecker(c echo.Context) error {
	return c.String(http.StatusOK, "proxy server is healthy.")
}

func healthyFullChecker(c echo.Context) error {
	//TODO: need check the downstream api healthy checking
	return c.String(http.StatusOK, "proxy server is healthy.")
}

func silly(c echo.Context) error {
	return c.String(http.StatusOK, "👋")
}
