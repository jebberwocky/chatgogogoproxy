package routes

import (
	responses "chatproxy/response"
	validators "chatproxy/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ChatRoute(e *echo.Echo) {
	g := e.Group("/chat")
	g.GET("/*", clickHandler)
}

func clickHandler(c echo.Context) error {
	var req models.AffiliateClickRequest
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	// Validate
	if err := validators.ValidateChatRequest(&req); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	}
	if resp, err := controllers.CreateAffiliateClick(req); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	} else {
		return responses.SendSuccessObj(c, resp)
	}
}
