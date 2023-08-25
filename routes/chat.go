package routes

import (
	"chatproxy/controllers"
	"chatproxy/models"
	responses "chatproxy/response"
	"chatproxy/util"
	validators "chatproxy/validator"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ChatRoute(e *echo.Echo) {
	g := e.Group("/chat")
	g.POST("/v4", v4Handler)
	g.POST("/newmonkey", v3Handler)
	g.POST("/finetune", vFineHandler)
	g.POST("*", defaultHandler)
}

func defaultHandler(c echo.Context) error {
	var req models.ChatRequest
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	app := c.Get(util.EchoAppContext).(models.AppContext)
	// Validate
	if err := validators.ValidateChatRequest(&req); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	}
	if resp, err := controllers.DefaultHandle(req, app); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	} else {
		return responses.SendSuccessObj(c, resp)
	}
}

func v3Handler(c echo.Context) error {
	var req models.ChatRequest
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	app := c.Get(util.EchoAppContext).(models.AppContext)
	// Validate
	if err := validators.ValidateChatRequest(&req); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	}
	if resp, err := controllers.V3Handle(req, app); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	} else {
		return responses.SendSuccessObj(c, resp)
	}
}

func v4Handler(c echo.Context) error {
	var req models.ChatRequest
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "bad request")
	}
	app := c.Get(util.EchoAppContext).(models.AppContext)
	// Validate
	if err := validators.ValidateChatRequest(&req); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	}
	if resp, err := controllers.V4Handle(req, app); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	} else {
		return responses.SendSuccessObj(c, resp)
	}
}

func vFineHandler(c echo.Context) error {
	var req models.ChatRequest
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "bad request")
	}
	app := c.Get(util.EchoAppContext).(models.AppContext)
	// Validate
	if err := validators.ValidateChatRequest(&req); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	}
	if resp, err := controllers.VFineHandle(req, app); err != nil {
		return responses.SendError(c, http.StatusBadRequest, err.Error())
	} else {
		return responses.SendSuccessObj(c, resp)
	}
}
