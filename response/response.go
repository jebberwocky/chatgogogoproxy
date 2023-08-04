// responses/responses.go

package responses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response struct
// - Used to structure API responses
// - Code: HTTP status code
// - Message: Simple message string
// - Data: Struct or map with response data
type Response struct {
	Message         string      `json:"message"`
	Data            interface{} `json:"data"`
	ChatbotResponse string      `json:"chatbotResponse"`
}

// CustomError struct
// - Used for custom error responses
// - Code: HTTP status code
// - Message: Error message
// - Details: Additional error details
type CustomError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"` // Single value
}

// SendSuccessObj function
// - Sends JSON success response
// - Status 200 OK
// - Response struct payload
func SendSuccessObj(c echo.Context, resp Response) error {
	return c.JSON(http.StatusOK, resp)
}

// SendError function
// - Sends standard error response
// - Uses HTTPError to set status code
func SendError(c echo.Context, status int, message string) error {
	return echo.NewHTTPError(status, message)
}

// SendCustomError function
// - Sends custom error format
// - Sets CustomError struct fields
// - Uses HTTPError to set code and error
func SendCustomError(c echo.Context, message string, details interface{}) error {

	customErr := &CustomError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Details: details, // Pass single value
	}

	return echo.NewHTTPError(customErr.Code, customErr)
}

// Usage:
// - Handler functions call these to generate responses
// - Success cases use SendSuccess
// - Errors use SendError or SendCustomError
