package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

func ValidateChatRequest(obj interface{}) error {
	return validateDefault(obj)
}
func validateDefault(obj interface{}) error {
	err := Validate.Struct(obj)
	if err != nil {
		var fieldErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, fmt.Sprintf("'%s' failed on the '%s'", e.Field(), e.Tag()))
		}

		errorMessage := fmt.Sprintf("validation failed on fields: %s", strings.Join(fieldErrors, ", "))
		return fmt.Errorf(errorMessage)
	}
	return nil
}
